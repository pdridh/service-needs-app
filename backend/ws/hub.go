package ws

import (
	"log"

	"github.com/coder/websocket"
)

// Hub maintains active clients and broadcasts messages
type Hub struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
}

// NewHub creates a new hub instance
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		// Handle register channel
		case client := <-h.register:
			h.clients[client.ID] = client
		// Handle unregister channel
		case client := <-h.unregister:
			if c, ok := h.clients[client.ID]; ok {
				c.Conn.Close(websocket.StatusNormalClosure, "byebye")
				delete(h.clients, client.ID)
				log.Printf("Client unregistered: %s", client.ID)
			}
			// TODO add message stuff
		}
	}
}
