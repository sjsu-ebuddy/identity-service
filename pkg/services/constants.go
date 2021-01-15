package services

type key int

// Common constants and errors
const (
	InternalServerError = "INTERNAL_SERVER_ERROR"
	ValidatationError   = "VALIDATION_ERROR"
	Success             = "SUCCESS"
	BadRequest          = "BAD_REQUEST"
	Conflict            = "CONFLICT"
)

// Error messages
const (
	UserAlreadyExists = "An user with the same email already exists"
)

const (
	userContextKey key = iota
)
