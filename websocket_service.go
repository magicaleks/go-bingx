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

func getAccountWsEndpoint(listenKey string) string {
	return baseAccountWsUrl + listenKey
}

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

type WsKLineEvent struct {
	Symbol string  `json:"s"`
	Open   float64 `json:"o"`
	Close  float64 `json:"c"`
	High   float64 `json:"h"`
	Low    float64 `json:"l"`
	Volume float64 `json:"v"`
}

type WsKLineHandler func(*WsKLineEvent)

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

			event := &WsKLineEvent{
				Symbol: _eventData.Symbol,
				Open:   o,
				Close:  c,
				High:   h,
				Low:    l,
				Volume: v,
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

type WsOrder struct {
	Symbol        string        `json:"s"`
	Side          SideType      `json:"S"`
	OrderType     OrderType     `json:"o"`
	Price         string        `json:"p"`
	AveragePrice  string        `json:"ap"`
	Quantity      string        `json:"q"`
	StopPrice     string        `json:"sp"`
	Status        OrderStatus   `json:"X"`
	Spec          OrderSpecType `json:"x"`
	Timestamp     int           `json:"T"`
	OrderId       int           `json:"i"`
	ClientOrderID string        `json:"c"`
}

type WsOrderUpdateEvent struct {
	EventType string   `json:"e"`
	Time      int      `json:"E"`
	Order     *WsOrder `json:"o"`
}

type WsOrderUpdateHandler func(*WsOrder)

func WsOrderUpdateServe(listenKey string, handler WsOrderUpdateHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	var wsHandler = func(data []byte) {

		var evMap map[string]interface{}
		err := json.Unmarshal(data, &evMap)
		if err != nil {
			errHandler(err)
			return
		}

		if evMap["e"].(string) == "ORDER_TRADE_UPDATE" {
			event := new(WsOrderUpdateEvent)
			err = json.Unmarshal(data, event)
			if err != nil {
				errHandler(err)
				return
			}
			handler(event.Order)
		}

	}

	return wsServe(nil, newWsConfig(getAccountWsEndpoint(listenKey)), wsHandler, errHandler)
}
