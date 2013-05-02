This is an example Go app which will scrape a user's tweets. Eventually it will
be the basis of a markov-powered "ebooks" Twitter bot.

Run the following SQL to create the necessary database and table.

    CREATE DATABASE messages;

    USE messages;

    CREATE TABLE `tweets` (
      `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
      `user_id` varchar(255) DEFAULT NULL,
      `created_at` datetime DEFAULT NULL,
      `remote_id` varchar(255) DEFAULT NULL,
      `text` text,
      PRIMARY KEY (`id`),
      UNIQUE KEY `unique_remote_id_on_tweets` (`remote_id`)
    ) ENGINE=InnoDB AUTO_INCREMENT=2293 DEFAULT CHARSET=utf8

1. Run `source bashrc` to set `$GOPATH` to the current directory.
2. Run `rake setup` to fetch the dependencies.
3. Finally run `rake` to build, run, then clean the app. Alternatively, you can run:
  1. `rake build`
  2. `rake run` or `./scraper`
  3. `rake clean`

