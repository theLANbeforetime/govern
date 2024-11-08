package main

import (
	"govern/twitch/twitchclient"

	"github.com/rs/zerolog"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	// log.Trace().Msg("this is a debug message")
	// log.Debug().Msg("this is a debug message")
	// log.Info().Msg("this is an info message")
	// log.Warn().Msg("this is a warning message")
	// log.Error().Msg("this is an error message")
	// log.Fatal().Msg("this is a fatal message")
	// log.Panic().Msg("This is a panic message")
	// startTwitchConnection()
	for shouldRestartTwitchConnection := true; shouldRestartTwitchConnection; {
		shouldRestartTwitchConnection = twitchclient.StartTwitchConnection()
	}
}
