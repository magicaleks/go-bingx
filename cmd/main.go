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

	// res1, err := client.NewOpenOrderService().Symbol("XRP-USDT").Quantity(16).Type(bingx.LimitOrderType).Side(bingx.BuySideType).Price(0.5).Do(context.Background())
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println(res1)

	res2, err := client.NewCancelOrderService().Order(1747931322549219328).Symbol("XRP-USDT").Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res2)

	res3, err := client.NewGetOpenPositionsService().Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res3)

	// var handler = func(event *bingx.WsOrder) {
	// 	fmt.Println(event)
	// }

	// var errHandler = func(err error) {
	// 	fmt.Println(err)
	// }

	// doneC, stopC, _ := bingx.WsOrderUpdateServe(res1, handler, errHandler)

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
