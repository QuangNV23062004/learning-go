package domain

import "errors"

var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrProductNotFound    = errors.New("product not found")
	ErrUserNotFound       = errors.New("user not found")
	ErrNotAllowed         = errors.New("operation not allowed")
	ErrInsufficientStock  = errors.New("insufficient stock for the product")
	ErrOldProductNotFound = errors.New("old product not found")
)
