package ws

import (
	"context"
	"fmt"
	"log"
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
		t := api.CurrentUserType(r)

		// Create a new client
		client := NewClient(conn, u, t, h.hub)
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
	type RequestPayload struct {
		Sender   string `json:"sender"`
		Receiver string `json:"receiver"`
		Message  string `json:"message"` // TODO probably handle attachments and shit (no idea how)
	}

	var p RequestPayload

	err := e.Event.ParsePayloadInto(&p)
	if err != nil {
		return
	}

	// TODO verify if the receiver is real and allows the sender to send message to it.
	// TODO something like: receiverStore.getReceiver(p.Receiver) -> check and then some kinda isAllowedTo(p.Sender, p.Receiver) or sum

	// If it passed verification the message sending responsibility is now to the server.

	msg := chat.NewChatMessage(e.Client.ID, p.Receiver, p.Message, chat.StatusMessageSent)
	if err := h.chatStore.CreateChatMessage(context.Background(), msg); err != nil {
		// TODO inform user that its server's fault the message wasnt delivered
		return
	}

	// If the receiver's connection is online we can send now
	if c, ok := h.clients[p.Receiver]; ok {

		msg.Status = chat.StatusMessageDelivered
		if err := h.chatStore.UpdateMessageStatus(context.Background(), msg.ID.Hex(), msg.Status); err != nil {
			// TODO Again, inform sender that the server failed to deliver
			log.Println(err)
			return
		}

		// Send the message to the receiver
		c.Send <- Event{Code: EventChat, Payload: msg}
	}

	// Update the sender about the message
	if c, ok := h.clients[e.Client.ID]; ok {
		c.Send <- Event{Code: EventChat, Payload: msg}
	}
}

func (h *Hub) HandleChatSeenEvent(e EventContext) {
	// TODO implement chat seen
}
