package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

type Message_Metadata struct {
	MessageID        string    `json:"message_id"`
	MessageType      string    `json:"message_type"`
	MessageTimestamp time.Time `json:"message_timestamp"`
}

type Message_Session struct {
	ID                      string    `json:"id"`
	Status                  string    `json:"status"`
	ConnectedAt             time.Time `json:"connected_at"`
	KeepaliveTimeoutSeconds int       `json:"keepalive_timeout_seconds"`
	ReconnectURL            string    `json:"reconnect_url"`
}

type Message_Payload struct {
	Session Message_Session
}

type Base_Message struct {
	Metadata Message_Metadata
	Payload  Message_Payload
}

type Persistant_Session_Information struct {
	TimeOut      int
	ReconnectURL string
}

// https://dev.twitch.tv/docs/eventsub/handling-websocket-events/
func message_type_handler(message Base_Message) {
	switch message.Metadata.MessageType {
	case "session_welcome":
		session_welcome(message)
	case "session_keepalive":
	case "notification":
	case "session_reconnect":
	case "revocation":

	}
}

// The welcome message gives two important things:
// 1. The URL at which we can reconnect if we are disconnected from the client.
// 2. The KeepAliveTimeout that will let us know how long we should be waiting before we get a session_keepalive
func session_welcome(message Base_Message) Persistant_Session_Information {
	current_session := Persistant_Session_Information{}
	current_session.ReconnectURL = message.Payload.Session.ReconnectURL
	current_session.TimeOut = message.Payload.Session.KeepaliveTimeoutSeconds
	return current_session
}

func session_keepalive(message Base_Message) {

}

func notification(message Base_Message) {

}

func session_reconnect(message Base_Message) {

}

func revocation(message Base_Message) {

}

var addr = flag.String("addr", "localhost:8080", "twitch-cli ws service address")

func connection() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {

				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
			// parse message json for the message_type then cast it to one of the structs

			// send message to different handler funcs

		}
	}()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
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

<<<<<<< HEAD
func Hello(name string) string {
	return "Hello, " + name
}

func main() {
	fmt.Println(Hello("Chris"))
=======
func main() {
	connection()
>>>>>>> 0ad6af4 (feat: initial commit for message parser)
}
