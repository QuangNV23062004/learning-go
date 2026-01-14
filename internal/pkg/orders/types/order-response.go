package types

import (
	"github.com/QuangNV23062004/learning-go/internal/pkg/orders/domain"
)

type OrderResponse struct {
	// Core order fields
	domain.Order

	// Enriched from User
	Username string `json:"username,omitempty" gorm:"column:username"`
	Email    string `json:"email,omitempty" gorm:"column:email"`

	// Enriched from Product
	ProductName  string  `json:"product_name,omitempty" gorm:"column:product_name"`
	ProductPrice float64 `json:"product_price,omitempty" gorm:"column:product_price"`
}
