package provider

import (
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
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

func (s *Service) GetProviders(filters bson.M, page int, limit int) ([]bson.M, error) {
	return s.store.GetProviders(filters, page, limit)
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
