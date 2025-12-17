package domain

import "time"

type User struct {
	ID        string    `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email     string    `json:"email" gorm:"uniqueIndex"`
	Password  string    `json:"password" gorm:"not null"`
	Username  string    `json:"username" gorm:"not null"`
	Role      string    `json:"role" gorm:"not null;default:'user'"`
	Birthdate string    `json:"birthdate"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	IsDeleted bool      `json:"is_deleted" gorm:"default:false"`
	DeletedAt string    `json:"deleted_at" gorm:"default:null"`
}
