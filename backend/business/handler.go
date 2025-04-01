package business

import (
	"net/http"

	"github.com/pdridh/service-needs-app/backend/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) GetBusinesses() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		queries := r.URL.Query()

		// Filter stuff
		validFilterKeys := []string{"location", "category"}
		filters := api.GetFiltersFromQuery(queries, validFilterKeys)

		findOptions := options.Find()

		// Sorting stuff
		if sortBy, sortOrder := queries.Get("sortBy"), queries.Get("sortOrder"); sortBy != "" && sortOrder != "" {
			order := 1
			if sortOrder == "desc" {
				order = -1
			}
			findOptions.SetSort(bson.D{{Key: sortBy, Value: order}})
		}

		// Pagination stuff
		page := api.GetIntParamFromQuery(queries, "page", 1, 1, 100)
		limit := api.GetIntParamFromQuery(queries, "limit", 10, 1, 50)
		skip := (page - 1) * limit
		findOptions.SetLimit(int64(limit)).SetSkip(int64(skip))

		// Finally after applying all the filters and options query the store
		ps, err := h.Service.GetBusinesses(filters, findOptions)
		if err != nil {
			api.WriteInternalError(w, r)
			return
		}

		api.WriteJSON(w, r, http.StatusOK, ps)
	}
}
