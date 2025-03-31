package auth

import (
	"github.com/go-playground/validator"
	"github.com/pdridh/service-needs-app/backend/config"
	"github.com/pdridh/service-needs-app/backend/user"
)

type service struct {
	store    user.UserStore
	validate *validator.Validate
}

// Simple wrapper to create a new user service given the user store (interface) and validator
func NewService(store user.UserStore, validate *validator.Validate) *service {
	return &service{
		store:    store,
		validate: validate,
	}
}

// Uses the user store to insert a new user based on the given email and password.
// Also responsible for hasing the password.
// If there are any errors it returns nil and the error.
// WARNING THIS FUNCTION IS REALLY STUPID AND DOESNT DO ANY CHECKS JUST STORES THE USER
// USING THE GIVEN ARGUMENTS.
func (s *service) RegisterUser(email string, password string) (*user.User, error) {

	// Actual registration process
	u := &user.User{
		Email:    email,
		Password: password, //TODO hash this shit
	}

	if err := s.store.InsertUser(u); err != nil {
		return nil, err
	}

	return u, nil
}

// Given a user email and password it checks stuff in the store and returns a jwt token (string) if everything went well.
// If anything didn't match it returns it as an error with a nil string. The errors are specific and can be a security vuln
// if used as a response in the api. The handler is supposed to handle the errors and show appropriate api responses.
func (s *service) AuthenticateUser(email string, password string) (string, error) {
	// Check if the indentifier is an email
	u, err := s.store.GetUserByEmail(email)
	if err != nil {
		return "", err
	}

	// If u is still nil then the identifier is just wrong
	if u == nil {
		return "", ErrUnknownEmail
	}

	// Check if the password is correct
	// TODO should be checking the hash and shit
	if password != u.Password {
		return "", ErrWrongPassword
	}

	// Everything went well so generate a token for the user
	t, err := GenerateJWT(u.ID.Hex(), config.Server().JWTExpiration)
	if err != nil {
		return "", err
	}

	return t, nil
}

// Given an email this function checks if theres already a user in the store
// that uses this email.
// Returns true if email is available otherwise false, also can return an error.
func (s *service) IsEmailAvailable(email string) (bool, error) {
	u, err := s.store.GetUserByEmail(email)

	if err != nil {
		return false, err
	}

	return u == nil, nil
}
