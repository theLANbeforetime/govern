package main

import (
	"govern/broker/messagebroker"
	"govern/twitch/twitchclient"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	broker := messagebroker.NewBroker()

	// log.Trace().Msg("this is a debug message")
	// log.Debug().Msg("this is a debug message")
	// log.Info().Msg("this is an info message")
	// log.Warn().Msg("this is a warning message")
	// log.Error().Msg("this is an error message")
	// log.Fatal().Msg("this is a fatal message")
	// log.Panic().Msg("This is a panic message")
	// startTwitchConnection()

	//Entire program hinges on connection to twitch.
	//If connection to twitch goes down rest of program should go down.
	for shouldRestartTwitchConnection := true; shouldRestartTwitchConnection; {
		discord_subscriber := broker.Subscribe("live_notifications")
		go func() {
			for {
				select {
				case msg, ok := <-discord_subscriber.Channel:
					if !ok {
						log.Info().Msg("Subscriber channel closed.")
						return
					}
					log.Info().Msgf("Received: %v\n", msg)
				case <-discord_subscriber.Unsubscribe:
					log.Info().Msg("Unsubscribed.")
					return
				}
			}
		}()
		shouldRestartTwitchConnection = twitchclient.StartTwitchConnection(broker)
	}
}
