package domain

import "net/http"

// CognitoContext holds the authenticated user's identity from Cognito.
type CognitoContext struct {
	Sub    string `json:"sub"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Scope  string `json:"scope"`
	Tenant string `json:"tenant"`
}

// MapCoreError maps an HTTP status code and body to a domain error.
func MapCoreError(statusCode int, body []byte) error {
	switch statusCode {
	case http.StatusNotFound:
		return &NotFoundError{Resource: "resource"}
	case http.StatusConflict:
		return &ConflictError{Message: string(body)}
	case http.StatusGatewayTimeout:
		return &TimeoutError{Service: "core"}
	default:
		return &BusinessError{
			Code:    "CORE_ERROR",
			Message: "core-api returned an unexpected error",
			Status:  statusCode,
		}
	}
}

// MapTimeoutError returns a timeout error for the given service.
func MapTimeoutError(service string) error {
	return &TimeoutError{Service: service}
}
