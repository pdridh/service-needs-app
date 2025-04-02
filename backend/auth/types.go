package auth

import "github.com/golang-jwt/jwt/v5"

type UserClaims struct {
	UserID   string
	UserType string
	jwt.RegisteredClaims
}
