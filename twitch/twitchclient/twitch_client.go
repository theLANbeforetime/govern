package twitchclient

import (
	"govern/twitch/twitchmessages"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

func NewTwitchConnection(addr string, msgChannel chan twitchmessages.BaseMessage) {
	for {
		u := addr
		c, _, err := websocket.DefaultDialer.Dial(u, nil)
		if err != nil {
			log.Fatal().Msgf("twitchclient:newtwitchconnection:dial:%v", err)
		}
		defer c.Close()
		log.Info().Msgf("twitchclient:newtwitchconnection: Successfully connected to %s", u)

		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Error().Msgf("twitchclient:newtwitchconnection:read: Client disconnected due to %v:", err)
				c.Close()
				break
			}
			log.Info().Msgf("twitchclient:newtwitchconnection:read: Successfully received a message. Sending message to handler.")
			converted_message, err := twitchmessages.ConvertToJson(msg)
			if err != nil {
				log.Error().Msgf("twitchclient:newtwitchconnection:convert: Failed to convert message to JSON due to: %v", err)
			}
			msgChannel <- converted_message
		}
	}
}
