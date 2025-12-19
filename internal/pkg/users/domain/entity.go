package domain

import "learning-go/internal/domain"

type User struct {
	Email     string `json:"email" gorm:"uniqueIndex"`
	Password  string `json:"password" gorm:"not null"`
	Username  string `json:"username" gorm:"not null"`
	Role      string `json:"role" gorm:"not null;default:'user'"`
	Birthdate string `json:"birthdate" gorm:"not null"`
	domain.BaseEntity
}

func (u *User) GetBaseEntity() *domain.BaseEntity {
	return &u.BaseEntity
}
