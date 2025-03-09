package main

import (
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"govern/controller"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Default addr to testing if flag is not passed in via CLI.
var addr = flag.String("addr", "ws://localhost:8080/ws", "twitch-cli ws service address")

func main() {
	flag.Parse()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	var wg sync.WaitGroup

	// Channel to listen for system signals (e.g., Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	wg.Add(1)

	go func() {
		controller.StartController(*addr)
	}()

	wg.Done()

	// Wait for an interrupt signal to initiate graceful shutdown
	select {
	case <-sigChan:
		// Handle shutdown signal (Ctrl+C or SIGTERM)
		log.Info().Msgf("Received shutdown signal. Shutting down gracefully...")
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Final cleanup before exiting
	log.Info().Msgf("Application shutdown complete.")
}
