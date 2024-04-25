package main

import (
	"context"
	"log"

	"github.com/magicaleks/go-bingx"
)

const (
	APIkey    = "API_KEY"
	SecretKey = "SECRET_KEY"
)

func main() {
	client := bingx.NewClient(APIkey, SecretKey)

	// perform account subscription
	listenKey, err := client.NewGetAccountListenKeyService().
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Account subscription listen key: %s", listenKey)

	doneC, _, err := bingx.WsOrderUpdateServe(listenKey, func(order *bingx.WsOrder) {
		log.Printf("WsOrderUpdateServe update: %+v", order)
	}, func(err error) {
		log.Printf("WsOrderUpdateServe error: %s\n", err)
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Account subscription has been established")

	// create limit orders
	symbol := "LINK-USDT"
	order, err := client.NewCreateOrderService().
		Symbol(symbol).
		Quantity(0.8).
		Type(bingx.LimitOrderType).
		Side(bingx.BuySideType).
		Price(15.8).
		ClientOrderID("my-order-id").
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Limit order created: %+v", order)

	// cancel it
	cancelResponse, err := client.NewCancelOrderService().
		Symbol(symbol).
		ClientOrderId("my-order-id").
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Limit order canceled: %+v", cancelResponse)

	<-doneC
}
