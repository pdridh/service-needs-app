package routes

import (
	"net/http"

	"github.com/pdridh/service-needs-app/backend/auth"
)

// Sets all the handlers to handle the appropriate routes
// All route handling is done inside this function so it acts as a map
// of all the routes (ideally)
func AddRoutes(
	mux *http.ServeMux,
	authHandler *auth.Handler,
) {
	mux.Handle("POST /auth/register", authHandler.Register())
	mux.Handle("POST /auth/login", authHandler.Login())
	mux.Handle("/", http.NotFoundHandler())
}
