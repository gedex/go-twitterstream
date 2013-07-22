// Copyright 2013 The go-twitterstream AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
TODO: put example usage here
*/

package twitterstream

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const (
	// Version represents this library version
	Version = "0.1"

	// DefaultBaseURL represents default Twitter Stream base URL
	DefaultBaseURL = "https://stream.twitter.com/1.1/"

	// UserAgent represents default client User-Agent
	DefaultUserAgent = "go-twitterstream/" + Version

	MaxReconnects = 10
)

// Client manages communication with Twitter stream.
type Client struct {
	// HTTP client used to communicate with the stream
	client *http.Client

	// Base URL for stream requests.
	baseURL *url.URL

	// Client's config
	config *Config

	// Streaming endpoints
	Public *PublicStreams
	User   *UserStreams
	Site   *SiteStreams

	streamHandleMux *ProcessStreamMux

	// Set to true by Disonnect
	closed bool

	// Reconnection
	reconnectCount   int
	reconnectTimeout int
}

// NewClient returns a new Twitter Streaming client.
func NewClient(conf *Config) *Client {

	_baseURL := DefaultBaseURL
	if conf.BaseURL != "" {
		_baseURL = conf.BaseURL
	}
	baseURL, _ := url.Parse(_baseURL)

	c := &Client{
		config:          conf,
		client:          http.DefaultClient,
		baseURL:         baseURL,
		streamHandleMux: &ProcessStreamMux{m: make(map[string]muxEntry)},
	}
	c.Public = &PublicStreams{client: c}
	c.User = &UserStreams{client: c}
	c.Site = &SiteStreams{client: c}
	return c
}

type RequestParams struct {
	Method   string
	Endpoint string
	Query    map[string]string
	Body     map[string]string
	OAuth    map[string]string
}

func (c *Client) NewRequest(method, urlStr string, body map[string]string) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.baseURL.ResolveReference(rel)

	params := new(RequestParams)
	params.Method = method
	params.Endpoint = u.Scheme + "://" + u.Host + u.Path

	q := u.Query()
	reqQuery := make(map[string]string)
	for k := range q {
		reqQuery[escape(k)] = escape(q.Get(k))
	}
	params.Query = reqQuery

	reqBody := ""
	if body != nil {
		bp := make(map[string]string)
		for k, v := range body {
			ek, ev := escape(k), escape(v)
			bp[ek] = ev
			reqBody += fmt.Sprintf("%s=%s&", url.QueryEscape(k), url.QueryEscape(body[k]))
		}
		params.Body = bp
		reqBody = strings.Trim(reqBody, "&")
	}

	req, err := http.NewRequest(method, u.String(), strings.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	ua := DefaultUserAgent
	if c.config.UserAgent != "" {
		ua = c.config.UserAgent
	}

	if method == "POST" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	req.Header.Add("User-Agent", ua)
	req.Header.Add("Authorization", c.config.authorizationHeader(params))

	return req, nil
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	err = CheckResponse(resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// DispatchResponse reads http.Response and dispatches the chunk
// to ProcessStream.
func (c *Client) DispatchResponse(r *http.Response) error {
	reader := bufio.NewReader(r.Body)
	for {
		if c.closed {
			r.Body.Close()
		}

		line, err := reader.ReadBytes('\n')
		if err != nil {
			// TODO: Check if connection stale

			// TODO: try reconnect if stall or connection error

			return err
		}
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		go c.streamSwitcher(line)
	}
}

func (c *Client) Disconnect() {
	c.closed = true
}

func (c *Client) streamSwitcher(raw []byte) {
	var v map[string]interface{}

	if err := json.Unmarshal(raw, &v); err != nil {
		log.Printf("twitterstream: error unmarshal stream: %v\n", err)
		return
	}

	stream := &Stream{Raw: raw}

	// TODO: Complete all stream types handlers
	if _, ok := v["control"]; ok {
	} else if _, ok := v["warning"]; ok {
	} else if d, ok := v["delete"]; ok {
		ds := d.(map[string]interface{})
		if _, ok = ds["status"]; ok {
			stream.Type = "delete"
			go c.handleStream(stream)
		}
	} else if sg, ok := v["scrub_geo"]; ok {
		sgc := sg.(map[string]interface{})
		if _, ok = sgc["up_to_status_id"]; ok {

		}
	} else if _, ok := v["limit"]; ok {
		stream.Type = "limit"
		go c.handleStream(stream)
	} else if _, ok := v["direct_message"]; ok {

	} else if _, ok := v["status_withheld"]; ok {

	} else if _, ok := v["user_withheld"]; ok {

	} else if _, ok := v["event"]; ok {

	} else if _, ok := v["friends"]; ok {
		stream.Type = "friends"
		go c.handleStream(stream)
	} else if _, ok := v["text"]; ok {
		if _, ok = v["user"]; ok {
			stream.Type = "tweet"
			go c.handleStream(stream)
		}
	} else if _, ok := v["for_user"]; ok {

	}
}

// ProcessStreamMux is stream multiplexer.
// It matches each incoming stream against a list of registered
// stream type and calls the handler for the pattern that matches
// the stream type.
type ProcessStreamMux struct {
	mu sync.RWMutex
	m  map[string]muxEntry
}

type muxEntry struct {
	streamType string
	h          Handler
}

type Stream struct {
	Raw          []byte
	Type         string
	Tweet        *Tweet
	LimitNotice  *LimitNotice
	DeleteNotice *TweetDeletionNotice
	FriendsLists *FriendsLists
	// TODO: event, friends, etc
}

var availableStreamTypes = map[string]bool{
	"tweet":   true,
	"limit":   true,
	"delete":  true,
	"friends": true,
}

var defaultStreamHandlers = map[string]func(*Stream){
	"tweet": defaultTweetStreamHandler,
}

func defaultTweetStreamHandler(s *Stream) {
	log.Printf("@%v: %v\n", s.Tweet.User.ScreenName, s.Tweet.Text)
}

// handle registers the stream handler for the given stream type.
// If a handler already exists for targetted stream type, handle
// panics. It also panics if a given stream type is an invalid
// stream type.
func (mux *ProcessStreamMux) handle(streamType string, handler Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	// TODO: Add helper to check if given streamType is a valid one
	if streamType == "" {
		panic(fmt.Sprintf("twitterstream: invalid stream type %v", streamType))
	}
	if handler == nil {
		panic("twitterstream: nil handler")
	}
	if _, exists := mux.m[streamType]; exists {
		panic(fmt.Sprintf("twitterstream: multiple registrations for %v", streamType))
	}

	mux.m[streamType] = muxEntry{h: handler, streamType: streamType}
}

// handler returns the handler to use for the given stream type.
func (mux *ProcessStreamMux) handler(streamType string) Handler {
	mux.mu.RLock()
	defer mux.mu.RUnlock()

	if h, exists := mux.m[streamType]; exists {
		return h.h
	}
	if h, exists := defaultStreamHandlers[streamType]; exists {
		return handlerFunc(h)
	}
	return nil
}

// Objects implementing the Handler interface can be
// registered to handle particular stream from Twitter.
type Handler interface {
	ProcessStream(*Stream)
}

// The handlerFunc type is an adapter to allow the use of
// odinary functions as stream handlers. If f is a function
// with the appropriate signature, handlerFunc(f) is a
// Handler object that calls f.
type handlerFunc func(*Stream)

// ProcessStream calls f(stream)
func (f handlerFunc) ProcessStream(stream *Stream) {
	f(stream)
}

func (c *Client) handleStream(stream *Stream) {
	h := c.streamHandleMux.handler(stream.Type)
	if h == nil {
		log.Printf("twitterstream: No handler for %v stream", stream.Type)
	} else {
		// TODO: supports all stream types
		var v interface{}
		switch stream.Type {
		case "limit":
			v = new(LimitNotice)
			stream.LimitNotice = v.(*LimitNotice)
		case "delete":
			v = new(TweetDeletionNotice)
			stream.DeleteNotice = v.(*TweetDeletionNotice)
		case "tweet":
			v = new(Tweet)
			stream.Tweet = v.(*Tweet)
		case "friend":
			v = new(FriendsLists)
			stream.FriendsLists = v.(*FriendsLists)
		}
		if v != nil {
			err := json.Unmarshal(stream.Raw, v)
			if err != nil {
				log.Printf("twitterstream: Error unmarshall: %v", err)
			} else {
				h.ProcessStream(stream)
			}
		}
	}
}

// HandleFunc registers the stream handler function for the given stream type.
func (c *Client) HandleFunc(streamType string, handler func(*Stream)) {
	v, exists := availableStreamTypes[streamType]
	if !v || !exists {
		panic("twitterstream: unknown stream type " + streamType)
	}
	c.streamHandleMux.handle(streamType, handlerFunc(handler))
}
