package main

import (
    "context"
    "fmt"
    "os"
    "os/signal"
    "syscall"
)

func main() {
    ctx, cancel := context.WithCancel(context.Background())

    // Start your application by calling the main function from the "app" package.
    if err := publicapi.Run(ctx); err != nil {
        fmt.Printf("Error: %v\n", err)
        cancel()
        os.Exit(1)
    }

    // Handle signals to trigger a graceful shutdown
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
    <-sigCh

    // Perform cleanup and shutdown logic here
    cancel()
    os.Exit(0)
}
