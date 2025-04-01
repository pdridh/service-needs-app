package auth

import (
	"errors"
)

var (
	ErrUnexpectedJWTSigningMethod = errors.New("unexpected signing method")
	ErrUnknownEmail               = errors.New("unknown email used")
	ErrWrongPassword              = errors.New("password doesnt match the identifier used")
	ErrUnknownUserType            = errors.New("unknown user type creation")
)
