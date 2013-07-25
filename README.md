go-twitterstream
================

Go library to access Twitter's Streaming API.

**Documentation**: http://godoc.org/github.com/gedex/go-twitterstream/twitterstream

**Build Status**: [![Build Status](https://travis-ci.org/gedex/go-twitterstream.png?branch=master)](https://travis-ci.org/gedex/go-twitterstream)

## Basic Usage

~~~go
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

err := client.Public.Sample()
if err != nil {
	log.Fatal(err)
}
~~~

Please see [examples](./examples) for a complete example.

## Credits

* [twitterstream](twitterstream) for Go
* [tweetstream](https://github.com/tweetstream/tweetstream) RubyGem
* [Twitter Platform Documentation](https://dev.twitter.com/docs)

## License

This library is distributed under the BSD-style license found in the LICENSE.md file.
