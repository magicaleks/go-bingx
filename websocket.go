package bingx

import (
	"encoding/json"
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

type WsClient struct {
	config     *WsConfig
	conn       *websocket.Conn
	stopC      chan struct{}
	doneC      chan struct{}
	subs       map[string]WsHandler
	errHandler ErrHandler
}

func (wc *WsClient) serve() (err error) {
	header := http.Header{}
	header.Add("Accept-Encoding", "gzip")

	c, _, err := websocket.DefaultDialer.Dial(wc.config.Endpoint, header)
	if err != nil {
		return err
	}
	wc.conn = c
	c.SetReadLimit(655350)
	wc.doneC = make(chan struct{})
	wc.stopC = make(chan struct{})
	go func() {
		defer close(wc.doneC)
		silent := false
		go func() {
			select {
			case <-wc.stopC:
				silent = true
			case <-wc.doneC:
			}
			c.Close()
		}()
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				if !silent {
					if wc.errHandler != nil {
						wc.errHandler(err)
					}
				}
				return
			}
			decodedMsg, err := common.DecodeGzip(message)
			if err != nil {
				if !silent {
					if wc.errHandler != nil {
						wc.errHandler(err)
					}
				}
				return
			}
			if string(decodedMsg) == "Ping" {
				err = c.WriteMessage(websocket.TextMessage, []byte("Pong"))
				if err != nil {
					if !silent {
						if wc.errHandler != nil {
							wc.errHandler(err)
						}
					}
					return
				}
			}
			dataTypeStrct := new(struct{ dataType string })
			err = json.Unmarshal(decodedMsg, dataTypeStrct)
			if err != nil {
				if !silent {
					if wc.errHandler != nil {
						wc.errHandler(err)
					}
				}
				return
			}
			wc.subs[dataTypeStrct.dataType](decodedMsg)

		}
	}()
	return
}

// Set error handler
func (wc *WsClient) SetErrHandler(errHandler ErrHandler) {
	wc.errHandler = errHandler
}

// Blocked function
func (wc *WsClient) Wait() {
	<-wc.doneC
}

func (wc *WsClient) Close() {
	wc.stopC <- struct{}{}
}
