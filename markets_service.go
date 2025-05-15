package bingx

import (
	"context"
	"encoding/json"
	"net/http"
)

type GetSymbolDataService struct {
	c      *Client
	symbol string
}

func (s *GetSymbolDataService) Symbol(symbol string) *GetSymbolDataService {
	s.symbol = symbol
	return s
}

type SymbolData struct {
	Symbol            string  `json:"symbol"`
	QuantityPrecision int     `json:"quantityPrecision"`
	PricePrecision    int     `json:"pricePrecision"`
	TradeMinQuantity  float64 `json:"tradeMinQuantity"`
}

func (s *GetSymbolDataService) Do(ctx context.Context, opts ...RequestOption) (res *SymbolData, err error) {
	r := &request{method: http.MethodGet, endpoint: "/openApi/swap/v2/quote/contracts"}

	if s.symbol != "" {
		r.addParam("symbol", s.symbol)
	}

	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}

	resp := new(struct {
		Code int          `json:"code"`
		Msg  string       `json:"msg"`
		Data []SymbolData `json:"data"`
	})

	err = json.Unmarshal(data, &resp)
	res = &resp.Data[0]

	if err != nil {
		return nil, err
	}

	return res, nil
}
