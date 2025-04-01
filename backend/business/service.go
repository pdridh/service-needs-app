package business

import (
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
	store    Store
	validate *validator.Validate
}

func NewService(store Store, validate *validator.Validate) *Service {
	return &Service{
		store:    store,
		validate: validate,
	}
}

func (s *Service) GetBusinesses(filters bson.M, options *options.FindOptions) ([]bson.M, error) {
	return s.store.GetBusinesses(filters, options)
}
