package main

<<<<<<< HEAD
import "testing"

func TestHello(t *testing.T) {
	got := Hello("Chris")
	want := "Hello, Chris"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
=======
import (
	"testing"
	"time"
)

// Start websocket server locally with below command:
// twitch event websocket start-server
// Websocket server should start on: ws://127.0.0.1:8080/ws

func TestMessageStruct(t *testing.T) {
	t.Run("checking welcome message", func(t *testing.T) {
		welcome := Base_Message{
			Metadata: Message_Metadata{
				MessageID:        "123",
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
		welcome := Base_Message{
			Metadata: Message_Metadata{
				MessageID:        "123",
				MessageType:      "session_welcome",
				MessageTimestamp: time.Now(),
			},
			Payload: Message_Payload{
				Session: Message_Session{
					ID:                      "AQoQexAWVYKSTIu4ec_2VAxyuhAB",
					Status:                  "connected",
					ConnectedAt:             time.Now(),
					KeepaliveTimeoutSeconds: 40,
					ReconnectURL:            "wss://eventsub.wss.twitch.tv?XYZ",
				},
			},
		}
		timeout := 40
		got := session_welcome(welcome).TimeOut
		want := timeout
		if got != want {
			t.Errorf("got %d want %d", got, want)
		}
	})

	t.Run("checking session_welcome message reconnect url", func(t *testing.T) {
		welcome := Base_Message{
			Metadata: Message_Metadata{
				MessageID:        "123",
				MessageType:      "session_welcome",
				MessageTimestamp: time.Now(),
			},
			Payload: Message_Payload{
				Session: Message_Session{
					ID:                      "AQoQexAWVYKSTIu4ec_2VAxyuhAB",
					Status:                  "connected",
					ConnectedAt:             time.Now(),
					KeepaliveTimeoutSeconds: 40,
					ReconnectURL:            "wss://eventsub.wss.twitch.tv?XYZ",
				},
			},
		}
		got := session_welcome(welcome).ReconnectURL
		want := "wss://eventsub.wss.twitch.tv?XYZ"
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

>>>>>>> 0ad6af4 (feat: initial commit for message parser)
}
