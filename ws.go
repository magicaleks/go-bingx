package bingx

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/magicaleks/go-bingx/common"
)

// WsHandler handle raw websocket message
type WsHandler func(string)

// ErrHandler handles errors
type ErrHandler func(err error)

// WsConfig webservice configuration
type WsConfig struct {
	Endpoint string
}

func newWsConfig(endpoint string) *WsConfig {
	return &WsConfig{
		Endpoint: endpoint,
	}
}

type WsClient struct {
	config *WsConfig
	conn   *websocket.Conn
	stopC  chan struct{}
	doneC  chan struct{}
}

func (wc *WsClient) serve(initMessage string, handler WsHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	header := http.Header{}
	header.Add("Accept-Encoding", "gzip")

	c, _, err := websocket.DefaultDialer.Dial(wc.config.Endpoint, header)
	if err != nil {
		return nil, nil, err
	}

	err = c.WriteMessage(websocket.TextMessage, []byte(initMessage))
	if err != nil {
		return nil, nil, err
	}

	wc.conn = c
	c.SetReadLimit(655350)
	doneC = make(chan struct{})
	stopC = make(chan struct{})
	go func() {
		defer close(doneC)
		silent := false
		go func() {
			select {
			case <-stopC:
				silent = true
			case <-doneC:
			}
			c.Close()
		}()
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				if !silent {
					errHandler(err)
				}
				return
			}
			decodedMsg, err := common.DecodeGzip(message)
			if err != nil {
				if !silent {
					errHandler(err)
				}
				return
			}
			if decodedMsg == "Ping" {
				err = c.WriteMessage(websocket.TextMessage, []byte("Pong"))
				if err != nil {
					if !silent {
						errHandler(err)
					}
					return
				}
			}
			handler(decodedMsg)
		}
	}()
	return
}

func (wc *WsClient) Wait() {
	<-wc.doneC
}

func (wc *WsClient) Close() {
	wc.conn.Close()
	wc.stopC <- struct{}{}
}