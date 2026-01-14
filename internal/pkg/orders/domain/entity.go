package domain

import (
	"github.com/QuangNV23062004/learning-go/internal/domain"
)

type Order struct {
	domain.BaseEntity
	ProductID string  `json:"product_id" gorm:"type:uuid;not null"`
	Quantity  int     `json:"quantity" gorm:"not null"`
	UserID    string  `json:"user_id" gorm:"type:uuid;not null"`
	Total     float64 `json:"total" gorm:"not null"`
}

func (o *Order) GetBaseEntity() *domain.BaseEntity {
	return &o.BaseEntity
}
