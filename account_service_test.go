package bingx

import (
	"net/http"
	"testing"
)

func clientTestDo(req *http.Request) (*http.Response, error) {
	code := http.StatusOK
	data := 
	return &http.Response{
		Body:       req.Body,
		StatusCode: code,
	}, nil
}

func TestAccountService(t *testing.T) {
	client := newTestClient()

}
