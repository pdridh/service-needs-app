package ws

import (
	"context"
	"fmt"
	"net/http"

	"github.com/coder/websocket"
	"github.com/pdridh/service-needs-app/backend/api"
	"github.com/pdridh/service-needs-app/backend/chat"
)

type Handler struct {
	hub *Hub
}

// Given the event context handles the event.
type EventHandler func(e EventContext)

func NewHandler(hub *Hub) *Handler {
	return &Handler{hub: hub}
}

func (h *Handler) Accept() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil) // TODO make this more secure and shti
		if err != nil {
			api.WriteInternalError(w, r)
			return
		}

		u := api.CurrentUserID(r)

		// Create a new client
		client := NewClient(conn, u, h.hub)
		h.hub.register <- client

		ctx, cancel := context.WithCancel(context.Background())

		// Start client pumps for reading and writing
		go client.WritePump(ctx, cancel)
		go client.ReadPump(ctx, cancel)
	}
}

// Sends hello back to the client with its id (tester function probably temproray)
func HandleHelloEvent(e EventContext) {
	e.Client.Send <- Event{Code: EventHello, Payload: EventHelloPayload{Message: fmt.Sprintf("Hello %s", e.Client.ID)}}
}

func (h *Hub) HandleChatEvent(e EventContext) {
	var p EventChatPayload

	err := e.Event.ParsePayloadInto(&p)
	if err != nil {
		return
	}
	// TODO verify if the receiver is real and allows the sender to send message to it.
	// TODO something like: receiverStore.getReceiver(p.Receiver) -> check and then some kinda isAllowedTo(p.Sender, p.Receiver) or sum

	p.Timestamp = e.Timestamp
	p.Sender = e.Client.ID

	if c, ok := h.clients[p.Receiver]; ok {
		c.Send <- Event{Code: EventChat, Payload: p}

		// Can store after sending to the conn cuz its more responsive
		h.chatStore.CreateChatMessage(context.Background(), chat.NewChatMessage(p.Sender, p.Receiver, p.Message, p.Timestamp))
	}
}
