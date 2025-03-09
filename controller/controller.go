package controller

import (
	"govern/twitch/twitchclient"
	"govern/twitch/twitchmessages"
	"github.com/rs/zerolog/log"
)

// Channel to trigger reconnection logic
var reconnectChan = make(chan string)

// Controller that manages the WebSocket connection and message handling
func StartController(initialAddr string) {
	msgChannel := make(chan twitchmessages.BaseMessage)

	// Start the WebSocket connection
	go twitchclient.NewTwitchConnection(initialAddr, msgChannel)

	// Handle incoming messages
	go func() {
		for {
			msg := <-msgChannel
			// Pass each message to the message handler
			twitchmessages.MessageTypeHandler(msg, reconnectChan)
		}
	}()

	// Monitor reconnection requests
	for {
		select {
		case newAddr := <-reconnectChan:
			log.Printf("Reconnecting to new address: %s", newAddr)
			// Stop the current WebSocket connection and reconnect
			twitchclient.NewTwitchConnection(newAddr, msgChannel)
		}
	}
}
