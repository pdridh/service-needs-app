package ws

import (
	"github.com/coder/websocket"
)

// Client represents a websocket connection
type Client struct {
	Conn *websocket.Conn
	ID   string
	Hub  *Hub
}
