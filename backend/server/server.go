package server

import (
	"net/http"

	"github.com/pdridh/service-needs-app/backend/routes"
)

// Creates a new server, assigns it its routes with routes.AddRoutes()
// and all top level middlewares should be added here.
// Returns the new server as a http.Handler.
func New() http.Handler {
	mux := http.NewServeMux()

	// Add all the routes
	routes.AddRoutes(mux)

	var handler http.Handler = mux

	// TODO add top level middlewares here (cors, ratelimiter, etc.)

	return handler
}
