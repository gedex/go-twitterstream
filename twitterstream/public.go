// Copyright 2013 The go-twitterstream AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package twitterstream

type PublicStreams struct {
	client *Client
}

func (s *PublicStreams) Sample() error {
	u := "statuses/sample.json?stall_warnings=true"

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

func (s *PublicStreams) Filter(f map[string]string) error {
	u := "statuses/filter.json"

	params := []string{"follow", "track", "locations"}
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

func (s *PublicStreams) Firehose() error {
	u := "statuses/firehose.json?stall_warnings=true"

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
