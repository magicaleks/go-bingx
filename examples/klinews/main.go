package main

import (
	"log"

	"github.com/magicaleks/go-bingx"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	symbol := "LINK-USDT"
	doneC, _, err := bingx.WsKLineServe(symbol, "1m", func(event *bingx.WsKLineEvent) {
		log.Printf("%s price update: %+v", symbol, event)
	}, func(err error) {
		log.Printf("WsKLine error: %s\n", err)
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	<-doneC
}
