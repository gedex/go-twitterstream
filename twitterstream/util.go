package twitterstream

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"math/rand"
	"sort"
	"strings"
)

func escape(s string) string {
	hexCount := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if shouldEscape(c) {
			hexCount++
		}
	}

	if hexCount == 0 {
		return s
	}

	t := make([]byte, len(s)+2*hexCount)
	j := 0
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case shouldEscape(c):
			t[j] = '%'
			t[j+1] = "0123456789ABCDEF"[c>>4]
			t[j+2] = "0123456789ABCDEF"[c&15]
			j += 3
		default:
			t[j] = s[i]
			j++
		}
	}
	return string(t)
}

func shouldEscape(c byte) bool {
	if 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z' || '0' <= c && c <= '9' {
		return false
	}
	switch c {
	case '-', '_', '.', '~':
		return false
	case '$', '&', '+', ',', '/', ':', ';', '=', '?', '@':
		return true
	}
	return true
}

func Nonce(n int) string {
	var alphanum = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	buf := make([]byte, n)
	for i := 0; i < n; i++ {
		buf[i] = alphanum[rand.Intn(len(alphanum))]
	}
	return string(buf)
}

func SignatureBaseString(rp *RequestParams) string {
	var length int

	length += len(rp.Query)
	length += len(rp.Body)
	length += len(rp.OAuth)

	var keys []string
	params := make(map[string]string, length)
	for k, v := range rp.Query {
		keys = append(keys, k)
		params[k] = v
	}
	for k, v := range rp.Body {
		keys = append(keys, k)
		params[k] = v
	}
	for k, v := range rp.OAuth {
		keys = append(keys, k)
		params[k] = v
	}
	sort.Strings(keys)

	ps := ""
	for _, k := range keys {
		ps += fmt.Sprintf("%s=%s&", k, params[k])
	}
	ps = strings.Trim(ps, "&")

	sbs := fmt.Sprintf("%s&%s&%s", rp.Method, escape(rp.Endpoint), escape(ps))

	return sbs
}

func Signature(consumerSecret, tokenSecret, sbs string) (string, error) {
	key := []byte(escape(consumerSecret) + "&" + escape(tokenSecret))
	h := hmac.New(sha1.New, key)
	defer h.Reset()

	_sbs := []byte(sbs)
	n, err := h.Write(_sbs)
	if n != len(_sbs) || err != nil {
		return "", err
	}

	bb := new(bytes.Buffer)
	enc := base64.NewEncoder(base64.StdEncoding, bb)
	enc.Write(h.Sum(nil))
	enc.Close()

	return bb.String(), nil
}
