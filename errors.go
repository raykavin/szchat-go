package szchat

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// APIError represents a non-2xx response from the SZChat API. The response
// body shape is not fully uniform across endpoints, so FieldErrors and
// Message are populated on a best-effort basis; Body always holds the raw
// response for callers that need the original payload.
type APIError struct {
	StatusCode  int
	Message     string
	FieldErrors map[string][]string
	Body        []byte
}

func (e *APIError) Error() string {
	if len(e.FieldErrors) > 0 {
		return fmt.Sprintf("szchat: %d %s: %v", e.StatusCode, e.Message, e.FieldErrors)
	}
	return fmt.Sprintf("szchat: %d %s", e.StatusCode, e.Message)
}

// parseAPIError builds an *APIError from a non-2xx response, tolerating the
// several error envelope shapes used across the SZChat API:
//
//	{"error": "Unauthorized"}
//	{"status": false, "errors": {"field": ["msg", ...]}}
//	{"status": "fail", "response": "existing name"}
//	{"success": false, "message": "..."}
func parseAPIError(statusCode int, body []byte) *APIError {
	apiErr := &APIError{StatusCode: statusCode, Body: body}

	var payload struct {
		Error    string              `json:"error"`
		Errors   map[string][]string `json:"errors"`
		Response string              `json:"response"`
		Message  string              `json:"message"`
	}

	switch {
	case json.Unmarshal(body, &payload) == nil && (payload.Error != "" || len(payload.Errors) > 0 || payload.Response != "" || payload.Message != ""):
		apiErr.FieldErrors = payload.Errors
		switch {
		case payload.Error != "":
			apiErr.Message = payload.Error
		case len(payload.Errors) > 0:
			apiErr.Message = "validation failed"
		case payload.Response != "":
			apiErr.Message = payload.Response
		default:
			apiErr.Message = payload.Message
		}
	case len(body) > 0:
		apiErr.Message = strings.TrimSpace(string(body))
	default:
		apiErr.Message = http.StatusText(statusCode)
	}

	return apiErr
}

func statusIs(err error, code int) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == code
	}
	return false
}

// IsNotFound reports whether err is an *APIError with a 404 status.
func IsNotFound(err error) bool { return statusIs(err, http.StatusNotFound) }

// IsUnauthorized reports whether err is an *APIError with a 401 status.
func IsUnauthorized(err error) bool { return statusIs(err, http.StatusUnauthorized) }

// IsValidationError reports whether err is an *APIError with a 422 status.
func IsValidationError(err error) bool { return statusIs(err, http.StatusUnprocessableEntity) }

// IsConflict reports whether err is an *APIError with a 409 status.
func IsConflict(err error) bool { return statusIs(err, http.StatusConflict) }

// IsRateLimited reports whether err is an *APIError with a 429 status.
func IsRateLimited(err error) bool { return statusIs(err, http.StatusTooManyRequests) }
