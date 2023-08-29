package exchanges

import (
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
)

// Define a custom type for ChannelType.
type ChannelType string

// Define constants for ChannelType.
const (
	OrderBook ChannelType = "orderbook"
	ObDiff    ChannelType = "obdiff"
)

var addr = flag.String("addr", "ws-feed-pro.btcturk.com", "BtcTurk Websocket Feed")

func wsChannel(channel ChannelType, pair string) string {

	u := url.URL{Scheme: "wss", Host: *addr, Path: "/"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Dial Error:", err)
	}
	message := []byte(`[151, {
		"type":151,
		"channel": "` + string(channel) + `",
		"event": ` + pair + `,
		"join":true}]`)

	err = c.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Println("ERROR:", err)
		return
	}
	for {
		messageType, p, err := c.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(string(p))
		if err = c.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}
