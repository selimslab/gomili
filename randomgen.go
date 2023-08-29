package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow any origin for WebSocket connections
		},
	}
)

func randomNumbersGenerator() <-chan int {
	rand.Seed(time.Now().UnixNano())
	ch := make(chan int)

	go func() {
		defer close(ch)
		for {
			randomNumber := rand.Intn(100) // Generate random number between 0 and 99
			ch <- randomNumber
			time.Sleep(1 * time.Second) // Send a new number every second
		}
	}()

	return ch
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	numbersCh := randomNumbersGenerator()
	for num := range numbersCh {
		err := conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Random Number: %d", num)))
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "public/index.html")
}

func main() {
	http.HandleFunc("/ws", websocketHandler)
	http.HandleFunc("/", homeHandler)

	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
