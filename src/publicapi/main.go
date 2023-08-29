package publicapi

import (
    "context"
    "fmt"
    "net/http"
)

// Run starts the application.
func Run(ctx context.Context) error {
    http.HandleFunc("/", indexHandler)
    http.HandleFunc("/ws", wsHandler)

    // Serve the client HTML file
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

    go startServer()

    generator := newKlineGenerator(ctx)

    fmt.Println("Web server started. Open http://localhost:8080/static/index.html in your browser to view real-time Kline updates.")

    for klineEvent := range generator {
        data := processKlineEvent(klineEvent)
        sendToClients(data)
    }

    return nil
}

// Add your indexHandler, wsHandler, and other functions here.
