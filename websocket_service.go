package bingx

import (
	"encoding/json"
	"fmt"
	"strconv"

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
	Code     int         `json:"code"`
	DataType string      `json:"dataType"`
	Data     interface{} `json:"data"`
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
	C      float64 `json:"c"`
	H      float64 `json:"h"`
	L      float64 `json:"l"`
	O      float64 `json:"o"`
	V      float64 `json:"v"`
	Symbol string  `json:"s"`
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

		fmt.Println("Low level handler: ", ev)

		if ev.DataType == reqEvent.DataType {
			_eventData := new(struct {
				Symbol string                   `json:"s"`
				Data   []map[string]interface{} `json:"data"`
			})
			err := json.Unmarshal(data, _eventData)
			if err != nil {
				errHandler(err)
				return
			}

			c, _ := strconv.ParseFloat(_eventData.Data[0]["c"].(string), 64)
			h, _ := strconv.ParseFloat(_eventData.Data[0]["h"].(string), 64)
			l, _ := strconv.ParseFloat(_eventData.Data[0]["l"].(string), 64)
			o, _ := strconv.ParseFloat(_eventData.Data[0]["o"].(string), 64)
			v, _ := strconv.ParseFloat(_eventData.Data[0]["v"].(string), 64)

			event := &KLineEvent{
				Symbol: _eventData.Symbol,
				C:      c,
				H:      h,
				L:      l,
				O:      o,
				V:      v,
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
