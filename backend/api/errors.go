package api

import (
	"fmt"
)

type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Errors  any    `json:"errors,omitempty"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("api error: status %d | message: %s | errors: %s", e.Status, e.Message, e.Errors)
}

func NewAPIError(status int, message string, errors any) *APIError {
	return &APIError{
		Status:  status,
		Message: message,
		Errors:  errors,
	}
}

type fieldError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value any    `json:"value"`
}

func NewFieldError(field string, tag string, value any) *fieldError {
	return &fieldError{
		Field: field,
		Tag:   tag,
		Value: value,
	}
}

func (e fieldError) Error() string {
	return fmt.Sprintf("field: %s, failed tag: %s, value: %s", e.Field, e.Tag, e.Value)
}
