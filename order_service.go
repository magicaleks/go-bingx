package bingx

import (
	"context"
	"encoding/json"
	"net/http"
)

// Only Market and Limit orders supported
type OpenOrderService struct {
	c            *Client
	symbol       string
	orderType    OrderType
	side         SideType
	positionSide PositionSideType
	reduceOnly   string
	price        float64
	quantity     float64
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

func (s *OpenOrderService) PositionSide(side PositionSideType) *OpenOrderService {
	s.positionSide = side
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

	if s.positionSide != "" {
		r.addParam("positionSide", s.positionSide)
	} else {
		r.addParam("positionSide", BothPositionSideType)
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
		Code int                           `json:"code"`
		Msg  string                        `json:"msg"`
		Data map[string]*OpenOrderResponse `json:"data"`
	})

	err = json.Unmarshal(data, &resp)
	res = resp.Data["order"]

	if err != nil {
		return nil, err
	}

	return res, nil
}

type CancelOrderService struct {
	c       *Client
	symbol  string
	orderId int
}

func (s *CancelOrderService) Symbol(symbol string) *CancelOrderService {
	s.symbol = symbol
	return s
}

func (s *CancelOrderService) Order(orderId int) *CancelOrderService {
	s.orderId = orderId

	return s
}

// Define order
type Order struct {
	Symbol        string      `json:"symbol"`
	Side          SideType    `json:"side"`
	OrderType     OrderType   `json:"type"`
	Status        OrderStatus `json:"status"`
	ReduceOnly    bool        `json:"reduceOnly"`
	Price         string      `json:"price"`
	Quantity      string      `json:"quantity"`
	StopPrice     string      `json:"stopPrice"`
	PriceRate     string      `json:"priceRate"`
	StopLoss      string      `json:"stopLoss"`   // mb later
	TakeProfit    string      `json:"takeProfit"` // mb later
	WorkingType   string      `json:"workingType"`
	Timestamp     string      `json:"timestamp"`
	OrderId       int         `json:"orderId"`
	ClientOrderID string      `json:"clientOrderID"`
}

// Define response of cancel order request
type CancelOrderResponse struct {
	Time          int              `json:"time"`
	Symbol        string           `json:"symbol"`
	Side          SideType         `json:"side"`
	OrderType     OrderType        `json:"type"`
	PositionSide  PositionSideType `json:"positionSide"`
	CumQuote      string           `json:"cumQuote"`
	Status        OrderStatus      `json:"status"`
	StopPrice     string           `json:"stopPrice"`
	Price         string           `json:"price"`
	OrigQty       string           `json:"origQty"`
	AvgPrice      string           `json:"avgPrice"`
	ExecutedQty   string           `json:"executedQty"`
	OrderId       int              `json:"orderId"`
	Profit        string           `json:"profit"`
	Commission    string           `json:"commission"`
	UpdateTime    int              `json:"updateTime"`
	ClientOrderID string           `json:"clientOrderID"`
}

func (s *CancelOrderService) Do(ctx context.Context, opts ...RequestOption) (res *CancelOrderResponse, err error) {
	r := &request{method: http.MethodDelete, endpoint: "/openApi/swap/v2/trade/order"}

	if s.symbol != "" {
		r.addParam("symbol", s.symbol)
	}

	if s.orderId != 0 {
		r.addParam("orderId", s.orderId)
	}

	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}

	resp := new(struct {
		Code int                             `json:"code"`
		Msg  string                          `json:"msg"`
		Data map[string]*CancelOrderResponse `json:"data"`
	})

	err = json.Unmarshal(data, &resp)
	res = resp.Data["order"]

	if err != nil {
		return nil, err
	}

	return res, nil
}
