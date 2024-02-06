package bingx

import (
	"context"
	"encoding/json"
	"net/http"
)

// Only Market and Limit orders supported
type CreateOrderService struct {
	c             *Client
	symbol        string
	orderType     OrderType
	side          SideType
	positionSide  PositionSideType
	clientOrderID string
	reduceOnly    string
	price         float64
	quantity      float64
}

func (s *CreateOrderService) Symbol(symbol string) *CreateOrderService {
	s.symbol = symbol
	return s
}

func (s *CreateOrderService) Type(orderType OrderType) *CreateOrderService {
	s.orderType = orderType
	return s
}

func (s *CreateOrderService) Side(side SideType) *CreateOrderService {
	s.side = side
	return s
}

func (s *CreateOrderService) PositionSide(side PositionSideType) *CreateOrderService {
	s.positionSide = side
	return s
}

func (s *CreateOrderService) ClientOrderID(clientOrderID string) *CreateOrderService {
	s.clientOrderID = clientOrderID
	return s
}

func (s *CreateOrderService) ReduceOnly() *CreateOrderService {
	s.reduceOnly = "true"
	return s
}

func (s *CreateOrderService) Price(price float64) *CreateOrderService {
	s.price = price
	return s
}

func (s *CreateOrderService) Quantity(quantity float64) *CreateOrderService {
	s.quantity = quantity
	return s
}

type CreateOrderResponse struct {
	OrderId int `json:"orderId"`
}

func (s *CreateOrderService) Do(ctx context.Context, opts ...RequestOption) (res *CreateOrderResponse, err error) {
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

	if s.clientOrderID != "" {
		r.addParam("clientOrderID", s.clientOrderID)
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

	//fmt.Print(string(data))

	resp := new(struct {
		Code int                             `json:"code"`
		Msg  string                          `json:"msg"`
		Data map[string]*CreateOrderResponse `json:"data"`
	})

	err = json.Unmarshal(data, &resp)
	res = resp.Data["order"]

	if err != nil {
		return nil, err
	}

	return res, nil
}

type CancelOrderService struct {
	c             *Client
	symbol        string
	orderId       int
	clientOrderID string
}

func (s *CancelOrderService) Symbol(symbol string) *CancelOrderService {
	s.symbol = symbol
	return s
}

func (s *CancelOrderService) OrderId(orderId int) *CancelOrderService {
	s.orderId = orderId

	return s
}

func (s *CancelOrderService) ClientOrderId(clientOrderID string) *CancelOrderService {
	s.clientOrderID = clientOrderID

	return s
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

	if s.clientOrderID != "" {
		r.addParam("clientOrderID", s.clientOrderID)
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

	if err != nil {
		return nil, err
	}

	res = resp.Data["order"]

	return res, nil
}

type GetOrderService struct {
	c             *Client
	symbol        string
	orderId       int
	clientOrderID string
}

// Define response of get order request
type GetOrderResponse struct {
	Time          int64            `json:"time"`
	Symbol        string           `json:"symbol"`
	Side          SideType         `json:"side"`
	OrderType     OrderType        `json:"type"`
	PositionSide  PositionSideType `json:"positionSide"`
	ReduceOnly    bool             `json:"reduceOnly"`
	CumQuote      string           `json:"cumQuote"`
	Status        OrderStatus      `json:"status"`
	StopPrice     string           `json:"stopPrice"`
	Price         string           `json:"price"`
	OrigQuantity  string           `json:"origQty"`
	AveragePrice  string           `json:"avgPrice"`
	Quantity      string           `json:"executedQty"`
	OrderId       int              `json:"orderId"`
	Profit        string           `json:"profit"`
	Fee           string           `json:"commission"`
	UpdateTime    int64            `json:"ppdateTime"`
	WorkingType   OrderWorkingType `json:"workingType"`
	ClientOrderID string           `json:"clientOrderID"`
}

func (s *GetOrderService) Symbol(symbol string) *GetOrderService {
	s.symbol = symbol
	return s
}

func (s *GetOrderService) OrderId(orderId int) *GetOrderService {
	s.orderId = orderId
	return s
}

func (s *GetOrderService) ClientOrderId(clientOrderID string) *GetOrderService {
	s.clientOrderID = clientOrderID
	return s
}

func (s *GetOrderService) Do(ctx context.Context, opts ...RequestOption) (res *GetOrderResponse, err error) {
	r := &request{method: http.MethodGet, endpoint: "/openApi/swap/v2/trade/order"}

	if s.symbol != "" {
		r.addParam("symbol", s.symbol)
	}

	if s.orderId != 0 {
		r.addParam("orderId", s.orderId)
	}

	if s.clientOrderID != "" {
		r.addParam("clientOrderID", s.clientOrderID)
	}

	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}

	resp := new(struct {
		Code int                          `json:"code"`
		Msg  string                       `json:"msg"`
		Data map[string]*GetOrderResponse `json:"data"`
	})

	err = json.Unmarshal(data, &resp)

	if err != nil {
		return nil, err
	}

	res = resp.Data["order"]

	return res, nil
}
