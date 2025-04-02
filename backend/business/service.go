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

func (s *Service) GetBusinesses(filters bson.M, options *options.FindOptions) ([]Business, error) {
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
func (s *Service) GetReviews(filters bson.M, options *options.FindOptions) ([]review.Review, error) {
	return s.reviewStore.GetReviews(filters, options)
}
