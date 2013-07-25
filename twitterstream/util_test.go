// Copyright 2013 The go-twitterstream AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package twitterstream

import (
	"testing"
)

type EscapeTest struct {
	in  string
	out string
	err error
}

var escapeTests = []EscapeTest{
	{
		"",
		"",
		nil,
	},
	{
		"Ladies + Gentlemen",
		"Ladies%20%2B%20Gentlemen",
		nil,
	},
	{
		"An encoded string!",
		"An%20encoded%20string%21",
		nil,
	},
	{
		"Dogs, Cats & Mice",
		"Dogs%2C%20Cats%20%26%20Mice",
		nil,
	},
	{
		"☃",
		"%E2%98%83",
		nil,
	},
}

func Test_escape(t *testing.T) {
	for _, tt := range escapeTests {
		actual := escape(tt.in)
		if tt.out != actual {
			t.Errorf("escape(%q) = %q, want %q", tt.in, actual, tt.out)
		}
	}
}

type ValidStreamTypeTest struct {
	in  string
	out bool
}

var validStreamTypeTests = []ValidStreamTypeTest{
	{
		"control",
		true,
	},
	{
		"warning",
		true,
	},
	{
		"scrub_geo",
		true,
	},
	{
		"tweet",
		true,
	},
	{
		"limit",
		true,
	},
	{
		"delete",
		true,
	},
	{
		"friends",
		true,
	},
	{
		"direct_message",
		true,
	},
	{
		"status_withheld",
		true,
	},
	{
		"user_withheld",
		true,
	},
	{
		"for_user",
		true,
	},
	{
		"invalid",
		false,
	},
}

func Test_isValidStreamType(t *testing.T) {
	for _, tt := range validStreamTypeTests {
		actual := isValidStreamType(tt.in)
		if tt.out != actual {
			t.Errorf("isValidStreamType(%q) = %q, want %q", tt.in, actual, tt.out)
		}
	}
}

func TestNonce(t *testing.T) {
	var prev string
	for i := 1; i <= 10; i++ {
		nonce := Nonce(42 - i)
		// Make sure Nonce always different
		if nonce == prev {
			t.Error("Nonce failed to generate unique string")
		}
		if len(nonce) != 42-i {
			t.Errorf("Nonce(%d) generate string with length %d, want %d", 42-i, len(nonce), 42-i)
		}
	}
}

type SignatureBaseStringTest struct {
	in  *RequestParams
	out string
}

var oauthParamTest = map[string]string{
	"oauth_nonce":            "BpLnfgDsc2WD8F2qNfHK5a84jjJkwzDkh9h2fhfUVu",
	"oauth_token":            "1106913162-fRKyqX9LcLINTMZ59w8fq0vmoA7Reh6eyuMcQzD",
	"oauth_signature_method": "HMAC-SHA1",
	"oauth_timestamp":        "1374766315",
	"oauth_consumer_key":     "ohBNaRJK7MQrRuBw0SbwQ",
	"oauth_version":          "1.0",
}

var signatureBaseStringTests = []SignatureBaseStringTest{
	{
		&RequestParams{
			Method:   "POST",
			Endpoint: "https://stream.twitter.com/1.1/statuses/filter.json",
			Body: map[string]string{
				"track":          escape("jakarta, macet"),
				"stall_warnings": "true",
			},
			OAuth: oauthParamTest,
		},
		"POST&https%3A%2F%2Fstream.twitter.com%2F1.1%2Fstatuses%2Ffilter.json&oauth_consumer_key%3DohBNaRJK7MQrRuBw0SbwQ%26oauth_nonce%3DBpLnfgDsc2WD8F2qNfHK5a84jjJkwzDkh9h2fhfUVu%26oauth_signature_method%3DHMAC-SHA1%26oauth_timestamp%3D1374766315%26oauth_token%3D1106913162-fRKyqX9LcLINTMZ59w8fq0vmoA7Reh6eyuMcQzD%26oauth_version%3D1.0%26stall_warnings%3Dtrue%26track%3Djakarta%252C%2520macet",
	},
	{
		&RequestParams{
			Method:   "GET",
			Endpoint: "https://stream.twitter.com/1.1/statuses/sample.json",
			Query:    map[string]string{"stall_warnings": "true"},
			OAuth:    oauthParamTest,
		},
		"GET&https%3A%2F%2Fstream.twitter.com%2F1.1%2Fstatuses%2Fsample.json&oauth_consumer_key%3DohBNaRJK7MQrRuBw0SbwQ%26oauth_nonce%3DBpLnfgDsc2WD8F2qNfHK5a84jjJkwzDkh9h2fhfUVu%26oauth_signature_method%3DHMAC-SHA1%26oauth_timestamp%3D1374766315%26oauth_token%3D1106913162-fRKyqX9LcLINTMZ59w8fq0vmoA7Reh6eyuMcQzD%26oauth_version%3D1.0%26stall_warnings%3Dtrue",
	},
	{
		&RequestParams{
			Method:   "GET",
			Endpoint: "https://stream.twitter.com/1.1/statuses/firehose.json",
			Query:    map[string]string{"stall_warnings": "true"},
			OAuth:    oauthParamTest,
		},
		"GET&https%3A%2F%2Fstream.twitter.com%2F1.1%2Fstatuses%2Ffirehose.json&oauth_consumer_key%3DohBNaRJK7MQrRuBw0SbwQ%26oauth_nonce%3DBpLnfgDsc2WD8F2qNfHK5a84jjJkwzDkh9h2fhfUVu%26oauth_signature_method%3DHMAC-SHA1%26oauth_timestamp%3D1374766315%26oauth_token%3D1106913162-fRKyqX9LcLINTMZ59w8fq0vmoA7Reh6eyuMcQzD%26oauth_version%3D1.0%26stall_warnings%3Dtrue",
	},
	{
		&RequestParams{
			Method:   "GET",
			Endpoint: "https://userstream.twitter.com/1.1/user.json",
			Query: map[string]string{
				"stall_warnings": "true",
				"track":          escape("golang, go-nuts"),
				"locations":      escape("-122.75,36.8,-121.75,37.8"),
			},
			OAuth: oauthParamTest,
		},
		"GET&https%3A%2F%2Fuserstream.twitter.com%2F1.1%2Fuser.json&locations%3D-122.75%252C36.8%252C-121.75%252C37.8%26oauth_consumer_key%3DohBNaRJK7MQrRuBw0SbwQ%26oauth_nonce%3DBpLnfgDsc2WD8F2qNfHK5a84jjJkwzDkh9h2fhfUVu%26oauth_signature_method%3DHMAC-SHA1%26oauth_timestamp%3D1374766315%26oauth_token%3D1106913162-fRKyqX9LcLINTMZ59w8fq0vmoA7Reh6eyuMcQzD%26oauth_version%3D1.0%26stall_warnings%3Dtrue%26track%3Dgolang%252C%2520go-nuts",
	},
	{
		&RequestParams{
			Method:   "POST",
			Endpoint: "https://sitestream.twitter.com/1.1/site.json",
			Query: map[string]string{
				"stall_warnings": "true",
				"follow":         escape("1,2,3,4,5"),
			},
			OAuth: oauthParamTest,
		},
		"POST&https%3A%2F%2Fsitestream.twitter.com%2F1.1%2Fsite.json&follow%3D1%252C2%252C3%252C4%252C5%26oauth_consumer_key%3DohBNaRJK7MQrRuBw0SbwQ%26oauth_nonce%3DBpLnfgDsc2WD8F2qNfHK5a84jjJkwzDkh9h2fhfUVu%26oauth_signature_method%3DHMAC-SHA1%26oauth_timestamp%3D1374766315%26oauth_token%3D1106913162-fRKyqX9LcLINTMZ59w8fq0vmoA7Reh6eyuMcQzD%26oauth_version%3D1.0%26stall_warnings%3Dtrue",
	},
}

func TestSignatureBaseString(t *testing.T) {
	for _, tt := range signatureBaseStringTests {
		actual := SignatureBaseString(tt.in)
		if actual != tt.out {
			t.Errorf("SignatureBaseString(%v) = %s, want %s", tt.in, actual, tt.out)
		}
	}
}

type SignatureTest struct {
	sbs string
	cs  string
	ts  string
	out string
}

var signatureTests = []SignatureTest{
	{
		"POST&https%3A%2F%2Fstream.twitter.com%2F1.1%2Fstatuses%2Ffilter.json&oauth_consumer_key%3DohBNaRJK7MQrRuBw0SbwQ%26oauth_nonce%3DBpLnfgDsc2WD8F2qNfHK5a84jjJkwzDkh9h2fhfUVu%26oauth_signature_method%3DHMAC-SHA1%26oauth_timestamp%3D1374766315%26oauth_token%3D1106913162-fRKyqX9LcLINTMZ59w8fq0vmoA7Reh6eyuMcQzD%26oauth_version%3D1.0%26stall_warnings%3Dtrue%26track%3Djakarta%252C%2520macet",
		"68M17oEE70Yg6ActFJgtulLu2NJi6ZjYDPVLKBAVwYc",
		"nKiH5o7ZTy0nGn0DaiNOEzF1pV5VitiWTbrsjK0nExM",
		"83fNTcywyCMiWjAwxBnQIakYQDA=",
	},
	{
		"GET&https%3A%2F%2Fstream.twitter.com%2F1.1%2Fstatuses%2Fsample.json&oauth_consumer_key%3DohBNaRJK7MQrRuBw0SbwQ%26oauth_nonce%3DBpLnfgDsc2WD8F2qNfHK5a84jjJkwzDkh9h2fhfUVu%26oauth_signature_method%3DHMAC-SHA1%26oauth_timestamp%3D1374766315%26oauth_token%3D1106913162-fRKyqX9LcLINTMZ59w8fq0vmoA7Reh6eyuMcQzD%26oauth_version%3D1.0%26stall_warnings%3Dtrue",
		"68M17oEE70Yg6ActFJgtulLu2NJi6ZjYDPVLKBAVwYc",
		"nKiH5o7ZTy0nGn0DaiNOEzF1pV5VitiWTbrsjK0nExM",
		"VdJmI7wEdXfPWuhfsClhkUBiZGA=",
	},
	{
		"GET&https%3A%2F%2Fstream.twitter.com%2F1.1%2Fstatuses%2Ffirehose.json&oauth_consumer_key%3DohBNaRJK7MQrRuBw0SbwQ%26oauth_nonce%3DBpLnfgDsc2WD8F2qNfHK5a84jjJkwzDkh9h2fhfUVu%26oauth_signature_method%3DHMAC-SHA1%26oauth_timestamp%3D1374766315%26oauth_token%3D1106913162-fRKyqX9LcLINTMZ59w8fq0vmoA7Reh6eyuMcQzD%26oauth_version%3D1.0%26stall_warnings%3Dtrue",
		"68M17oEE70Yg6ActFJgtulLu2NJi6ZjYDPVLKBAVwYc",
		"nKiH5o7ZTy0nGn0DaiNOEzF1pV5VitiWTbrsjK0nExM",
		"+tWwuSCnZHV/tMKFSFWvTXNY+AM=",
	},
	{
		"GET&https%3A%2F%2Fuserstream.twitter.com%2F1.1%2Fuser.json&locations%3D-122.75%252C36.8%252C-121.75%252C37.8%26oauth_consumer_key%3DohBNaRJK7MQrRuBw0SbwQ%26oauth_nonce%3DBpLnfgDsc2WD8F2qNfHK5a84jjJkwzDkh9h2fhfUVu%26oauth_signature_method%3DHMAC-SHA1%26oauth_timestamp%3D1374766315%26oauth_token%3D1106913162-fRKyqX9LcLINTMZ59w8fq0vmoA7Reh6eyuMcQzD%26oauth_version%3D1.0%26stall_warnings%3Dtrue%26track%3Dgolang%252C%2520go-nuts",
		"68M17oEE70Yg6ActFJgtulLu2NJi6ZjYDPVLKBAVwYc",
		"nKiH5o7ZTy0nGn0DaiNOEzF1pV5VitiWTbrsjK0nExM",
		"iMqjTUp/5sk69H8Ojbv4XCIboIo=",
	},
	{
		"POST&https%3A%2F%2Fsitestream.twitter.com%2F1.1%2Fsite.json&follow%3D1%252C2%252C3%252C4%252C5%26oauth_consumer_key%3DohBNaRJK7MQrRuBw0SbwQ%26oauth_nonce%3DBpLnfgDsc2WD8F2qNfHK5a84jjJkwzDkh9h2fhfUVu%26oauth_signature_method%3DHMAC-SHA1%26oauth_timestamp%3D1374766315%26oauth_token%3D1106913162-fRKyqX9LcLINTMZ59w8fq0vmoA7Reh6eyuMcQzD%26oauth_version%3D1.0%26stall_warnings%3Dtrue",
		"68M17oEE70Yg6ActFJgtulLu2NJi6ZjYDPVLKBAVwYc",
		"nKiH5o7ZTy0nGn0DaiNOEzF1pV5VitiWTbrsjK0nExM",
		"5c+kD/dIEoTHiDKl00rDC37AecE=",
	},
}

func TestSignature(t *testing.T) {
	for _, tt := range signatureTests {
		actual, _ := Signature(tt.cs, tt.ts, tt.sbs)
		if actual != tt.out {
			t.Errorf("Signature(%s, %s, %s) = %s, want %s", tt.cs, tt.ts, tt.sbs, actual, tt.out)
		}
	}
}
