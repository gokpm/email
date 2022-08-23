package email

import "errors"

var (
	ErrInvalidEmail    = errors.New("invalid email address")
	ErrDisposableEmail = errors.New("disposable email address")
	ErrNoMXRecords     = errors.New("mx record not found")
)
