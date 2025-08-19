package server

import (
	"net/http"

	"github.com/pdridh/service-needs-app/backend/auth"
	"github.com/pdridh/service-needs-app/backend/business"
	"github.com/pdridh/service-needs-app/backend/routes"
	"github.com/pdridh/service-needs-app/backend/ws"
	"github.com/rs/cors"
)

// Creates a new server, assigns it its routes with routes.AddRoutes()
// and all top level middlewares should be added here.
// Returns the new server as a http.Handler.
func New(authHandler *auth.Handler, businessHandler *business.Handler, wsHandler *ws.Handler) http.Handler {
	mux := http.NewServeMux()

	// Add all the routes
	routes.AddRoutes(mux, authHandler, businessHandler, wsHandler)

	// TODO add top level middlewares here (cors, ratelimiter, etc.)

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // TODO read this from config
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	}).Handler(mux)

	return handler
}
