package provider

import (
	"log"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/pdridh/service-needs-app/backend/api"
	"github.com/pdridh/service-needs-app/backend/auth"
)

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) GetProviders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		queries := r.URL.Query()

		page := api.GetIntParamFromQuery(queries, "page", 1, 1, 100)
		limit := api.GetIntParamFromQuery(queries, "limit", 10, 1, 50)

		validFilterKeys := []string{"location", "category"}
		filters := api.GetFiltersFromQuery(queries, validFilterKeys)

		ps, err := h.Service.GetProviders(filters, page, limit)
		if err != nil {
			api.WriteInternalError(w, r)
			return
		}

		log.Println(ps)
		api.WriteJSON(w, r, http.StatusOK, ps)
	}
}

func (h *Handler) RegisterProvider() http.HandlerFunc {
	type ProviderPayload struct {
		Name        string `json:"name" validate:"required"`
		Category    string `json:"category" validate:"required"`
		Location    string `json:"location" validate:"required"` // TODO change this to somehting better like coords or something idk
		Description string `json:"description"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var payload ProviderPayload
		u := auth.CurrentUserID(r)

		// Parse the payload
		if err := api.ParseJSON(r, &payload); err != nil {
			api.WriteError(w, r, http.StatusBadRequest, "Bad json request", nil)
			return
		}

		// Validate the payload
		var allErrs []error
		if err := h.Service.validate.Struct(payload); err != nil {
			for _, e := range err.(validator.ValidationErrors) {
				allErrs = append(allErrs, api.NewFieldError(e.Field(), e.Tag(), e.Value()))
			}
		}

		// TODO check if the name is available or already registered

		if len(allErrs) > 0 {
			// Handle all errors
			api.WriteError(w, r, http.StatusBadRequest, "Invalid form body", allErrs)
			return
		}

		p, err := h.Service.RegisterProvider(u, payload.Name, payload.Category, payload.Location, payload.Description)
		if err != nil {
			api.WriteInternalError(w, r)
			return
		}

		api.WriteJSON(w, r, http.StatusOK, p)
	}
}
