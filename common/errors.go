package common

import "fmt"

// API error when response status is 4xx or 5xx
type APIError struct {
	Code    int64  `json:"code"`
	Message string `json:"msg"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("<APIError> code=%d, msg=%s", e.Code, e.Message)
}
