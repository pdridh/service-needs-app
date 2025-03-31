package provider

import (
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
	store    ProviderStore
	validate *validator.Validate
}

func NewService(store ProviderStore, validate *validator.Validate) *Service {
	return &Service{
		store:    store,
		validate: validate,
	}
}

func (s *Service) GetProviders(filters bson.M, options *options.FindOptions) ([]bson.M, error) {
	return s.store.GetProviders(filters, options)
}

func (s *Service) RegisterProvider(userid string, name string, category string, location string, description string) (*Provider, error) {
	p := &Provider{
		Name:        name,
		Category:    category,
		Location:    location,
		Description: description,
	}

	if err := s.store.InsertProvider(userid, p); err != nil {
		return nil, err
	}

	return p, nil
}
