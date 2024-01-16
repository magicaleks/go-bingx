package bingx

import (
	"context"
	"encoding/json"
	"net/http"
)

// Only Market and Limit orders supported
type OpenOrderService struct {
	c          *Client
	symbol     string
	orderType  OrderType
	side       SideType
	reduceOnly string
	price      float64
	quantity   float64
}

func (s *OpenOrderService) Symbol(symbol string) *OpenOrderService {
	s.symbol = symbol
	return s
}

func (s *OpenOrderService) Type(orderType OrderType) *OpenOrderService {
	s.orderType = orderType
	return s
}

func (s *OpenOrderService) Side(side SideType) *OpenOrderService {
	s.side = side
	return s
}

func (s *OpenOrderService) ReduceOnly() *OpenOrderService {
	s.reduceOnly = "true"
	return s
}

func (s *OpenOrderService) Price(price float64) *OpenOrderService {
	s.price = price
	return s
}

func (s *OpenOrderService) Quantity(quantity float64) *OpenOrderService {
	s.quantity = quantity
	return s
}

type OpenOrderResponse struct {
	OrderId int `json:"orderId"`
}

func (s *OpenOrderService) Do(ctx context.Context, opts ...RequestOption) (res *OpenOrderResponse, err error) {
	r := &request{method: http.MethodPost, endpoint: "/openApi/swap/v2/trade/order"}

	if s.symbol != "" {
		r.addParam("symbol", s.symbol)
	}

	if s.orderType != "" {
		r.addParam("type", s.orderType)
	}

	if s.side != "" {
		r.addParam("side", s.side)
	}

	if s.reduceOnly != "" {
		r.addParam("reduceOnly", s.reduceOnly)
	}

	if s.price != 0 {
		r.addParam("price", s.price)
	}

	if s.quantity != 0 {
		r.addParam("quantity", s.quantity)
	}

	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}

	resp := new(struct {
		code int
		msg  string
		data map[string]*OpenOrderResponse
	})

	err = json.Unmarshal(data, &resp)
	res = resp.data["order"]

	if err != nil {
		return nil, err
	}

	return res, nil
}

type CancelOrderService struct {
	c           *Client
	symbol      string
	orderIdList []int
}

func (s *CancelOrderService) Symbol(symbol string) *CancelOrderService {
	s.symbol = symbol
	return s
}

func (s *CancelOrderService) Orders(orderIds ...int) *CancelOrderService {
	for orderId := range orderIds {
		s.orderIdList = append(s.orderIdList, orderId)
	}

	return s
}

// Define order
type Order struct {
	Time          int    `json:"time"`
	Symbol        string `json:"symbol"`
	Side          string `json:"side"`
	OrderType     string `json:"type"`
	PositionSide  string `json:"positionSide"`
	CumQuote      string `json:"cumQuote"`
	Status        string `json:"status"`
	StopPrice     string `json:"stopPrice"`
	Price         string `json:"price"`
	OrigQty       string `json:"origQty"`
	AvgPrice      string `json:"avgPrice"`
	ExecutedQty   string `json:"executedQty"`
	OrderId       int    `json:"orderId"`
	Profit        string `json:"profit"`
	Commission    string `json:"commission"`
	UpdateTime    int    `json:"updateTime"`
	ClientOrderID string `json:"clientOrderID"`
}

// Define response of cancel order request
type CancelOrderResponse struct {
	Success []*Order `json:"success"`
	Failed  []*Order `json:"failed"`
}

func (s *CancelOrderService) Do(ctx context.Context, opts ...RequestOption) (res *CancelOrderResponse, err error) {
	r := &request{method: http.MethodPost, endpoint: "/openApi/swap/v2/trade/order"}

	if s.symbol != "" {
		r.addParam("symbol", s.symbol)
	}

	if s.orderIdList != nil {
		r.addParam("orderIdList", s.orderIdList)
	}

	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}

	resp := new(struct {
		code int
		msg  string
		data *CancelOrderResponse
	})

	err = json.Unmarshal(data, &resp)
	res = resp.data

	if err != nil {
		return nil, err
	}

	return res, nil
}
