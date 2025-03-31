package auth

import "errors"

var (
	ErrUnexpectedJWTSigningMethod = errors.New("unexpected signing method")
)
