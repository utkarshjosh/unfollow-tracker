package domain

import "errors"

// Domain errors
var (
	// User errors
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")

	// Account errors
	ErrAccountNotFound      = errors.New("account not found")
	ErrAccountAlreadyExists = errors.New("account already tracked")
	ErrAccountLimitReached  = errors.New("account limit reached for your plan")
	ErrAccountNotPublic     = errors.New("account is not public")

	// Scan errors
	ErrScanInProgress   = errors.New("scan already in progress")
	ErrScanFailed       = errors.New("scan failed")
	ErrRateLimited      = errors.New("rate limited by platform")
	ErrProxyUnavailable = errors.New("no proxies available")

	// Validation errors
	ErrInvalidUsername = errors.New("invalid username")
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidPassword = errors.New("password must be at least 8 characters")
	ErrInvalidPlatform = errors.New("unsupported platform")

	// Authorization errors
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
)

// DomainError wraps an error with additional context
type DomainError struct {
	Err     error
	Message string
	Code    string
}

func (e *DomainError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}

func (e *DomainError) Unwrap() error {
	return e.Err
}

// NewDomainError creates a new domain error
func NewDomainError(err error, message, code string) *DomainError {
	return &DomainError{
		Err:     err,
		Message: message,
		Code:    code,
	}
}
