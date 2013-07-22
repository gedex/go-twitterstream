package main

import (
	"github.com/gedex/go-twitterstream/twitterstream"
	"log"
)

func main() {
	config := &twitterstream.Config{
		ConsumerKey:      "ohBNaRJK7MQrRuBw0SbwQ",
		ConsumerSecret:   "68M17oEE70Yg6ActFJgtulLu2NJi6ZjYDPVLKBAVwYc",
		OAuthToken:       "1106913162-fRKyqX9LcLINTMZ59w8fq0vmoA7Reh6eyuMcQzD",
		OAuthTokenSecret: "nKiH5o7ZTy0nGn0DaiNOEzF1pV5VitiWTbrsjK0nExM",
	}
	client := twitterstream.NewClient(config)
	client.HandleFunc("tweet", func(s *twitterstream.Stream) {
		log.Printf("user %v tweets: %v\n", s.Tweet.User.ScreenName, s.Tweet.Text)
	})
	client.HandleFunc("limit", func(s *twitterstream.Stream) {
		log.Printf("limit notice: %v\n", s.LimitNotice)
	})
	client.HandleFunc("delete", func(s *twitterstream.Stream) {
		log.Printf("delete notice %+v\n", s.DeleteNotice.Delete.Status)
	})

	sf := map[string]string{
		"track": "macet, jakarta",
	}
	err := client.Public.Filter(sf)
	if err != nil {
		log.Fatal(err)
	}
}
