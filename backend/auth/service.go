package auth

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/pdridh/service-needs-app/backend/business"
	"github.com/pdridh/service-needs-app/backend/common"
	"github.com/pdridh/service-needs-app/backend/config"
	"github.com/pdridh/service-needs-app/backend/consumer"
	"github.com/pdridh/service-needs-app/backend/user"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	client        *mongo.Client
	userStore     user.Store
	businessStore business.Store
	consumerStore consumer.Store
	validate      *validator.Validate
}

// Simple wrapper to create a new user service given the stores and validator
func NewService(client *mongo.Client, userStore user.Store, businessStore business.Store, consumerStore consumer.Store, validate *validator.Validate) *service {
	return &service{
		client:        client,
		userStore:     userStore,
		businessStore: businessStore,
		consumerStore: consumerStore,
		validate:      validate,
	}
}

func (s *service) RegisterBusiness(email string, password string, name string, category string, longitude float64, latitude float64, description string) (*business.Business, error) {
	// Start session for transaction
	session, err := s.client.StartSession()
	if err != nil {
		return nil, err
	}
	// End session after this closes
	defer session.EndSession(context.Background())

	var b *business.Business
	err = mongo.WithSession(context.Background(), session, func(sc mongo.SessionContext) error {
		// Start the transaction
		if err := session.StartTransaction(); err != nil {
			return err
		}

		// Register the user
		// Hash the password before storing
		hashedPassword, err := HashPassword(password)
		if err != nil {
			session.AbortTransaction(sc)
			return err
		}

		// Actual registration process
		u := &user.User{
			Email:    email,
			Password: hashedPassword,
			Type:     user.UserTypeBusiness,
		}

		if err := s.userStore.CreateUser(sc, u); err != nil {
			session.AbortTransaction(sc)
			return err
		}

		// Get user ID
		uid := u.ID

		b = &business.Business{
			ID:       uid,
			Name:     name,
			Category: category,
			Location: common.GeoLocation{
				Type:        "Point",
				Coordinates: []float64{longitude, latitude},
			},
			Description: description,
			Verified:    false,
		}
		if err := s.businessStore.CreateBusiness(sc, b); err != nil {
			session.AbortTransaction(sc)
			return err
		}

		// Commit the transaction
		if err := session.CommitTransaction(sc); err != nil {
			session.AbortTransaction(sc)
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return b, err
}

// Given the detials for the user account and the details for consumer.
// This functions starts a transaction and updates the consumerstore & userStore together, meaning if one fails the other also fails
// And the records are consistent.
func (s *service) RegisterConsumer(email string, password string, firstName string, lastName string) (*consumer.Consumer, error) {
	// Start session for transaction
	session, err := s.client.StartSession()
	if err != nil {
		return nil, err
	}

	// End session after this closes
	defer session.EndSession(context.Background())

	var c *consumer.Consumer
	err = mongo.WithSession(context.Background(), session, func(sc mongo.SessionContext) error {
		// Start the transaction
		if err := session.StartTransaction(); err != nil {
			return err
		}

		// Register the user
		// Hash the password before storing
		hashedPassword, err := HashPassword(password)
		if err != nil {
			session.AbortTransaction(sc)
			return err
		}

		// Actual registration process
		u := &user.User{
			Email:    email,
			Password: hashedPassword,
			Type:     user.UserTypeBusiness,
		}

		if err := s.userStore.CreateUser(sc, u); err != nil {
			session.AbortTransaction(sc)
			return err
		}

		// Get user ID
		uid := u.ID

		c = &consumer.Consumer{
			ID:        uid,
			FirstName: firstName,
			LastName:  lastName,
			Verified:  false,
		}

		if err := s.consumerStore.CreateConsumer(sc, c); err != nil {
			session.AbortTransaction(sc)
			return err
		}
		// Commit the transaction
		if err := session.CommitTransaction(sc); err != nil {
			session.AbortTransaction(sc)
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return c, nil
}

// Given a user email and password it checks stuff in the store and returns a jwt token (string) if everything went well.
// If anything didn't match it returns it as an error with a nil string. The errors are specific and can be a security vuln
// if used as a response in the api. The handler is supposed to handle the errors and show appropriate api responses.
func (s *service) AuthenticateUser(email string, password string) (string, error) {
	// Check if the indentifier is an email
	u, err := s.userStore.GetUserByEmail(email)
	if err != nil {
		return "", err
	}

	// If u is still nil then the identifier is just wrong
	if u == nil {
		return "", ErrUnknownEmail
	}

	// Check if the password is correct
	if err := CompareHashedPasswords(u.Password, password); err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return "", ErrWrongPassword
		default:
			return "", err
		}
	}

	// Everything went well so generate a token for the user
	t, err := GenerateJWT(u.ID.Hex(), u.Type, config.Server().JWTExpiration)
	if err != nil {
		return "", err
	}

	return t, nil
}

// Given an email this function checks if theres already a user in the store
// that uses this email.
// Returns true if email is available otherwise false, also can return an error.
func (s *service) IsEmailAvailable(email string) (bool, error) {
	u, err := s.userStore.GetUserByEmail(email)

	if err != nil {
		return false, err
	}

	return u == nil, nil
}
