package bingx

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/magicaleks/go-bingx/common"
)

// API Endpoints
const (
	baseApiUrl = "https://open-api.bingx.com"
	// baseTestApiUrl = "" Unactual
)

// Side type of order
type SideType string

// PositionSide type of order
type PositionSideType string

// Type of order
type OrderType string

type OrderStatus string

type OrderSpecType string

const (
	timestampKey  = "timestamp"
	signatureKey  = "signature"
	recvWindowKey = "recvWindow"

	BuySideType  SideType = "BUY"
	SellSideType SideType = "SELL"

	ShortPositionSideType PositionSideType = "SHORT"
	LongPositionSideType  PositionSideType = "LONG"
	BothPositionSideType  PositionSideType = "BOTH"

	LimitOrderType  OrderType = "LIMIT"
	MarketOrderType OrderType = "MARKET"

	NewOrderStatus             OrderStatus = "NEW"
	PartiallyFilledOrderStatus OrderStatus = "PARTIALLY_FILLED"
	FilledOrderStatus          OrderStatus = "FILLED"
	CanceledOrderStatus        OrderStatus = "CANCELED"
	ExpiredOrderStatus         OrderStatus = "EXPIRED"

	NewOrderSpecType        OrderSpecType = "NEW"
	CanceledOrderSpecType   OrderSpecType = "CANCELED"
	CalculatedOrderSpecType OrderSpecType = "CALCULATED"
	ExpiredOrderSpecType    OrderSpecType = "EXPIRED"
	TradeOrderSpecType      OrderSpecType = "TRADE"
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

func (c *Client) parseRequest(r *request, opts ...RequestOption) (err error) {
	for _, opt := range opts {
		opt(r)
	}

	err = r.validate()
	if err != nil {
		return err
	}

	r.setParam(recvWindowKey, r.recvWindow)

	timestamp := time.Now().UnixNano() / 1e6
	if r.query != nil {
		r.setParam(timestampKey, timestamp)
		sign := computeHmac256(r.query.Encode(), c.SecretKey)
		r.setParam(signatureKey, sign)
	} else {
		r.setFormParam(timestampKey, timestamp)
		sign := computeHmac256(r.form.Encode(), c.SecretKey)
		r.setFormParam(signatureKey, sign)
	}

	queryString := r.query.Encode()
	body := &bytes.Buffer{}
	bodyString := r.form.Encode()
	header := http.Header{}
	if r.header != nil {
		header = r.header.Clone()
	}
	if bodyString != "" {
		header.Set("Content-Type", "application/x-www-form-urlencoded")
		body = bytes.NewBufferString(bodyString)
	}

	fullUrl := fmt.Sprintf("%s%s", c.BaseURL, r.endpoint)

	if queryString != "" {
		fullUrl = fmt.Sprintf("%s?%s", fullUrl, queryString)
	}

	header.Add("X-BX-APIKEY", c.APIKey)

	r.fullUrl = fullUrl
	r.header = header
	r.body = body
	return nil
}

func computeHmac256(strMessage string, strSecret string) string {
	key := []byte(strSecret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(strMessage))
	return hex.EncodeToString(h.Sum(nil))
}

func (c *Client) callAPI(ctx context.Context, r *request, opts ...RequestOption) (data []byte, err error) {
	err = c.parseRequest(r, opts...)
	if err != nil {
		return []byte{}, err
	}
	req, err := http.NewRequest(r.method, r.fullUrl, r.body)
	if err != nil {
		return []byte{}, err
	}
	req = req.WithContext(ctx)
	req.Header = r.header
	c.debug("request: %#v", req)
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
	c.debug("response: %#v", res)
	c.debug("response body: %s", string(data))
	c.debug("response status code: %d", res.StatusCode)

	apiErr := new(common.APIError)
	json.Unmarshal(data, apiErr)

	if apiErr.Code != 0 {
		return nil, apiErr
	}
	return data, nil
}

func (c *Client) NewGetBalanceService() *GetBalanceService {
	return &GetBalanceService{c}
}

func (c *Client) NewGetAccountListenKeyService() *GetAccountListenKeyService {
	return &GetAccountListenKeyService{c: c}
}

func (c *Client) NewGetOpenPositionsService() *GetOpenPositionsService {
	return &GetOpenPositionsService{c: c}
}

func (c *Client) NewOpenOrderService() *CreateOrderService {
	return &CreateOrderService{c: c}
}

func (c *Client) NewCancelOrderService() *CancelOrderService {
	return &CancelOrderService{c: c}
}
