package business

import (
	"net/http"

	"github.com/pdridh/service-needs-app/backend/api"
	"github.com/pdridh/service-needs-app/backend/review"
	"github.com/pdridh/service-needs-app/backend/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		var params QueryOptions

		api.ParseQueryParams(r.URL.Query(), &params)

		b, _, err := h.Service.GetBusinesses(r.Context(), params)
		if err != nil {
			api.WriteInternalError(w, r)
			return
		}

		api.WriteJSON(w, r, http.StatusOK, b)
	}
}

func (h *Handler) AddReview() http.HandlerFunc {

	type ReviewPayload struct {
		Rating  float32 `json:"rating" validate:"required,min=0,max=5"`
		Comment string  `json:"comment" validate:"required,min=3,max=100"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		businessIDStr := r.PathValue("id")
		consumerIDstr := user.CurrentUserID(r)

		var p ReviewPayload

		if err := api.ParseJSON(r, &p); err != nil {
			api.WriteError(w, r, http.StatusBadRequest, "Bad json request", nil)
			return
		}

		bid, err := primitive.ObjectIDFromHex(businessIDStr)
		if err != nil {
			api.WriteError(w, r, http.StatusNotFound, "Business not found", nil)
			return
		}

		cid, err := primitive.ObjectIDFromHex(consumerIDstr)
		if err != nil {
			api.WriteInternalError(w, r)
			return
		}

		// Check if the business id is valid
		valid, err := h.Service.IsValidID(bid.Hex())
		if err != nil {
			api.WriteInternalError(w, r)
			return
		}

		if !valid {
			api.WriteError(w, r, http.StatusNotFound, "Business not found", nil)
			return
		}

		// Check if theres already a review with this business id and consumer id combo
		filters := bson.M{"business_id": bid, "consumer_id": cid}
		// All reviews that have business id and consumer id that matches
		results, err := h.Service.GetReviews(filters, nil)
		if err != nil {
			api.WriteInternalError(w, r)
			return
		}

		// If found that means theres already a review made by this consumer
		if len(results) > 0 {
			api.WriteError(w, r, http.StatusConflict, "already reviewed, edit review in some other path (cuz u not using frontend u bad)", nil)
			return
		}

		// If it got to this point we are good to go ahead
		review := &review.Review{
			BusinessID: bid,
			ConsumerID: cid,
			Rating:     p.Rating,
			Comment:    p.Comment,
		}

		if err := h.Service.AddReview(review); err != nil {
			api.WriteInternalError(w, r)
			return
		}

		// TODO replace WriteJSON with some kinda WriteSuccess type shit
		// TODO also make the response more standard like "data" or something idk
		api.WriteJSON(w, r, http.StatusCreated, review)
	}
}

func (h *Handler) GetBusinessReviews() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		businessIDStr := r.PathValue("id")
		bid, err := primitive.ObjectIDFromHex(businessIDStr)
		if err != nil {
			api.WriteError(w, r, http.StatusNotFound, "Business not found", nil)
			return
		}

		// Check if the business id is valid
		valid, err := h.Service.IsValidID(bid.Hex())
		if err != nil {
			api.WriteInternalError(w, r)
			return
		}

		if !valid {
			api.WriteError(w, r, http.StatusNotFound, "Business not found", nil)
			return
		}

		// TODO FIX THIS, this whole query shit is very redundant LITERALLY copy pasted from mathi ko "GetBusinesses" maybe fix that ionno

		queries := r.URL.Query()

		filters := bson.M{"business_id": bid}
		findOptions := options.Find()

		// TODO icl this feels like a security vuln since im not checking what we are sorting by, idk its just my cybersec side tingling
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
		reviews, err := h.Service.GetReviews(filters, findOptions)
		if err != nil {
			api.WriteInternalError(w, r)
			return
		}

		// TODO again, make ts more standard
		api.WriteJSON(w, r, http.StatusOK, reviews)
	}
}
