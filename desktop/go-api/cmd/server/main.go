package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"foundry-tunnel/internal/app"
)

func main() {
	application, err := app.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := application.StartWebServer(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting web server: %v\n", err)
		os.Exit(1)
	}

	url := application.WebServer.URL()
	fmt.Printf("🎲 Foundry Tunnel Manager Web Server\n")
	fmt.Printf("🌐 Dashboard running at: %s\n", url)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\nShutting down...")
}
