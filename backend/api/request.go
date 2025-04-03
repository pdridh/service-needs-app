package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
)

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

// Helper function to get a int param from query params.
// Given the key this function gets the string value using queries.Get().
// The value is defaultValue if there was an error converting or the string was not convertable.
// The value is clamped using minValue and maxValue.
func GetIntParamFromQuery(queries url.Values, key string, defaultValue int, minValue int, maxValue int) int {
	value, err := strconv.Atoi(queries.Get(key))
	if err != nil || value < minValue {
		return defaultValue
	}

	if value < minValue {
		return minValue
	}

	if value > maxValue {
		return maxValue
	}

	return value
}

// Helper function to extract bson.M filters from query params.
// Iterates over each validKeys []string array and checks if a param with that key was used.
// If the key is empty the filter is not applied otherwise the param's value is used as a filter with that key.
func GetFiltersFromQuery(queries url.Values, validKeys []string) bson.M {
	filters := bson.M{}

	for _, key := range validKeys {
		value := queries.Get(key)
		if value != "" {
			filters[key] = value
		}
	}

	return filters
}
