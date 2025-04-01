package api

import (
	"net/http"
)

type contextKey string

const ContextUserKey contextKey = "user"

type CurrentUser struct {
	ID   string
	Type string
}

// Given a request extracts the value of the userID (string) from the context using the ContextUserKey
func CurrentUserID(r *http.Request) string {
	return r.Context().Value(ContextUserKey).(string)
}
