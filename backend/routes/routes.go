package routes

import (
	"net/http"

	"github.com/pdridh/service-needs-app/backend/auth"
	"github.com/pdridh/service-needs-app/backend/business"
	"github.com/pdridh/service-needs-app/backend/user"
	"github.com/pdridh/service-needs-app/backend/ws"
)

// Sets all the handlers to handle the appropriate routes
// All route handling is done inside this function so it acts as a map
// of all the routes (ideally)
func AddRoutes(
	mux *http.ServeMux,
	wsHandler *ws.Handler,
	userHandler *user.Handler,
	businessHandler *business.Handler,
	authHandler *auth.Handler,
) {
	mux.Handle("POST /auth/register/businesses", authHandler.RegisterBusiness())
	mux.Handle("POST /auth/register/consumers", authHandler.RegisterConsumer())

	mux.Handle("POST /auth/login", authHandler.Login())
	mux.Handle("/ws", auth.Middleware(wsHandler.Accept()))

	mux.Handle("GET /api/v1/businesses", auth.Middleware(businessHandler.GetBusinesses()))

	mux.Handle("GET /api/v1/businesses/{id}/reviews", auth.Middleware(businessHandler.GetBusinessReviews()))
	mux.Handle("POST /api/v1/businesses/{id}/reviews", auth.Middleware(businessHandler.AddReview(), user.UserTypeConsumer))

	// Admin routes
	mux.Handle("GET /api/v1/users", auth.Middleware(userHandler.GetUsers(), user.UserTypeAdmin))

	mux.Handle("/", http.NotFoundHandler())
}
