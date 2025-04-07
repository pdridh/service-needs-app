package ws

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/coder/websocket"
)

// Constants
const (
	maxMessageSize   = 512
	clientSendBuffer = 64
)

// Client represents a websocket connection
type Client struct {
	Conn *websocket.Conn
	ID   string
	Hub  *Hub
	Send chan Event
}

func NewClient(conn *websocket.Conn, id string, hub *Hub) *Client {
	return &Client{
		Conn: conn,
		ID:   id,
		Hub:  hub,
		Send: make(chan Event, clientSendBuffer),
	}
}

// Constantly reads from the connection and sends it to the hub.
// * Is responsible for cleaning up both it's and it's brother's (WritePump) client.
// TODO use actual structured, typed events instead of random payloads
// TODO also handle message types
func (c *Client) ReadPump(ctx context.Context, cancel context.CancelFunc) {
	defer func() {
		cancel()
		c.Hub.unregister <- c
	}()

	c.Conn.SetReadLimit(maxMessageSize)

	for {
		messageType, p, err := c.Conn.Read(ctx)
		if websocket.CloseStatus(err) != -1 {
			log.Println("Client closed connection?")
			return
		}

		if errors.Is(err, context.Canceled) {
			return
		}

		if err != nil {
			log.Println("Connection unexpectedly closed")
			return
		}

		// TODO handle other messagetypes for now only text
		if messageType == websocket.MessageText {

			// json.Decode
			var e Event

			if err := json.Unmarshal(p, &e); err != nil {
				// If it doesnt follow the protocol just drop it for now
				// TODO maybe track how many times the client didnt follow protocol
				continue
			}

			// If it does follow the protocol
			eventContext := EventContext{
				Event:     e,
				Client:    c,
				Timestamp: time.Now().UTC(),
			}

			c.Hub.eventRouter <- eventContext
		}
	}
}

// Constantly tries writing from the client's send channel to the ws connection.
// * WARNING WritePump doesnt really care about the connection and expects its brother readpump to handle it.
// * All WritePump does in cases of errors or if conn somehow closed is it cancels the given context. It expects whoever
// * Shares the context to handle cleanup
func (c *Client) WritePump(ctx context.Context, cancel context.CancelFunc) {
	defer func() {
		cancel()
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			if !ok {
				return
			}

			raw, err := json.Marshal(msg)
			if err != nil {
				// TODO THis is the apis problem and client should be notified probably so handle this error
				return // FOR NOW CLOSE THE PUMP AND THE CONN
			}

			if err := c.Conn.Write(ctx, websocket.MessageText, raw); err != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
