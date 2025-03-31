package auth

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/pdridh/service-needs-app/backend/api"
)

type Handler struct {
	Service *service
}

// Simple wrapper to create a new auth handler that uses the given service (can be changed for mock and shit)
func NewHandler(service *service) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) Register() http.HandlerFunc {
	type RegisterPayload struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8,max=80"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Read body into the payload
		var p RegisterPayload

		if err := api.ParseJSON(r, &p); err != nil {
			api.WriteError(w, r, http.StatusBadRequest, "Bad json request", nil)
			return
		}

		var allErrs []error
		if err := h.Service.validate.Struct(p); err != nil {
			for _, e := range err.(validator.ValidationErrors) {
				allErrs = append(allErrs, api.NewFieldError(e.Field(), e.Tag(), e.Value()))
			}
		}

		available, err := h.Service.IsEmailAvailable(p.Email)
		if err != nil {
			api.WriteInternalError(w, r)
			return
		}

		if !available {
			allErrs = append(allErrs, api.NewFieldError("Email", "conflict", p.Email))
		}

		if len(allErrs) > 0 {
			// Handle all errors
			api.WriteError(w, r, http.StatusBadRequest, "Invalid form body", allErrs)
			return
		}

		u, err := h.Service.RegisterUser(p.Email, p.Password)
		if err != nil {
			api.WriteInternalError(w, r)
			return
		}

		// Passed everything and therefore user is registered
		// Inform the user that it was succesful
		api.WriteJSON(w, r, http.StatusCreated, u)
	}
}

func (h *Handler) Login() http.HandlerFunc {
	type LoginPayload struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8,max=80"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Read body into the payload
		var p LoginPayload

		if err := api.ParseJSON(r, &p); err != nil {
			api.WriteError(w, r, http.StatusBadRequest, "Bad json request", nil)
			return
		}

		var validationErrs []error
		if err := h.Service.validate.Struct(p); err != nil {
			for _, e := range err.(validator.ValidationErrors) {
				validationErrs = append(validationErrs, api.NewFieldError(e.Field(), e.Tag(), e.Value()))
			}
		}

		if len(validationErrs) > 0 {
			// Handle all errors
			api.WriteError(w, r, http.StatusBadRequest, "Invalid form body", validationErrs)
			return
		}

		t, err := h.Service.AuthenticateUser(p.Email, p.Password)
		if err != nil {
			switch err {
			case ErrUnknownEmail, ErrWrongPassword:
				api.WriteError(w, r, http.StatusUnauthorized, "Invalid credentials", nil)
				return
			default:
				api.WriteInternalError(w, r)
				return
			}
		}

		// Passed everything meaning the authentication was succesful, lets give the user a token
		SetJWTCookie(w, t)
		api.WriteJSON(w, r, http.StatusOK, "Login succesfull!")
	}
}
