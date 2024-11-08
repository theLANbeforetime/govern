package twitchmessages

import (
	"testing"
	"time"
)

// Start websocket server locally with below command:
// twitch event websocket start-server
// Websocket server should start on: ws://127.0.0.1:8080/ws

//twitch event trigger stream.online --transport=websocket

func TestMessageStruct(t *testing.T) {
	t.Run("checking welcome message", func(t *testing.T) {
		welcome := BaseMessage{
			Metadata: MessageMetadata{
				MessageId:        "123",
				MessageType:      "session_welcome",
				MessageTimestamp: time.Now(),
			},
		}

		got := welcome.Metadata.MessageType
		want := "session_welcome"
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})
}

func TestMessageTypeHandler(t *testing.T) {
	t.Run("checking session_welcome message timeout", func(t *testing.T) {
		welcome := BaseMessage{
			Metadata: MessageMetadata{
				MessageId:        "123",
				MessageType:      "session_welcome",
				MessageTimestamp: time.Now(),
			},
			Payload: MessagePayload{
				Session: MessageSession{
					Id:                      "AQoQexAWVYKSTIu4ec_2VAxyuhAB",
					Status:                  "connected",
					ConnectedAt:             time.Now(),
					KeepaliveTimeoutSeconds: 40,
					ReconnectURL:            "wss://eventsub.wss.twitch.tv?XYZ",
				},
			},
		}
		timeout := 40
		got := sessionWelcome(welcome).TimeOut
		want := timeout
		if got != want {
			t.Errorf("got %d want %d", got, want)
		}
	})

	t.Run("checking session_welcome message reconnect url", func(t *testing.T) {
		welcome := BaseMessage{
			Metadata: MessageMetadata{
				MessageId:        "123",
				MessageType:      "session_welcome",
				MessageTimestamp: time.Now(),
			},
			Payload: MessagePayload{
				Session: MessageSession{
					Id:                      "AQoQexAWVYKSTIu4ec_2VAxyuhAB",
					Status:                  "connected",
					ConnectedAt:             time.Now(),
					KeepaliveTimeoutSeconds: 40,
					ReconnectURL:            "wss://eventsub.wss.twitch.tv?XYZ",
				},
			},
		}
		got := sessionWelcome(welcome).ReconnectURL
		want := "wss://eventsub.wss.twitch.tv?XYZ"
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("should push live notification to Discord", func(t *testing.T) {
		notificationMessage := BaseMessage{
			Metadata: MessageMetadata{
				MessageId:        "123",
				MessageType:      "notification",
				MessageTimestamp: time.Now(),
			},
			Payload: MessagePayload{
				Subscription: MessageSubscription{
					Type: "channel.follow",
				},
				Event: MessageEvent{
					UserId:               1337,
					UserLogin:            "awesome_user",
					UserName:             "Awsome_User",
					BroadcasterUserId:    "12826",
					BroadcasterUserLogin: "awesome_broadcaster",
					BroadcasterUserName:  "Awesome_Broadcaster",
				},
			},
		}
		gotType := notification(notificationMessage).Type
		gotBroadcaster := notification(notificationMessage).BroadcasterName
		wantType := "channel.follow"
		wantBroadcaster := "Awesome_Broadcaster"
		if gotType != wantType {
			t.Errorf("got %q want %q", gotType, wantType)
		}
		if gotBroadcaster != wantBroadcaster {
			t.Errorf("got %q want %q", gotBroadcaster, wantBroadcaster)
		}
	})

	t.Run("should convert raw message to json", func(t *testing.T) {
		//Converted message of the raw below.
		//{{c16b4820-56a0-924a-bee1-2352c3265a65 session_welcome 2024-10-29 01:58:11.828583253 +0000 UTC  0} {{27856544_2e317b72 connected 2024-10-29 01:58:11.828465553 +0000 UTC 10 } {   0 0 {} { } 0001-01-01 00:00:00 +0000 UTC} {0   0   0001-01-01 00:00:00 +0000 UTC}}}"}
		rawMessage := []byte{123, 34, 109, 101, 116, 97, 100, 97, 116, 97, 34, 58, 123, 34, 109, 101, 115, 115, 97, 103, 101, 95, 105, 100, 34, 58, 34, 99, 49, 54, 98, 52, 56, 50, 48, 45, 53, 54, 97, 48, 45, 57, 50, 52, 97, 45, 98, 101, 101, 49, 45, 50, 51, 53, 50, 99, 51, 50, 54, 53, 97, 54, 53, 34, 44, 34, 109, 101, 115, 115, 97, 103, 101, 95, 116, 121, 112, 101, 34, 58, 34, 115, 101, 115, 115, 105, 111, 110, 95, 119, 101, 108, 99, 111, 109, 101, 34, 44, 34, 109, 101, 115, 115, 97, 103, 101, 95, 116, 105, 109, 101, 115, 116, 97, 109, 112, 34, 58, 34, 50, 48, 50, 52, 45, 49, 48, 45, 50, 57, 84, 48, 49, 58, 53, 56, 58, 49, 49, 46, 56, 50, 56, 53, 56, 51, 50, 53, 51, 90, 34, 125, 44, 34, 112, 97, 121, 108, 111, 97, 100, 34, 58, 123, 34, 115, 101, 115, 115, 105, 111, 110, 34, 58, 123, 34, 105, 100, 34, 58, 34, 50, 55, 56, 53, 54, 53, 52, 52, 95, 50, 101, 51, 49, 55, 98, 55, 50, 34, 44, 34, 115, 116, 97, 116, 117, 115, 34, 58, 34, 99, 111, 110, 110, 101, 99, 116, 101, 100, 34, 44, 34, 107, 101, 101, 112, 97, 108, 105, 118, 101, 95, 116, 105, 109, 101, 111, 117, 116, 95, 115, 101, 99, 111, 110, 100, 115, 34, 58, 49, 48, 44, 34, 114, 101, 99, 111, 110, 110, 101, 99, 116, 95, 117, 114, 108, 34, 58, 110, 117, 108, 108, 44, 34, 99, 111, 110, 110, 101, 99, 116, 101, 100, 95, 97, 116, 34, 58, 34, 50, 48, 50, 52, 45, 49, 48, 45, 50, 57, 84, 48, 49, 58, 53, 56, 58, 49, 49, 46, 56, 50, 56, 52, 54, 53, 53, 53, 51, 90, 34, 125, 125, 125}
		convertedMessage, _ := ConvertToJson(rawMessage)
		want := "session_welcome"
		got := convertedMessage.Metadata.MessageType

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}

	})

}
