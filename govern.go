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

func startTwitchConnection() bool {
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
	timeout := 15 * time.Second //Likely need to adjust this based on session_start field.
	// Create a channel to receive data
	timerCh := make(chan bool)

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			timerCh <- true
			if err != nil {
				log.Error().Msgf("Main:ReadMessage:Error:%v", err)
				return
			}
			log.Trace().Msgf("Main:ReadMessage:RawMessage:%v", message)
			// May need to add concurrency to the below message parsing/handler in the future.
			converted_message, err := twitchmessages.ConvertToJson(message)
			if err != nil {
				// If we fail to parse message correctly we just log an error but processing continues.
				log.Error().Msgf("Main:ConvertToJson:%v", err)
			}
			log.Info().Msgf("Main:ConvertToJson:Converted:%v", converted_message)
			twitchmessages.MessageTypeHandler(converted_message)
		}
	}()

	for {
		select {
		case <-timerCh:
			log.Info().Msgf("Main:Connection:Timeout: Message recevied within set time-out: %v", timeout)
		case <-time.After(timeout):
			log.Info().Msgf("Main:Connection:Timeout: No messaged received within set time-out: %v", timeout)
			return true
		case <-done:
			return false
		case <-interrupt:
			log.Info().Msg("Main:Connection:Interrupt occured, closing connection.")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Error().Msgf("write close:%v", err)
				return false
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return false
		}
	}
}

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
		shouldRestartTwitchConnection = startTwitchConnection()
	}
}
