package domain

import "time"

type Order struct {
	ID        string    `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ProductID string    `json:"product_id" gorm:"type:uuid;not null"`
	Quantity  int       `json:"quantity" gorm:"not null"`
	UserID    string    `json:"user_id" gorm:"type:uuid;not null"`
	Total     float64   `json:"total" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	IsDeleted bool      `json:"is_deleted" gorm:"default:false"`
	DeletedAt string    `json:"deleted_at" gorm:"default:null"`
}
