package auth

import (
	"net/http"
)

type contextKey string

const ContextUserKey contextKey = "user"

// Given a request extracts the value of the userID (string) from the context using the ContextUserKey
func CurrentUserID(r *http.Request) string {
	return r.Context().Value(ContextUserKey).(string)
}
