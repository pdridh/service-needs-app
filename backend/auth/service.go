package auth

import (
	"context"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/pdridh/service-needs-app/backend/api"
	"github.com/pdridh/service-needs-app/backend/business"
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

// Given the detials for the user account and the details for either consumer type or business type user
// This functions starts a transaction and updates either the consumerstore & userStore OR businessStore & userStore
func (s *service) RegisterUser(email string, password string, userType string, businessInfo *api.BusinessPayload, consumerInfo *api.ConsumerPayload) (any, error) {
	// Start session for transaction
	session, err := s.client.StartSession()
	if err != nil {
		return nil, err
	}
	// End session after this closes
	defer session.EndSession(context.Background())

	var result any

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
			Type:     userType,
		}

		log.Println("User created")
		log.Println(u)

		if err := s.userStore.CreateUser(sc, u); err != nil {
			session.AbortTransaction(sc)
			return err
		}

		// Get user ID
		uid := u.ID

		// Create the associated entity (Business or Consumer)
		switch userType {
		case user.UserTypeBusiness:
			b := &business.Business{
				UserID:      uid,
				Name:        businessInfo.Name,
				Category:    businessInfo.Category,
				Location:    businessInfo.Location,
				Description: businessInfo.Description,
				Verified:    false,
			}
			if err := s.businessStore.CreateBusiness(sc, b); err != nil {
				session.AbortTransaction(sc)
				return err
			}
			result = b
		case user.UserTypeConsumer:
			c := &consumer.Consumer{
				UserID:    uid,
				FirstName: consumerInfo.FirstName,
				LastName:  consumerInfo.LastName,
				Verified:  false,
			}
			if err := s.consumerStore.CreateConsumer(sc, c); err != nil {
				session.AbortTransaction(sc)
				return err
			}
			result = c
		default:
			return ErrUnknownUserType
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

	return result, nil
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
	u, err := s.userStore.GetUserByEmail(email)

	if err != nil {
		return false, err
	}

	return u == nil, nil
}
