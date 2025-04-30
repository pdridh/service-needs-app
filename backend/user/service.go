package user

import (
	"context"

	"github.com/go-playground/validator/v10"
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

func (s *Service) GetUsers(ctx context.Context, options QueryOptions) ([]User, int64, error) {
	return s.store.GetUsers(ctx, options)
}
