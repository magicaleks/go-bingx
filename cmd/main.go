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

	res1, err := client.NewGetBalanceService().Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res1.Balance)

	res2, err := client.NewGetOpenPositionsService().Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, position := range *res2 {
		fmt.Println(position.Symbol)
		fmt.Println(position.LiquidationPrice)
		fmt.Println(position.InitialMargin)
	}
}
