package bingx

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/magicaleks/go-bingx/common"
)

// WsHandler handle raw websocket message
type WsHandler func([]byte)

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

var wsServe = func(initMessage []byte, config *WsConfig, handler WsHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	header := http.Header{}
	header.Add("Accept-Encoding", "gzip")

	c, _, err := websocket.DefaultDialer.Dial(config.Endpoint, header)
	if err != nil {
		return nil, nil, err
	}

	if initMessage != nil {
		err = c.WriteMessage(websocket.TextMessage, initMessage)
		if err != nil {
			return nil, nil, err
		}
	}

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
			if string(decodedMsg) == "Ping" {
				err = c.WriteMessage(websocket.TextMessage, []byte("Pong"))
				if err != nil {
					if !silent {
						errHandler(err)
					}
					return
				}
				continue
			}
			handler(decodedMsg)
		}
	}()
	return
}
