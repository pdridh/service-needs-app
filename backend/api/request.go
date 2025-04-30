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
func ParseQueryParams(q url.Values, dataHolder any) {
	dhType := reflect.TypeOf(dataHolder)
	dhVal := reflect.ValueOf(dataHolder)

	for i := range dhType.Elem().NumField() {
		field := dhType.Elem().Field(i)
		key := field.Tag.Get("json")
		kind := field.Type.Kind()

		queryVal := q.Get(key)

		fieldInput := dhVal.Elem().Field(i)

		if !fieldInput.CanSet() {
			continue
		}

		switch kind {
		case reflect.Int:
			intVal, err := strconv.ParseInt(queryVal, 10, 64)
			if err == nil {
				fieldInput.SetInt(intVal)
			}
		case reflect.String:
			fieldInput.SetString(queryVal)
		case reflect.Bool:
			val, err := strconv.ParseBool(queryVal)
			if err == nil {
				fieldInput.SetBool(val)
			}
		case reflect.Float64:
			val, err := strconv.ParseFloat(queryVal, 64)
			if err == nil {
				fieldInput.SetFloat(val)
			}
		case reflect.Struct, reflect.Map:
			val := reflect.New(field.Type)
			err := json.Unmarshal([]byte(queryVal), val.Interface())
			if err == nil {
				fieldInput.Set(val.Elem())
			}
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
