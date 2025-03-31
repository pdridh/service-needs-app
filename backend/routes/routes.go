package routes

import "net/http"

// Sets all the handlers to handle the appropriate routes
// All route handling is done inside this function so it acts as a map
// of all the routes (ideally)
func AddRoutes(
	mux *http.ServeMux,
) {

	mux.Handle("GET /", helloWorldHandler())

}

func helloWorldHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world"))
	}
}
