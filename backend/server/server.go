package server

import (
	"net/http"

	"github.com/pdridh/service-needs-app/backend/auth"
	"github.com/pdridh/service-needs-app/backend/routes"
)

// Creates a new server, assigns it its routes with routes.AddRoutes()
// and all top level middlewares should be added here.
// Returns the new server as a http.Handler.
func New(authHandler *auth.Handler) http.Handler {
	mux := http.NewServeMux()

	// Add all the routes
	routes.AddRoutes(mux, authHandler)

	var handler http.Handler = mux

	// TODO add top level middlewares here (cors, ratelimiter, etc.)

	return handler
}
