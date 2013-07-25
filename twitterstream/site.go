// Copyright 2013 The go-twitterstream AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package twitterstream

import (
	"net/url"
)

type SiteStreams struct {
	client *Client
}

func (s *SiteStreams) Get(f map[string]string) error {
	baseURL, _ := url.Parse("https://sitestream.twitter.com/1.1/")
	s.client.baseURL = baseURL

	u := "site.json"

	params := []string{"follow", "with", "replies"}
	body := make(map[string]string)
	for _, p := range params {
		if v, exists := f[p]; exists && v != "" {
			body[p] = v
		}
	}
	body["stall_warnings"] = "true"

	req, err := s.client.NewRequest("POST", u, body)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}

	return s.client.DispatchResponse(resp)
}
