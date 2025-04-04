package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"reflect"
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

// Given the query params in the url (for GETs) and a struct that has strongly typed fields.
// Extracts the field from the query param if it exists and parses it into the type of the field it is.
func ParseQueryParams(q url.Values, params any) {
	v := reflect.ValueOf(params).Elem()

	for i := range v.NumField() {
		field := v.Field(i)
		fieldType := v.Type().Field(i)
		paramValue := q.Get(fieldType.Tag.Get("json")) // Get from query

		if paramValue == "" {
			continue
		}

		switch field.Kind() {
		case reflect.Int64, reflect.Int32, reflect.Int:
			val, err := strconv.Atoi(paramValue)
			if err == nil {
				field.SetInt(int64(val))
			}
		case reflect.Float64:
			val, err := strconv.ParseFloat(paramValue, 64)
			if err == nil {
				field.SetFloat(val)
			}
		case reflect.Bool:
			val, err := strconv.ParseBool(paramValue)
			if err == nil {
				field.SetBool(val)
			}
		default:
			field.SetString(paramValue)
		}
	}
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
