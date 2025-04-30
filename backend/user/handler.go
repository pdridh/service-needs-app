package user

import (
	"net/http"

	"github.com/pdridh/service-needs-app/backend/api"
)

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) GetUsers() http.HandlerFunc {
	type ResponsePayload struct {
		Users []User `json:"users"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var params QueryOptions
		api.ParseQueryParams(r.URL.Query(), &params)

		u, _, err := h.Service.GetUsers(r.Context(), params)
		if err != nil {
			api.WriteInternalError(w, r)
			return
		}

		api.WriteJSON(w, r, http.StatusOK, ResponsePayload{Users: u})
	}
}
