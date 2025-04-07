package ws

import (
	"github.com/coder/websocket"
	"github.com/pdridh/service-needs-app/backend/chat"
)

// Hub maintains active clients and handles events
type Hub struct {
	clients       map[string]*Client
	register      chan *Client
	unregister    chan *Client
	eventRouter   chan EventContext
	eventHandlers map[EventCode]EventHandler
	chatStore     chat.Store
}

// NewHub creates a new hub instance.
func NewHub(chatStore chat.Store) *Hub {
	h := &Hub{
		clients:       make(map[string]*Client),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		eventRouter:   make(chan EventContext),
		eventHandlers: make(map[EventCode]EventHandler),
		chatStore:     chatStore,
	}

	// Assign handlers here
	h.On(EventHello, HandleHelloEvent)
	h.On(EventChat, h.HandleChatEvent)
	h.On(EventChatSeen, h.HandleChatSeenEvent)

	return h
}

// Simple wrapper function that assigns the given event string to be handled by the given EventHandler.
// Overwrites the previous handler if it has already been assigned.
func (h *Hub) On(e EventCode, f EventHandler) {
	h.eventHandlers[e] = f
}

// Given the event context, checks the eventHandlers map and finds the assigned event to the give event code.
// If found it executes the handler otherwise its a nop.
func (h *Hub) RouteEvent(e EventContext) {
	if handler, ok := h.eventHandlers[e.Event.Code]; ok {
		handler(e)
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
		case e := <-h.eventRouter:
			h.RouteEvent(e)
		}
	}
}
