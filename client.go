package bingx

import (
	"context"
	"encoding/json"
	"go-bingx/common"
	"io/ioutil"
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

func (c *Client) debug(message string, args ...interface{}) {
	if c.Debug {
		c.Logger.Printf(message, args...)
	}
}

func (c *Client) callAPI(ctx context.Context, r *request) (data []byte, err error) {
	req, err := http.NewRequest(r.method, r.fullUrl, r.body)
	if err != nil {
		return []byte{}, err
	}
	req = req.WithContext(ctx)
	req.Header = r.header

	f := c.do
	if f == nil {
		f = c.HTTPClient.Do
	}
	res, err := f(req)
	if err != nil {
		return []byte{}, err
	}
	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}
	defer func() {
		cerr := res.Body.Close()
		if err == nil && cerr != nil {
			err = cerr
		}
	}()

	if res.StatusCode >= http.StatusBadRequest {
		apiErr := new(common.APIError)
		e := json.Unmarshal(data, apiErr)
		if e != nil {
			c.debug("failed to unmarshal json: %s", e)
		}
		return nil, apiErr
	}
	return data, nil
}
