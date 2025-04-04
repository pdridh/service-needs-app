package server

import (
	"net/http"

	"github.com/pdridh/service-needs-app/backend/auth"
	"github.com/pdridh/service-needs-app/backend/business"
	"github.com/pdridh/service-needs-app/backend/routes"
	"github.com/pdridh/service-needs-app/backend/ws"
)

// Creates a new server, assigns it its routes with routes.AddRoutes()
// and all top level middlewares should be added here.
// Returns the new server as a http.Handler.
func New(authHandler *auth.Handler, businessHandler *business.Handler, wsHandler *ws.Handler) http.Handler {
	mux := http.NewServeMux()

	// Add all the routes
	routes.AddRoutes(mux, authHandler, businessHandler, wsHandler)

	var handler http.Handler = mux

	// TODO add top level middlewares here (cors, ratelimiter, etc.)

	return handler
}
