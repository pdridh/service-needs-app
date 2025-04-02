package routes

import (
	"net/http"

	"github.com/pdridh/service-needs-app/backend/auth"
	"github.com/pdridh/service-needs-app/backend/business"
	"github.com/pdridh/service-needs-app/backend/user"
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

	mux.Handle("GET /api/v1/businesses", auth.Middleware(businessHandler.GetBusinesses(), ""))

	mux.Handle("GET /api/v1/businesses/{id}/reviews", auth.Middleware(businessHandler.GetBusinessReviews(), ""))
	mux.Handle("POST /api/v1/businesses/{id}/reviews", auth.Middleware(businessHandler.AddReview(), user.UserTypeConsumer))

	mux.Handle("/", http.NotFoundHandler())
}
