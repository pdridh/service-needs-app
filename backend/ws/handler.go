package ws

import (
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
		client := &Client{
			ID:   u,
			Conn: conn,
			Hub:  h.hub,
		}

		h.hub.register <- client

	}
}
