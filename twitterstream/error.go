package twitterstream

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type ErrorReponse struct {
	Response *http.Response
	Message  string
}

func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	err := &ErrorReponse{Response: r}

	var msg string
	switch r.StatusCode {
	case http.StatusUnauthorized:
		msg = "Unauthorized"
	case http.StatusForbidden:
		msg = "Forbidden"
	case http.StatusNotAcceptable:
		msg = "Not acceptable"
	case http.StatusRequestEntityTooLarge:
		msg = "A parameter list is too long"
	case http.StatusRequestedRangeNotSatisfiable:
		msg = "Range Unacceptable"
	case 420:
		msg = "Rate Limited"
	default:
		msg = "Unknown"
	}

	err.Message = msg

	return err
}

func (e *ErrorReponse) Error() string {
	defer e.Response.Body.Close()
	body, _ := ioutil.ReadAll(e.Response.Body)
	return fmt.Sprintf("twitterstream: response error: %v, response body: %v", e.Message, string(body))
}
