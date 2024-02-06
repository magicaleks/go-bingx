package main

import (
	"context"
	"fmt"

	"github.com/magicaleks/go-bingx"
)

const (
	APIkey    = "kAETP0GtSvvjZjJCzgjyoo6kaOAZQsFurlppV2cyQWUK1ETs55QSkf9CNwtBf0tHAo5jiNZzhWmUo02g9w"
	SecretKey = "4S3tymbzGKUa1OT9n7sI83El5O0xGHwNMNSOlMtMGuowjVJkY983ZtlZ7Qf04I1lM39l74cD7MDIYk6bLyg"
)

func main() {
	client := bingx.NewClient(APIkey, SecretKey)

	res1, err := client.NewGetKlinesService().Symbol("XRP-USDT").Interval(bingx.Interval1).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res1[0])

	// res1, err := client.NewCreateOrderService().Symbol("XRP-USDT").Quantity(16).Type(bingx.LimitOrderType).Side(bingx.BuySideType).Price(0.5).Do(context.Background())
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println(res1)

	// res6, err := client.NewGetOrderService().Symbol("XRP-USDT").OrderId(res1.OrderId).Do(context.Background())
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println(res6)

	// res3, err := client.NewGetOpenPositionsService().Do(context.Background())
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println(res3)

	// res2, err := client.NewCancelOrderService().OrderId(res1.OrderId).Symbol("XRP-USDT").Do(context.Background())
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println(res2)

	// res4, err := client.NewGetBalanceService().Do(context.Background())
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println(res4)

	// res5, err := client.NewGetAccountListenKeyService().Do(context.Background())
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println(res5)

	// var handler = func(event *bingx.WsOrder) {
	// 	fmt.Println(event)
	// }

	// var errHandler = func(err error) {
	// 	fmt.Println(err)
	// }

	// doneC, stopC, _ := bingx.WsOrderUpdateServe(res5, handler, errHandler)

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// go func() {
	// 	time.Sleep(20 * time.Second)
	// 	stopC <- struct{}{}
	// }()

	// <-doneC
}
