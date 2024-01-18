package bingx

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

const (
	baseWsUrl        = "wss://open-api-swap.bingx.com/swap-market"
	baseAccountWsUrl = "wss://open-api-swap.bingx.com/swap-market?listenKey="
)

func getWsEndpoint() string {
	return baseWsUrl
}

// func getAccountWsEndpoint(listenKey string) string {
// 	return baseAccountWsUrl + listenKey
// }

type Event struct {
	Code     int    `json:"code"`
	DataType string `json:"dataType"`
	Data     string `json:"data"`
}

type RequestType string

const (
	SubscribeRequestType  RequestType = "sub"
	UnubscribeRequestType RequestType = "unsub"
)

type RequestEvent struct {
	Id       uuid.UUID   `json:"id"`
	ReqType  RequestType `json:"reqType"`
	DataType string      `json:"dataType"`
}

type KLineEvent struct {
	Data   interface{} `json:"data"`
	Asks   interface{} `json:"asks"`
	Bids   interface{} `json:"bids"`
	Price  float64     `json:"p"`
	Volume float64     `json:"v"`
}

type WsKLineHandler func(*KLineEvent)

func WsKLineServe(symbol string, interval string, handler WsKLineHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	// Symbol e.g. "BTC-USDT"
	// Interval e.g. "1m", "3h"
	reqEvent := RequestEvent{
		Id:       uuid.New(),
		ReqType:  SubscribeRequestType,
		DataType: fmt.Sprintf("%s@kline_%s", symbol, interval),
	}

	var wsHandler = func(data []byte) {

		ev := new(Event)
		err := json.Unmarshal(data, ev)
		if err != nil {
			errHandler(err)
			return
		}

		if ev.DataType == reqEvent.DataType {
			event := new(KLineEvent)
			err := json.Unmarshal(data, event)
			if err != nil {
				errHandler(err)
				return
			}

			handler(event)
		}

	}

	initMessage, err := json.Marshal(reqEvent)
	if err != nil {
		return nil, nil, err
	}

	return wsServe(initMessage, newWsConfig(getWsEndpoint()), wsHandler, errHandler)

}
