package errors

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
)

type APIError struct {
	StatusCode int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
}

func (e APIError) Error() string {
	return e.Message
}

func RespondWithError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case ValidationError:
		RespondWithValidationError(w, e)
		return
	case APIError:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(e.StatusCode)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]string{
				"code":    e.Code,
				"message": e.Message,
			},
		})
	default:
		apiErr := InternalServerError(err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(apiErr.StatusCode)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]string{
				"code":    apiErr.Code,
				"message": apiErr.Message,
			},
		})
	}
}

func BadRequest(message string) APIError {
	return APIError{
		StatusCode: http.StatusBadRequest,
		Code:       "BAD_REQUEST",
		Message:    message,
	}
}

func NotFound(message string) APIError {
	return APIError{
		StatusCode: http.StatusNotFound,
		Code:       "NOT_FOUND",
		Message:    message,
	}
}

func Conflict(message string) APIError {
	return APIError{
		StatusCode: http.StatusConflict,
		Code:       "CONFLICT",
		Message:    message,
	}
}

func InternalServerError(message string) APIError {
	return APIError{
		StatusCode: http.StatusInternalServerError,
		Code:       "INTERNAL_SERVER_ERROR",
		Message:    message,
	}
}

func Unauthorized(message string) APIError {
	return APIError{
		StatusCode: http.StatusUnauthorized,
		Code:       "UNAUTHORIZED",
		Message:    message,
	}
}

func Forbidden(message string) APIError {
	return APIError{
		StatusCode: http.StatusForbidden,
		Code:       "FORBIDDEN",
		Message:    message,
	}
}

type ValidationError struct {
	APIError
	Fields map[string]string `json:"fields"`
}

func NewValidationError(fields map[string]string) ValidationError {
	messages := make([]string, 0, len(fields))
	for field, msg := range fields {
		messages = append(messages, field+": "+msg)
	}
	return ValidationError{
		APIError: APIError{
			StatusCode: http.StatusBadRequest,
			Code:       "VALIDATION_ERROR",
			Message:    "Validation failed: " + strings.Join(messages, ", "),
		},
		Fields: fields,
	}
}

func RespondWithValidationError(w http.ResponseWriter, validationErr ValidationError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(validationErr.StatusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]interface{}{
			"code":    validationErr.Code,
			"message": validationErr.Message,
			"fields":  validationErr.Fields,
		},
	})
}

func ErrorHandlingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v\nStack trace: %s", err, debug.Stack())
				apiErr := InternalServerError("Internal server error")
				RespondWithError(w, apiErr)
			}
		}()
		next.ServeHTTP(w, r)
	})
}