// Package domain defines errors and value objects for bff-api-go-collections-client-late-paying.
package domain

import "fmt"

// BusinessError represents a business-level error returned to the caller.
type BusinessError struct {
	Code    string
	Message string
	Status  int
}

func (e *BusinessError) Error() string {
	return fmt.Sprintf("business error [%s]: %s", e.Code, e.Message)
}

// NotFoundError indicates a resource was not found.
type NotFoundError struct {
	Resource string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("resource not found: %s", e.Resource)
}

// ConflictError indicates a conflict (e.g., duplicate key).
type ConflictError struct {
	Message string
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("conflict: %s", e.Message)
}

// TimeoutError indicates a upstream timeout.
type TimeoutError struct {
	Service string
}

func (e *TimeoutError) Error() string {
	return fmt.Sprintf("timeout: service %s did not respond in time", e.Service)
}
