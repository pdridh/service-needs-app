package ws

import (
	"log"
	"sync"
)

// Hub maintains active clients and broadcasts messages
type Hub struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	mu         sync.Mutex
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
			h.mu.Lock()
			h.clients[client.ID] = client
			h.mu.Unlock()
			log.Println("Client registered", client.ID)
		// Handle unregister channel
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)
				log.Printf("Client unregistered: %s", client.ID)
			}
			h.mu.Unlock()
			// TODO add message stuff
		}
	}
}
