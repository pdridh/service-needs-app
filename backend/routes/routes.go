package routes

import (
	"fmt"
	"net/http"

	"github.com/pdridh/service-needs-app/backend/api"
	"github.com/pdridh/service-needs-app/backend/auth"
	"github.com/pdridh/service-needs-app/backend/business"
)

// Sets all the handlers to handle the appropriate routes
// All route handling is done inside this function so it acts as a map
// of all the routes (ideally)
func AddRoutes(
	mux *http.ServeMux,
	authHandler *auth.Handler,
	businessHandler *business.Handler,
) {
	mux.Handle("POST /auth/register", authHandler.Register())
	mux.Handle("POST /auth/login", authHandler.Login())

	mux.Handle("GET /protected", auth.Middleware(ProtectedHandler()))

	mux.Handle("GET /api/v1/businesses", auth.Middleware(businessHandler.GetBusinesses()))

	mux.Handle("/", http.NotFoundHandler())
}

func ProtectedHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		api.WriteJSON(w, r, http.StatusOK, fmt.Sprintf("Hello %s", auth.CurrentUserID(r)))
	}
}
