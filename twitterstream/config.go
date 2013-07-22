// Copyright 2013 The go-twitterstream AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package twitterstream

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	ConsumerKey      string
	ConsumerSecret   string
	OAuthToken       string
	OAuthTokenSecret string
	UserAgent        string
	BaseURL          string
}

func (conf *Config) authorizationHeader(rp *RequestParams) string {
	op := map[string]string{
		"oauth_nonce":            Nonce(42),
		"oauth_token":            conf.OAuthToken,
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        strconv.FormatInt(time.Now().Unix(), 10),
		"oauth_consumer_key":     conf.ConsumerKey,
		"oauth_version":          "1.0",
	}
	rp.OAuth = op

	// This will be sorted later
	var keys []string

	for k, v := range op {
		ke := escape(k)
		op[ke] = escape(v)
		keys = append(keys, ke)
	}

	signature, err := Signature(conf.ConsumerSecret, conf.OAuthTokenSecret, SignatureBaseString(rp))
	if err != nil {
		log.Printf("twitterstream: error generating signature %s\n", err)
	}
	op["oauth_signature"] = escape(signature)
	keys = append(keys, "oauth_signature")

	sort.Strings(keys)

	authStr := "OAuth "
	for _, k := range keys {
		authStr += fmt.Sprintf("%s=\"%s\", ", k, op[k])
	}
	authStr = strings.Trim(authStr, ", ")
	return authStr
}
