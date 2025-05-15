package bingx

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
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

type OrderWorkingType string

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

	MarkOrderWorkingType     OrderWorkingType = "MARK_PRICE"
	ContractOrderWorkingType OrderWorkingType = "CONTRACT_PRICE"
	IndexOrderWorkingType    OrderWorkingType = "INDEX_PRICE"
)

type Interval string

const (
	Interval1  Interval = "1m"
	Interval3  Interval = "3m"
	Interval5  Interval = "5m"
	Interval15 Interval = "15m"
	Interval30 Interval = "30m"

	Interval60  Interval = "1h"
	Interval2h  Interval = "2h"
	Interval4h  Interval = "4h"
	Interval6h  Interval = "6h"
	Interval8h  Interval = "8h"
	Interval12h Interval = "12h"

	Interval1d Interval = "1d"
	Interval3d Interval = "3d"

	Interval1w Interval = "1w"

	Interval1M Interval = "1M"
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

	recvWindow := r.recvWindow
	if recvWindow == 0 {
		recvWindow = 10000
	}

	r.setParam(recvWindowKey, recvWindow)

	timestamp := time.Now().UnixNano() / 1e6
	if r.query != nil {
		r.setParam(timestampKey, timestamp)
		c.debug(r.query.Encode())
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
	// c.debug("request: %#v", req)
	c.debug("request url: %#v", req.URL.String())
	f := c.do
	if f == nil {
		f = c.HTTPClient.Do
	}
	res, err := f(req)
	if err != nil {
		return []byte{}, err
	}
	data, err = io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}
	defer func() {
		cerr := res.Body.Close()
		if err == nil && cerr != nil {
			err = cerr
		}
	}()
	// c.debug("response: %#v", res)
	c.debug("response body: %s", string(data))
	c.debug("response status code: %d", res.StatusCode)

	apiErr := new(common.APIError)
	json.Unmarshal(data, apiErr)

	if apiErr.Code != 0 {
		return nil, apiErr
	}
	return data, nil
}

type GetServerTimeService struct {
	c *Client
}

func (s *GetServerTimeService) Do(ctx context.Context, opts ...RequestOption) (res int64, err error) {
	r := &request{method: http.MethodGet, endpoint: "/openApi/swap/v2/server/time"}

	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return 0, err
	}

	resp := new(struct {
		Code int              `json:"code"`
		Msg  string           `json:"msg"`
		Data map[string]int64 `json:"data"`
	})

	err = json.Unmarshal(data, &resp)

	if err != nil {
		return 0, err
	}

	res = resp.Data["serverTime"]

	return res, nil
}

func (c *Client) NewGetServerTimeService() *GetServerTimeService {
	return &GetServerTimeService{c: c}
}

func (c *Client) NewGetBalanceService() *GetBalanceService {
	return &GetBalanceService{c: c}
}

func (c *Client) NewGetAccountListenKeyService() *GetAccountListenKeyService {
	return &GetAccountListenKeyService{c: c}
}

func (c *Client) NewGetOpenPositionsService() *GetOpenPositionsService {
	return &GetOpenPositionsService{c: c}
}

func (c *Client) NewCreateOrderService() *CreateOrderService {
	return &CreateOrderService{c: c}
}

func (c *Client) NewCancelOrderService() *CancelOrderService {
	return &CancelOrderService{c: c}
}

func (c *Client) NewGetOrderService() *GetOrderService {
	return &GetOrderService{c: c}
}

func (c *Client) NewGetOpenOrdersService() *GetOpenOrdersService {
	return &GetOpenOrdersService{c: c}
}

func (c *Client) NewGetKlinesService() *GetKlinesService {
	return &GetKlinesService{c: c}
}

func (c *Client) NewGetSymbolDataService() *GetSymbolDataService {
	return &GetSymbolDataService{c: c}
}

func (c *Client) NewCancelAllOrdersService() *CancelAllOrdersService {
	return &CancelAllOrdersService{c: c}
}
