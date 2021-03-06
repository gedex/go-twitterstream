package main

import (
	"github.com/gedex/go-twitterstream/twitterstream"
	"log"
)

func main() {
	config := &twitterstream.Config{
		ConsumerKey:      "YOUR CONSUMER KEY",
		ConsumerSecret:   "YOUR CONSUMER SECRET",
		OAuthToken:       "YOUR OAUTH TOKEN",
		OAuthTokenSecret: "YOUR OAUTH TOKEN SECRET",
	}
	client := twitterstream.NewClient(config)
	client.HandleFunc("tweet", func(s *twitterstream.Stream) {
		log.Printf("user %v tweets: %v\n", s.Tweet.User.ScreenName, s.Tweet.Text)
	})

	sf := map[string]string{
		"track": "macet, jakarta",
	}
	err := client.Public.Filter(sf)
	if err != nil {
		log.Fatal(err)
	}
}
