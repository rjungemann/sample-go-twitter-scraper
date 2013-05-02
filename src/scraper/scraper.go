package main

import (
	"anaconda"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	_ "mysql"
	"net/url"
	"strconv"
	"time"
)

type Config struct {
	Twitter_consumer_key        string
	Twitter_consumer_secret     string
	Twitter_access_token        string
	Twitter_access_token_secret string
	Twitter_username            string
	Database                    string
}

func FormatTwitterDateForMysql(twitter_date string) string {
	// In Go, you use strings representing a "standard date" to define date
	// layouts, instead of using something like strftime
	twitter_date_layout := "Mon Jan 02 15:04:05 -0700 2006"
	mysql_date_layout := "2006-01-02 15:04:05"

	// Get the date relative to a timezone
	relative_date, err := time.Parse(twitter_date_layout, twitter_date)
	if err != nil {
		log.Fatal(err)
	}

	// Make the date UTC before converting it to a MySQL representation
	return relative_date.UTC().Format(mysql_date_layout)
}

func main() {
	// Load JSON config file
	json_string, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}

	// Populate contents of Config struct with data from JSON config file
	var config Config
	err = json.Unmarshal(json_string, &config)
	if err != nil {
		log.Fatal(err)
	}

	// Load the database
	db, err := sql.Open("mysql", config.Database)
	if err != nil {
		log.Fatal(err)
	}

	// Compile insert statement to be used to populate a tweet in the database
	//
	// Note that REPLACE will deduplicate tweets on insert
	//
	stmtIns, err := db.Prepare("REPLACE tweets VALUES(NULL, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmtIns.Close()

	rows, err := db.Query("SELECT remote_id FROM tweets ORDER BY created_at DESC LIMIT 1;")
	if err != nil {
		log.Fatal(err)
	}

	// Initialize variables used for paginating tweets
	var max_id int64 = -1
	var last_max_id int64 = 0
	var since_id int64 = -1

	// Fetch the since_id if the database is populated
	var remote_id_string string
	for rows.Next() {
		err = rows.Scan(&remote_id_string)
		if err != nil {
			log.Fatal(err)
		}
		since_id, err = strconv.ParseInt(remote_id_string, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		break
	}

	// Load Twitter client
	anaconda.SetConsumerKey(config.Twitter_consumer_key)
	anaconda.SetConsumerSecret(config.Twitter_consumer_secret)
	api := anaconda.NewTwitterApi(config.Twitter_access_token, config.Twitter_access_token_secret)

	for {
		// If last_max_id equals max_id then we're done
		if last_max_id == max_id {
			break
		}

		// Prepare our request to Twitter API
		v := url.Values{}
		v.Set("screen_name", config.Twitter_username)
		v.Set("count", "200")
		if since_id != -1 {
			// Expects since_id to be a string for some stupid reason
			v.Set("since_id", fmt.Sprintf("%v", since_id))
		}
		if max_id != -1 {
			// Expects max_id to be a string for some stupid reason
			v.Set("max_id", fmt.Sprintf("%v", max_id-1))
		}

		// Get next 200 tweets
		searchResult, err := api.GetUserTimeline(v)
		if err != nil {
			log.Fatal(err)
		}

		// Insert fetched tweets into the database
		last_max_id = max_id
		for _, tweet := range searchResult {
			_, err = stmtIns.Exec(tweet.User.Id, FormatTwitterDateForMysql(tweet.Created_at), tweet.Id, tweet.Text)
			if err != nil {
				log.Fatal(err)
			}
			max_id = tweet.Id
		}
	}
}
