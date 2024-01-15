package bingx

import (
	"io"
	"net/http"
	"net/url"
)

// API request
type request struct {
	method   string
	endpoint string
	query    url.Values
	header   http.Header
	body     io.Reader
	fullUrl  string
}

// Init new get request
func NewGetRequest(endpoint string) *request {
	r := &request{
		method:   "GET",
		endpoint: endpoint,
	}

	return r
}

// Init new post request
func NewPostRequest(endpoint string, body io.Reader) *request {
	r := &request{
		method:   "POST",
		endpoint: endpoint,
		body:     body,
	}

	return r
}
