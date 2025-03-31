package utils

import (
	"compress/gzip"
	"encoding/json"
	"io"
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

// Read a request body and parse it.
// The parsed json is loaded into v unless an error occurs
func ParseJSON(r *http.Request, v any) error {

	b, err := io.ReadAll(io.LimitReader(r.Body, 1024*1024))
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, v)
	if err != nil {
		return err
	}

	return nil
}
