package main

import (
	"flag"
	"net/url"
	"os"
	"os/signal"
	"time"

	"govern/twitch/twitchmessages"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var addr = flag.String("addr", "localhost:8080", "twitch-cli ws service address")

func startTwitchConnection() {
	flag.Parse()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Info().Msgf("Main:Connection:Connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal().Msgf("Main:Connection:Dial:%v", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Error().Msgf("Main:ReadMessage:Raw:%v", err)
				return
			}
			// May need to add concurrency to the below message parsing/handler in the future.
			converted_message, err := twitchmessages.ConvertToJson(message)
			if err != nil {
				log.Error().Msgf("Main:ConvertToJson:%v", err)
				return
			}
			log.Info().Msgf("Main:ConvertToJson:Converted:%v", converted_message)
			twitchmessages.MessageTypeHandler(converted_message)
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Info().Msg("Main:Connection:Interrupt occured, closing connection.")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Error().Msgf("write close:%v", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func main() {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	// log.Trace().Msg("this is a debug message")
	// log.Debug().Msg("this is a debug message")
	// log.Info().Msg("this is an info message")
	// log.Warn().Msg("this is a warning message")
	// log.Error().Msg("this is an error message")
	// log.Fatal().Msg("this is a fatal message")
	// log.Panic().Msg("This is a panic message")
	startTwitchConnection()
}
