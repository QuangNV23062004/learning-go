package domain

import (
	"learning-go/internal/domain"
)

type Product struct {
	Name   string  `json:"name" gorm:"not null"`
	Price  float64 `json:"price" gorm:"not null"`
	Stock  int     `json:"stock" gorm:"not null"`
	UserID string  `json:"user_id" gorm:"type:uuid;not null"`
	domain.BaseEntity
}

func (p *Product) GetBaseEntity() *domain.BaseEntity {
	return &p.BaseEntity
}
