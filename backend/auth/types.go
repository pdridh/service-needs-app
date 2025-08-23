package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/pdridh/service-needs-app/backend/user"
)

type UserClaims struct {
	UserID    string
	UserType  user.UserType
	UserEmail string
	jwt.RegisteredClaims
}
