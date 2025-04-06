package ws

import (
	"context"
	"fmt"
	"net/http"

	"github.com/coder/websocket"
	"github.com/pdridh/service-needs-app/backend/api"
)

type Handler struct {
	hub *Hub
}

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
