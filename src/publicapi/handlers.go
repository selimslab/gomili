package publicapi

import (
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
)

func startServer() {
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	addClient(conn)

	wg.Add(1)
	defer func() {
		removeClient(conn)
		wg.Done()
	}()

	go func() {
		for update := range klineUpdates {
			err := conn.WriteMessage(websocket.TextMessage, []byte(update))
			if err != nil {
				fmt.Println("Error sending message to client:", err)
				break
			}
		}
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}
