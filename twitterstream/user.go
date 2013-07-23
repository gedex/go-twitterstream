// Copyright 2013 The go-twitterstream AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package twitterstream

import (
	"net/url"
)

type UserStreams struct {
	client *Client
}

func (s *UserStreams) Get(f map[string]string) error {
	baseURL, _ := url.Parse("https://userstream.twitter.com/1.1/")
	s.client.baseURL = baseURL

	u := "user.json"

	params := url.Values{
		"stall_warnings": {"true"},
	}
	if v, exists := f["with"]; exists {
		params.Add("with", v)
	}
	if v, exists := f["replies"]; exists {
		params.Add("replies", v)
	}
	if v, exists := f["track"]; exists {
		params.Add("track", v)
	}
	if v, exists := f["locations"]; exists {
		params.Add("locations", v)
	}
	u += "?" + params.Encode()

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}

	return s.client.DispatchResponse(resp)
}
