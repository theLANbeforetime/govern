package twitchmessages

import (
	"encoding/json"
	"time"

	"github.com/rs/zerolog/log"
)

var CurrentSession PersistantSessionInformation

type MessageMetadata struct {
	MessageId           string    `json:"message_id"`
	MessageType         string    `json:"message_type"`
	MessageTimestamp    time.Time `json:"message_timestamp"`
	SubscriptionType    string    `json:"subscription_type"`
	SubscriptionVersion string    `json:"subscription_version"`
}

type MessageSession struct {
	Id                      string    `json:"id"`
	Status                  string    `json:"status"`
	ConnectedAt             time.Time `json:"connected_at"`
	KeepaliveTimeoutSeconds int       `json:"keepalive_timeout_seconds"`
	ReconnectURL            string    `json:"reconnect_url"`
}

type MessageSubscription struct {
	Id        string `json:"id"`
	Status    string `json:"status"`
	Type      string `json:"type"`
	Version   string `json:"version"`
	Cost      int    `json:"cost"`
	Condition struct {
		BroadcasterUserId string `json:"broadcaster_user_id"`
	} `json:"condition"`
	Transport struct {
		Method    string `json:"method"`
		SessionId string `json:"session_id"`
	} `json:"transport"`
	CreatedAt time.Time `json:"created_at"`
}

type MessageEvent struct {
	UserId               int       `json:"user_id"`
	UserLogin            string    `json:"user_login"`
	UserName             string    `json:"user_name"`
	BroadcasterUserId    string    `json:"broadcaster_user_id"`
	BroadcasterUserLogin string    `json:"broadcaster_user_login"`
	BroadcasterUserName  string    `json:"broadcaster_user_name"`
	FollowedAt           time.Time `json:"followed_at"`
}

type MessagePayload struct {
	Session      MessageSession
	Subscription MessageSubscription
	Event        MessageEvent
}

type BaseMessage struct {
	Metadata MessageMetadata
	Payload  MessagePayload
}

type PersistantSessionInformation struct {
	TimeOut      int
	ReconnectURL string
}

type DiscordNotification struct {
	Type            string
	BroadcasterName string
}

func ConvertToJson(message []byte) (BaseMessage, error) {
	var jsonData map[string]interface{}
	err := json.Unmarshal(message, &jsonData)
	if err != nil {
		log.Error().Msgf("Messages:ConvertToJson:Unmarshal:Byte:%v", err)
		return BaseMessage{}, err
	}

	jsonStr, err := json.Marshal(jsonData)
	if err != nil {
		log.Error().Msgf("Messages:ConvertToJson:Marshal:%v", err)
		return BaseMessage{}, err
	}
	var convertedMessage BaseMessage
	if err := json.Unmarshal(jsonStr, &convertedMessage); err != nil {
		log.Error().Msgf("Messages:ConvertToJson:Unmarshal:Json:%v", err)
		return BaseMessage{}, err
	}
	return convertedMessage, nil
}

// https://dev.twitch.tv/docs/eventsub/handling-websocket-events/
func MessageTypeHandler(message BaseMessage, reconnectChan chan string) {
	switch message.Metadata.MessageType {
	case "session_welcome":
		sessionWelcome(message)
		log.Info().Msgf("Messages:SessionWelcome: ReconnectUrl: %s, TimeOut: %v", CurrentSession.ReconnectURL, CurrentSession.TimeOut)
	case "session_keepalive":
		log.Info().Msgf("Messages:SessionKeepAlive: Connection good.")
	case "notification":
		notification := notification(message)
		log.Info().Msgf("Messages:Notification: Received notification of type: %s, for broadcaster %s", notification.Type, notification.BroadcasterName)
	case "session_reconnect":
		reconnect := sessionReconnect(message)
		if reconnect != "" {
			log.Info().Msgf("Recevied reconnect message for address: %v", reconnect)
		}
		reconnectChan <- reconnect
	case "revocation":

	}
}

func sessionWelcome(message BaseMessage) PersistantSessionInformation {
	CurrentSession.ReconnectURL = message.Payload.Session.ReconnectURL
	CurrentSession.TimeOut = message.Payload.Session.KeepaliveTimeoutSeconds
	return CurrentSession
}

func notification(message BaseMessage) DiscordNotification {
	notification := DiscordNotification{
		Type:            message.Payload.Subscription.Type,
		BroadcasterName: message.Payload.Event.BroadcasterUserName,
	}
	return notification
}

func sessionReconnect(message BaseMessage) string {
	reconnect_address := message.Payload.Session.ReconnectURL
	if reconnect_address != "" {
		log.Info().Msgf("Received reconnect message, attempting to reconnect to: %s", reconnect_address)
		return reconnect_address
	}
	return ""
}

func revocation(message BaseMessage) {

}
