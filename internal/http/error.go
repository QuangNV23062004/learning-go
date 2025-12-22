package http

import (
	"errors"
	productDomain "learning-go/internal/pkg/products/domain"
	userDomain "learning-go/internal/pkg/users/domain"
)

// globalish errors
var (
	ErrInvalidBody         = errors.New("invalid request body")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden           = errors.New("forbidden")
	ErrInvalidQuery        = errors.New("invalid query parameters")
	ErrMissingRefreshToken = errors.New("missing refresh token in cookies")
)

func GetStatusCode(err error) int {
	switch {
	case errors.Is(err, userDomain.ErrUserAlreadyExists):
		return 400
	case errors.Is(err, userDomain.ErrInvalidCredentials):
		return 400
	case errors.Is(err, ErrUnauthorized):
		return 401
	case errors.Is(err, ErrForbidden):
		return 403
	case errors.Is(err, ErrInvalidBody):
		return 400
	case errors.Is(err, ErrMissingRefreshToken):
		return 400
	case errors.Is(err, productDomain.ErrProductNotFound):
		return 404
	case errors.Is(err, productDomain.ErrUserNotFound):
		return 404
	default:
		return 500
	}
}
