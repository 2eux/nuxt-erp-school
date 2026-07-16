package domain

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound          = errors.New("resource not found")
	ErrDuplicate         = errors.New("resource already exists")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrInvalidInput      = errors.New("invalid input")
	ErrInternalServer    = errors.New("internal server error")
	ErrTokenExpired      = errors.New("token expired")
	ErrTokenInvalid      = errors.New("token invalid")
	ErrPasswordMismatch  = errors.New("password mismatch")
	ErrUserInactive      = errors.New("user is inactive")
	ErrSchoolNotActive   = errors.New("school is not active")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrSessionExpired    = errors.New("session expired")
	ErrPaymentFailed     = errors.New("payment failed")
	ErrInvoiceNotFound   = errors.New("invoice not found")
	ErrQuotaExceeded     = errors.New("quota exceeded")
	ErrFileTooLarge      = errors.New("file too large")
	ErrUnsupportedFile   = errors.New("unsupported file type")
	ErrRateLimited       = errors.New("rate limit exceeded")
	ErrValidation        = errors.New("validation failed")
)

type DomainError struct {
	Code    string
	Message string
	Err     error
}

func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *DomainError) Unwrap() error {
	return e.Err
}

func NewNotFoundError(entity string, id string) *DomainError {
	return &DomainError{
		Code:    "NOT_FOUND",
		Message: fmt.Sprintf("%s with id %s not found", entity, id),
		Err:     ErrNotFound,
	}
}

func NewDuplicateError(entity string, field string) *DomainError {
	return &DomainError{
		Code:    "DUPLICATE",
		Message: fmt.Sprintf("%s with %s already exists", entity, field),
		Err:     ErrDuplicate,
	}
}

func NewForbiddenError(msg string) *DomainError {
	return &DomainError{
		Code:    "FORBIDDEN",
		Message: msg,
		Err:     ErrForbidden,
	}
}

func NewUnauthorizedError(msg string) *DomainError {
	return &DomainError{
		Code:    "UNAUTHORIZED",
		Message: msg,
		Err:     ErrUnauthorized,
	}
}

func NewInvalidInputError(msg string) *DomainError {
	return &DomainError{
		Code:    "INVALID_INPUT",
		Message: msg,
		Err:     ErrInvalidInput,
	}
}

func NewInternalError(msg string, err error) *DomainError {
	return &DomainError{
		Code:    "INTERNAL_ERROR",
		Message: msg,
		Err:     fmt.Errorf("%w: %v", ErrInternalServer, err),
	}
}

func NewValidationError(msg string) *DomainError {
	return &DomainError{
		Code:    "VALIDATION_ERROR",
		Message: msg,
		Err:     ErrValidation,
	}
}

func NewRateLimitError() *DomainError {
	return &DomainError{
		Code:    "RATE_LIMIT",
		Message: "too many requests, please try again later",
		Err:     ErrRateLimited,
	}
}

func NewTokenExpiredError() *DomainError {
	return &DomainError{
		Code:    "TOKEN_EXPIRED",
		Message: "authentication token has expired",
		Err:     ErrTokenExpired,
	}
}

func NewTokenInvalidError() *DomainError {
	return &DomainError{
		Code:    "TOKEN_INVALID",
		Message: "authentication token is invalid",
		Err:     ErrTokenInvalid,
	}
}
