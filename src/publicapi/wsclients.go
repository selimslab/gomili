package publicapi

import (
	"fmt"
	"sync"
	"github.com/gorilla/websocket"
)

func sendToClients(message string) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			fmt.Println("Error sending message to client:", err)
			removeClient(client)
		}
	}
}

func closeClients() {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	for client := range clients {
		removeClient(client)
	}
}

func waitAndCloseClients() {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	var wg sync.WaitGroup
	for client := range clients {
		wg.Add(1)
		go func(c *websocket.Conn) {
			defer wg.Done()
			c.Close()
		}(client)
	}
	wg.Wait()
}

func addClient(client *websocket.Conn) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	clients[client] = true
}

func removeClient(client *websocket.Conn) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	client.Close()
	delete(clients, client)
}
