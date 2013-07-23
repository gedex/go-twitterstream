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

// NewClient returns a new Twitter Streaming client. It expects
// conf with valid credentials.
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

// RequestParams represents parameters used when requesting stream
// to any stream endpoints.
type RequestParams struct {
	Method   string
	Endpoint string
	Query    map[string]string
	Body     map[string]string
	OAuth    map[string]string
}

// NewRequest cretes a strema request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the baseURL of the Client.
// Relative URLs should always be specified without a preceding slash. The value
// of body is url encoded and included as the request body if specified.
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

// Do sends a stream request and returns the stream response. The stream
// response consists of a series of newline-delimited messages, where
// "newline" is considered to be \r\n (in hex, 0x0D 0x0A) and "message"
// is a JSON encoded data structure or a blank line. The return values
// should always be consumed by DispatchResponse.
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
// to ProcessStream until client is closed.
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

// Disconnect closes the client from the stream.
func (c *Client) Disconnect() {
	c.closed = true
}

// streamSwitcher unmarshall the raw into general container
// which then decoded into more specific type if matches with
// any defined stream type.
func (c *Client) streamSwitcher(raw []byte) {
	var v map[string]interface{}

	if err := json.Unmarshal(raw, &v); err != nil {
		log.Printf("twitterstream: error unmarshal stream: %v\n", err)
		return
	}

	stream := &Stream{Raw: raw}

	var container interface{}
	if _, ok := v["control"]; ok {
	} else if _, ok := v["warning"]; ok {
		stream.Type = "warning"
		container = new(WarningNotice)
		stream.WarningNotice = container.(*WarningNotice)
	} else if d, ok := v["delete"]; ok {
		ds := d.(map[string]interface{})
		if _, ok = ds["status"]; ok {
			stream.Type = "delete"
			container = new(TweetDeletionNotice)
			stream.TweetDeletionNotice = container.(*TweetDeletionNotice)
		}
	} else if sg, ok := v["scrub_geo"]; ok {
		sgc := sg.(map[string]interface{})
		if _, ok = sgc["up_to_status_id"]; ok {
			stream.Type = "scrub_geo"
			container = new(LocationDeletionNotice)
			stream.LocationDeletionNotice = container.(*LocationDeletionNotice)
		}
	} else if _, ok := v["limit"]; ok {
		stream.Type = "limit"
		container = new(LimitNotice)
		stream.LimitNotice = container.(*LimitNotice)
	} else if _, ok := v["direct_message"]; ok {
		stream.Type = "direct_message"
		container = new(DirectMessageNotice)
		stream.DirectMessageNotice = container.(*DirectMessageNotice)
	} else if _, ok := v["status_withheld"]; ok {
		stream.Type = "status_withheld"
		container = new(StatusWithheldNotice)
		stream.StatusWithheldNotice = container.(*StatusWithheldNotice)
	} else if _, ok := v["user_withheld"]; ok {
		stream.Type = "user_withheld"
		container = new(UserWithheldNotice)
		stream.UserWithheldNotice = container.(*UserWithheldNotice)
	} else if _, ok := v["event"]; ok {
		stream.Type = "event"
		container = new(Event)
		stream.Event = container.(*Event)
	} else if _, ok := v["friends"]; ok {
		stream.Type = "friends"
		container = new(FriendsLists)
		stream.FriendsLists = container.(*FriendsLists)
	} else if _, ok := v["text"]; ok {
		if _, ok = v["user"]; ok {
			stream.Type = "tweet"
			container = new(Tweet)
			stream.Tweet = container.(*Tweet)
		}
	} else if _, ok := v["for_user"]; ok {
		stream.Type = "for_user"
		container = new(ForUser)
		stream.ForUser = container.(*ForUser)
	}

	go c.handleStream(stream, container)
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

// Stream represents a twitter stream.
type Stream struct {
	Raw                    []byte
	Type                   string
	Tweet                  *Tweet
	WarningNotice          *WarningNotice
	LimitNotice            *LimitNotice
	TweetDeletionNotice    *TweetDeletionNotice
	LocationDeletionNotice *LocationDeletionNotice
	DirectMessageNotice    *DirectMessageNotice
	FriendsLists           *FriendsLists
	ForUser                *ForUser
	UserWithheldNotice     *UserWithheldNotice
	Event                  *Event
	StatusWithheldNotice   *StatusWithheldNotice
}

var availableStreamTypes = map[string]bool{
	"control":         true,
	"warning":         true,
	"scrub_geo":       true,
	"tweet":           true,
	"limit":           true,
	"delete":          true,
	"friends":         true,
	"direct_message":  true,
	"status_withheld": true,
	"user_withheld":   true,
	"for_user":        true,
}

var defaultStreamHandlers = map[string]func(*Stream){
	"tweet":   defaultTweetStreamHandler,
	"friends": defaultFriendsStreamHandler,
}

func defaultTweetStreamHandler(s *Stream) {
	log.Printf("@%v: %v\n", s.Tweet.User.ScreenName, s.Tweet.Text)
}
func defaultFriendsStreamHandler(s *Stream) {
	log.Printf("@%+v\n", s.FriendsLists.Friends)
}

// handle registers the stream handler for the given stream type.
// If a handler already exists for targetted stream type, handle
// panics.
func (mux *ProcessStreamMux) handle(streamType string, handler Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	if handler == nil {
		panic("twitterstream: nil handler")
	}
	if _, exists := mux.m[streamType]; exists {
		panic(fmt.Sprintf("twitterstream: multiple registrations for %v", streamType))
	}

	mux.m[streamType] = muxEntry{h: handler, streamType: streamType}
}

// handler returns the handler to use for the given stream type.
// If no registered handler found, it checks for default stream
// handler. Otherwise nil is returned.
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

// handleStream handles the stream.
func (c *Client) handleStream(stream *Stream, container interface{}) {
	h := c.streamHandleMux.handler(stream.Type)
	if h == nil {
		log.Printf("twitterstream: No handler for %v stream", stream.Type)
	} else {
		if container != nil {
			err := json.Unmarshal(stream.Raw, container)
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
	valid := isValidStreamType(streamType)
	if !valid {
		panic("twitterstream: unknown stream type " + streamType)
	}
	c.streamHandleMux.handle(streamType, handlerFunc(handler))
}
