package ws

import (
	"github.com/coder/websocket"
)

// Hub maintains active clients and handles events
type Hub struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	events     chan []byte // TODO replace this with actual events
}

// NewHub creates a new hub instance
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		events:     make(chan []byte),
	}
}

// The only reason these functions dont require locks or anything is because the hub itself is using channles

// Simple helper to add a client to the client map.
// Handles the case when the same client joins from multiple devices (TODO)
func (h *Hub) RegisterClient(c *Client) {
	// TODO handle the case where client is already connected (for example two tabs opened idkididkdidkdidkd)
	h.clients[c.ID] = c
}

// Removes the given client c from the clients map if it exists.
// If it doesnt, it like a nop idk
func (h *Hub) UnregisterClient(c *Client) {
	if _, ok := h.clients[c.ID]; ok {
		c.Conn.Close(websocket.StatusNormalClosure, "byebye")
		delete(h.clients, c.ID)
	}
}

// This is the core function of ws pkg, it basically handles all the channels and also manages the clients and connections.
// * WARNING THIS SHOULD BE RUN DURING THE APPLICATION STARTUP
// TODO probably should add a context thingy idk
func (h *Hub) Run() {
	for {
		select {
		// Handle register channel
		case client := <-h.register:
			h.RegisterClient(client)
		// Handle unregister channel
		case client := <-h.unregister:
			h.UnregisterClient(client)
		case e := <-h.events:
			// FOR NOW ONLY BROADCASTS THE MESSAGE TO ALL (EVEN ECHO)
			for _, c := range h.clients {
				c.Send <- e
			}
		}
	}
}
