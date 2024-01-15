package bingx

import (
	"log"
	"net/http"
	"os"
)

// API Endpoints
const (
	baseApiUrl = "https://open-api.bingx.com"
	// baseTestApiUrl = "" Unactual
)

func getApiEndpoint() string {
	return baseApiUrl
}

type doFunc func(*http.Request) (*http.Response, error)

// Client define API client
type Client struct {
	APIKey     string
	SecretKey  string
	BaseURL    string
	UserAgent  string
	HTTPClient *http.Client
	Debug      bool
	Logger     *log.Logger
	TimeOffset int64
	do         doFunc
}

// Init Api Client from apiKey & secretKey
func NewClient(apiKey, secretKey string) *Client {
	return &Client{
		APIKey:     apiKey,
		SecretKey:  secretKey,
		BaseURL:    getApiEndpoint(),
		UserAgent:  "Bingx/golang",
		HTTPClient: http.DefaultClient,
		Logger:     log.New(os.Stderr, "bingx-golang", log.LstdFlags),
	}
}
