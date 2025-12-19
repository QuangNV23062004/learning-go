package utils

import (
	"learning-go/internal/types"
)

func Success[T any](data T, status int) types.Response {
	return types.Response{
		Status:  status,
		Success: true,
		Data:    data,
	}
}

func Error(message string, status int) types.Response {
	return types.Response{
		Status:  status,
		Success: false,
		Error:   message,
	}
}
