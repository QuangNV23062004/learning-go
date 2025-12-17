package domain

import "time"

type Product struct {
	ID        string    `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string    `json:"username" gorm:"not null"`
	Price     float64   `json:"price" gorm:"not null"`
	Stock     int       `json:"stock" gorm:"not null"`
	UserID    string    `json:"user_id" gorm:"type:uuid;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	IsDeleted bool      `json:"is_deleted" gorm:"default:false"`
	DeletedAt string    `json:"deleted_at" gorm:"default:null"`
}
