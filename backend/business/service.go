package business

import (
	"github.com/go-playground/validator/v10"
	"github.com/pdridh/service-needs-app/backend/review"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
	businessStore Store
	reviewStore   review.Store
	validate      *validator.Validate
}

func NewService(bstore Store, rstore review.Store, validate *validator.Validate) *Service {
	return &Service{
		businessStore: bstore,
		reviewStore:   rstore,
		validate:      validate,
	}
}

func (s *Service) GetBusinesses(filters bson.M, options *options.FindOptions) ([]bson.M, error) {
	return s.businessStore.GetBusinesses(filters, options)
}

func (s *Service) IsValidID(id string) (bool, error) {
	b, err := s.businessStore.GetBusinessByID(id)

	if err != nil {
		return false, err
	}

	// If a business is found by this id (not nil) then its valid
	return b != nil, nil
}

func (s *Service) AddReview(r *review.Review) error {
	return s.reviewStore.CreateReview(r)
}

// TODO this kind of shares some shit with consumer and shit and review and shit, idk maybe we can move it idk idk idk
func (s *Service) GetReviews(filters bson.M, options *options.FindOptions) ([]bson.M, error) {
	return s.reviewStore.GetReviews(filters, options)
}

func (s *Service) GetBusinessReviews() {

	// validFilterKeys := []string{"location", "category"}
	// filters := api.GetFiltersFromQuery(queries, validFilterKeys)

	// findOptions := options.Find()

	// // Sorting stuff
	// if sortBy, sortOrder := queries.Get("sortBy"), queries.Get("sortOrder"); sortBy != "" && sortOrder != "" {
	// 	order := 1
	// 	if sortOrder == "desc" {
	// 		order = -1
	// 	}
	// 	findOptions.SetSort(bson.D{{Key: sortBy, Value: order}})
	// }

	// // Pagination stuff
	// page := api.GetIntParamFromQuery(queries, "page", 1, 1, 100)
	// limit := api.GetIntParamFromQuery(queries, "limit", 10, 1, 50)
	// skip := (page - 1) * limit
	// findOptions.SetLimit(int64(limit)).SetSkip(int64(skip))

	// // Finally after applying all the filters and options query the store
	// ps, err := h.Service.GetBusinesses(filters, findOptions)
	// if err != nil {
	// 	api.WriteInternalError(w, r)
	// 	return
	// }
}
