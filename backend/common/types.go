package common

import "github.com/golang-jwt/jwt/v5"

type UserClaims struct {
	UserID string
	jwt.RegisteredClaims
}
