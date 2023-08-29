package exchanges

import (
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2"
)

func newKlineGenerator(ctx context.Context) <-chan *binance.WsKlineEvent {
	c := make(chan *binance.WsKlineEvent)

	wsKlineHandler := func(event *binance.WsKlineEvent) {
		select {
		case c <- event:
		case <-ctx.Done():
		}
	}

	errHandler := func(err error) {
		fmt.Println(err)
	}

	go func() {
		doneC, _, err := binance.WsKlineServe("BTCUSDT", "1m", wsKlineHandler, errHandler)
		if err != nil {
			fmt.Println(err)
			close(c)
			return
		}

		select {
		case <-ctx.Done():
			close(c)
		case <-doneC:
			close(c)
		}
	}()

	return c
}
