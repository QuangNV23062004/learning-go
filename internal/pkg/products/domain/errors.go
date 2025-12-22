package domain

import "errors"

var (
	ErrProductNotFound = errors.New("product not found")
	ErrUserNotFound    = errors.New("user not found")
)
