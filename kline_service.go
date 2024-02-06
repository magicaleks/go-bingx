package bingx

import (
	"context"
	"encoding/json"
	"net/http"
)

type GetKlinesService struct {
	c         *Client
	symbol    string
	interval  Interval
	startTime int64
	endTime   int64
	limit     int64
}

// Define Kline model
type Kline struct {
	Open   string `json:"open"`
	Close  string `json:"close"`
	High   string `json:"hign"`
	Low    string `json:"low"`
	Volume string `json:"volume"`
	Time   int64  `json:"time"`
}

func (s *GetKlinesService) Symbol(symbol string) *GetKlinesService {
	s.symbol = symbol
	return s
}

func (s *GetKlinesService) Interval(interval Interval) *GetKlinesService {
	s.interval = interval
	return s
}

func (s *GetKlinesService) StartTime(startTime int64) *GetKlinesService {
	s.startTime = startTime
	return s
}

func (s *GetKlinesService) EndTime(endTime int64) *GetKlinesService {
	s.endTime = endTime
	return s
}

func (s *GetKlinesService) Limit(limit int64) *GetKlinesService {
	s.limit = limit
	return s
}

func (s *GetKlinesService) Do(ctx context.Context, opts ...RequestOption) (res []*Kline, err error) {
	r := &request{method: http.MethodGet, endpoint: "/openApi/swap/v3/quote/klines"}

	if s.symbol != "" {
		r.addParam("symbol", s.symbol)
	}

	if s.interval != "" {
		r.addParam("interval", s.interval)
	}

	if s.startTime != 0 {
		r.addParam("startTime", s.startTime)
	}

	if s.endTime != 0 {
		r.addParam("endTime", s.endTime)
	}

	if s.limit != 0 {
		r.addParam("limit", s.limit)
	}

	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}

	resp := new(struct {
		Code int      `json:"code"`
		Msg  string   `json:"msg"`
		Data []*Kline `json:"data"`
	})

	err = json.Unmarshal(data, &resp)

	if err != nil {
		return nil, err
	}

	res = resp.Data

	return res, nil
}
