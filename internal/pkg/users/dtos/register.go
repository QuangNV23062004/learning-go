package dtos

type RegisterDto struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	Username  string `json:"username" validate:"required"`
	Birthdate string `json:"birthdate" validate:"omitempty,datetime=2006-01-02"`
}
