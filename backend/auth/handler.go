package auth

import (
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
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

func (h *Handler) RegisterBusiness() http.HandlerFunc {

	type RequestPayload struct {
		Email       string  `json:"email" validate:"required,email"`
		Password    string  `json:"password" validate:"required,min=8,max=70"`
		Name        string  `json:"name" validate:"required,min=3,max=30"`
		Category    string  `json:"category" validate:"required"`
		Longitude   float64 `json:"longitude" validate:"required,min=-180,max=180"`
		Latitude    float64 `json:"latitude" validate:"required,min=-90,max=90"`
		Description string  `json:"description"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Read body into the payload
		var p RequestPayload

		if err := api.ParseJSON(r, &p); err != nil {
			log.Println(err)
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

		b, err := h.Service.RegisterBusiness(p.Email, p.Password, p.Name, p.Category, p.Longitude, p.Latitude, p.Description)
		if err != nil {
			api.WriteInternalError(w, r)
			return
		}

		// Passed everything and therefore let the frontend know it was a success
		api.WriteJSON(w, r, http.StatusCreated, b)
	}
}

// TODO make this less redudant? idk.. chose simplicity and redudancy over complexity and no redudancy; but theres probably a better solution here
func (h *Handler) RegisterConsumer() http.HandlerFunc {

	type RequestPayload struct {
		Email     string `json:"email" validate:"required,email"`
		Password  string `json:"password" validate:"required,min=8,max=70"`
		FirstName string `json:"firstName" validate:"required,min=3,max=20"`
		LastName  string `json:"lastName" validate:"required,min=2,max=20"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Read body into the payload
		var p RequestPayload

		if err := api.ParseJSON(r, &p); err != nil {
			log.Println(err)
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

		c, err := h.Service.RegisterConsumer(p.Email, p.Password, p.FirstName, p.LastName)
		if err != nil {
			api.WriteInternalError(w, r)
			return
		}

		// Passed everything and therefore let the frontend know it was a success
		api.WriteJSON(w, r, http.StatusCreated, c)
	}
}

func (h *Handler) Login() http.HandlerFunc {
	type LoginPayload struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8,max=70"`
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

		api.WriteSuccess(w, r, http.StatusOK, "Login succesfull!", nil)
	}
}

func (h *Handler) GetAuth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jCookie, err := r.Cookie("jwt")
		if err != nil {
			api.WriteError(w, r, http.StatusBadRequest, "Bad json request", nil)
			return
		}

		j := jCookie.Value

		t, err := ValidateJWT(j)
		if err != nil {
			api.WriteError(w, r, http.StatusUnauthorized, "invalid token", nil)
			return
		}

		c, err := UserClaimsFromJWT(t)
		if err != nil {
			api.WriteError(w, r, http.StatusUnauthorized, "invalid token", nil)
			return
		}

		auth := AuthUser{
			ID:   c.UserID,
			Type: c.UserType,
		}

		api.WriteSuccess(w, r, http.StatusOK, "Retrival succesful", auth)
	}
}
