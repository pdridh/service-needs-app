package api

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

// Write "v" (value that can be json marshaled) and send as response
// If the requestor accepts encoding then this encodes with gzip
func WriteJSON(w http.ResponseWriter, r *http.Request, status int, v any) error {

	b, err := json.Marshal(v)
	if err != nil {
		return errors.Wrap(err, "encode json")
	}

	var out io.Writer = w
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
		gzw := gzip.NewWriter(w)
		out = gzw
		defer gzw.Close()
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if _, err := out.Write(b); err != nil {
		return err
	}

	return nil
}

// A wrapper that takes creates an APIError with the given status, message and errors and writes that to the user as json with an appropriate status
func WriteError(w http.ResponseWriter, r *http.Request, status int, message string, errors any) {
	if err := WriteJSON(w, r, status, NewAPIError(status, message, errors)); err != nil {
		// TODO weird ahh error need to handle ts
	}
}

func WriteSuccess(w http.ResponseWriter, r *http.Request, status int, message string, data any) {
	s := NewSuccessResponse(status, message, data)

	if err := WriteJSON(w, r, status, s); err != nil {
		log.Println("failed to write to request")
	}
}

// Helper that calls WriteError() with args for an internal server error
func WriteInternalError(w http.ResponseWriter, r *http.Request) {
	WriteError(w, r, http.StatusInternalServerError, "Internal server error :(", nil)
}
