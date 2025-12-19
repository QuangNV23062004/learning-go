package domain

import "time"

type BaseEntity struct {
	ID string `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	IsDeleted bool      `json:"is_deleted" gorm:"default:false"`
	DeletedAt string    `json:"deleted_at" gorm:"default:null"`
}
