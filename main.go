package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/tgdrive/teldrive/cmd"
)

func main() {
	// Setup a context that listens for the interrupt signal
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Run your command with context
	go func() {
		if err := cmd.New().ExecuteContext(ctx); err != nil {
			panic(err)
		}
	}()

	// Start the HTTP server to handle requests (needed for Vercel)
	http.HandleFunc("/", handler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		ReadHeaderTimeout: 5 * time.Second,
	}

	fmt.Printf("Server started at port %s\n", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("Could not listen on port %s: %v\n", port, err)
	}
}

// Handler is the function Vercel looks for to process HTTP requests
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from Vercel with Go!")
}
