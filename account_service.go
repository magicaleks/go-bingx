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
	UserId           string `json:"userId"`
	Asset            string `json:"asset"`
	Balance          string `json:"balance"`
	Equity           string `json:"equity"`
	UnrealizedProfit string `json:"unrealizedProfit"`
	RealisedProfit   string `json:"realisedProfit"`
	AavailableMargin string `json:"availableMargin"`
	UsedMargin       string `json:"usedMargin"`
	FreezedMargin    string `json:"freezedMargin"`
}

func (s *GetBalanceService) Do(ctx context.Context, opts ...RequestOption) (res *Balance, err error) {
	r := &request{method: http.MethodGet, endpoint: "/openApi/swap/v2/user/balance"}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}

	resp := new(struct {
		Code int                 `json:"code"`
		Msg  string              `json:"msg"`
		Data map[string]*Balance `json:"data"`
	})

	err = json.Unmarshal(data, resp)
	res = resp.Data["balance"]

	if err != nil {
		return nil, err
	}

	return res, nil
}

type GetAccountListenKeyService struct {
	c *Client
}

func (s *GetAccountListenKeyService) Do(ctx context.Context, opts ...RequestOption) (res string, err error) {
	r := &request{method: http.MethodPost, endpoint: "/openApi/user/auth/userDataStream"}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return "", err
	}

	resp := new(struct {
		ListenKey string `json:"listenKey"`
	})

	err = json.Unmarshal(data, resp)
	if err != nil {
		return "", err
	}

	return resp.ListenKey, nil
}
