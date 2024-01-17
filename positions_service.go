package bingx

import (
	"context"
	"encoding/json"
	"net/http"
)

type GetOpenPositionsService struct {
	c      *Client
	symbol string
}

func (s *GetOpenPositionsService) Symbol(symbol string) *GetOpenPositionsService {
	s.symbol = symbol
	return s
}

type Position struct {
	Symbol             string `json:"symbol"`
	PositionId         string `json:"positionId"`
	PositionSide       string `json:"positionSide"`
	Isolated           bool   `json:"isolated"`
	PositionAmt        string `json:"positionAmt"`
	AvailableAmt       string `json:"availableAmt"`
	UnrealizedProfit   string `json:"unrealizedProfit"`
	RealisedProfit     string `json:"realisedProfit"`
	InitialMargin      string `json:"initialMargin"`
	AvgPrice           string `json:"avgPrice"`
	LiquidationPrice   string `json:"liquidationPrice"`
	Leverage           int    `json:"leverage"`
	PositionValue      string `json:"positionValue"`
	MarkPrice          string `json:"markPrice"`
	RiskRate           string `json:"riskRate"`
	MaxMarginReduction string `json:"maxMarginReduction"`
	PnlRatio           string `json:"pnlRatio"`
}

func (s *GetOpenPositionsService) Do(ctx context.Context, opts ...RequestOption) (res *[]Position, err error) {
	r := &request{method: http.MethodGet, endpoint: "/openApi/swap/v2/user/positions"}

	if s.symbol != "" {
		r.addParam("symbol", s.symbol)
	}

	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}

	resp := new(struct {
		Code int         `json:"code"`
		Msg  string      `json:"msg"`
		Data *[]Position `json:"data"`
	})

	err = json.Unmarshal(data, &resp)
	res = resp.Data

	if err != nil {
		return nil, err
	}

	return res, nil
}
