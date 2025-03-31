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
	mux.Handle("GET /", helloWorldHandler())
	mux.Handle("POST /auth/register", authHandler.Register())
	mux.Handle("/", http.NotFoundHandler())
}

func helloWorldHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world"))
	}
}
