package main

import (
	"anaconda"
	"database/sql"
	"fmt"
	"log"
	_ "mysql"
	"net/url"
)

func main() {
	db, err := sql.Open("mysql", "root:@/messages?charset=utf8")
	if err != nil {
		log.Fatal(err)
	}

	stmtIns, err := db.Prepare("INSERT INTO tweets VALUES(NULL, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmtIns.Close()

	anaconda.SetConsumerKey("Y13sFl9oNt7sS6NdUtQw")
	anaconda.SetConsumerSecret("d2D4qFilNRQL3ho9uClKzzxERhkpxRah8a45NT5OiI")

	api := anaconda.NewTwitterApi("19683-Wd51DLxYGmQmNqGxOVKQqs7w31Ho2GfazxUkwCHds", "vSfKzTlQz3HZXSO0bMjg4yXVYgUodMyc7l378USeM8")

	var max_id int64 = -1
	var last_max_id int64 = 0
	var since_id int64 = -1 // TODO: handle since_id

	for {
		if last_max_id == max_id {
			break
		}

		v := url.Values{}
		v.Set("screen_name", "rjungemann")
		v.Set("count", "200")

		if since_id != -1 {
			v.Set("since_id", fmt.Sprintf("%v", since_id))
		}

		if max_id != -1 {
			v.Set("max_id", fmt.Sprintf("%v", max_id - 1))
		}

		searchResult, err := api.GetUserTimeline(v)
		if err != nil {
			log.Fatal(err)
		}

		last_max_id = max_id

		for _, tweet := range searchResult {
			user := tweet.User
			user_id := user.Id
			created_at := tweet.Created_at
			remote_id := tweet.Id
			text := tweet.Text

			_, err = stmtIns.Exec(user_id, created_at, remote_id, text)
			if err != nil {
				log.Fatal(err)
			}

			max_id = tweet.Id
		}
	}
}
