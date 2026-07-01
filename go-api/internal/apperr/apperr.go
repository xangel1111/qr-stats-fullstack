// Package apperr defines a typed application error that carries the HTTP
// status, a machine-readable code and a human-readable message. Handlers and
// services return these; the central error middleware turns them into the
// shared JSON error envelope.
package apperr

// APIError is an error with an associated HTTP status and API error code.
type APIError struct {
	Status  int
	Code    string
	Message string
	Details []string
}

func (e *APIError) Error() string { return e.Message }

// New builds an APIError.
func New(status int, code, message string) *APIError {
	return &APIError{Status: status, Code: code, Message: message}
}

// The constructors below map 1:1 to the error catalog documented in the API
// contract, so the HTTP status and code stay consistent across the service.

func BadRequest(message string) *APIError {
	return New(400, "BAD_REQUEST", message)
}

func Unauthorized(message string) *APIError {
	return New(401, "UNAUTHORIZED", message)
}

func NotFound(message string) *APIError {
	return New(404, "NOT_FOUND", message)
}

func Validation(message string) *APIError {
	return New(422, "VALIDATION_ERROR", message)
}

func Upstream(message string) *APIError {
	return New(502, "UPSTREAM_ERROR", message)
}

func Internal(message string) *APIError {
	return New(500, "INTERNAL_ERROR", message)
}
