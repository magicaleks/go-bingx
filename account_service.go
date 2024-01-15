package bingx

import (
	"context"
	"encoding/json"
	"net/http"
)

type GetBalanceService struct {
	c *Client
}

type Balance struct {
	Asset            string  `json:"asset"`
	Balance          float64 `json:"balance"`
	Equity           float64 `json:"equity"`
	UnrealizedProfit float64 `json:"unrealizedProfit"`
	RealisedProfit   float64 `json:"realisedProfit"`
	AavailableMargin float64 `json:"availableMargin"`
	UsedMargin       float64 `json:"usedMargin"`
	FreezedMargin    float64 `json:"freezedMargin"`
}

func (s *GetBalanceService) Do(ctx context.Context, opts ...RequestOption) (res *Balance, err error) {
	r := &request{method: http.MethodGet, endpoint: "/openApi/swap/v2/user/balance"}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}

	res = new(Balance)
	err = json.Unmarshal(data, &res)

	if err != nil {
		return nil, err
	}

	return res, nil
}
