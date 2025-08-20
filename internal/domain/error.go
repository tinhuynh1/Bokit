package domain

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrInvalidMethod      = errors.New("invalid login method")
)
