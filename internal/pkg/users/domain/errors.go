package domain

import "errors"

var (
	ErrUserAlreadyExists        = errors.New("user already exists")
	ErrInvalidCredentials       = errors.New("invalid credentials")
	ErrFailedToRenderHTML       = errors.New("failed to render HTML")
	ErrInvalidVerificationToken = errors.New("invalid or expired verification token")
	ErrUserNotFound             = errors.New("user not found")
)
