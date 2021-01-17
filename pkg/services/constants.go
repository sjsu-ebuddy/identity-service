package services

type key int

// Common constants and errors
const (
	InternalServerError = "INTERNAL_SERVER_ERROR"
	ValidatationError   = "VALIDATION_ERROR"
	Success             = "SUCCESS"
	BadRequest          = "BAD_REQUEST"
	Conflict            = "CONFLICT"
	NotFound            = "NOT_FOUND"
	UnauthorizedAccess  = "UNAUTHORIZED_ACCESS"
	AuthorizationHeader = "Authorization"
)

// Error messages
const (
	UserAlreadyExists = "An user with the same email already exists"
	NoSuchUser        = "No such user exists"
	InvalidPassword   = "Invalid password"
	InvalidToken      = "Token passed is invalid"
)

const (
	userContextKey key = iota
	loginContextKey
)
