package dtos

type UpdateUserDto struct {
	Username  string `json:"username" validate:"required"`
	Birthdate string `json:"birthdate" validate:"required,datetime=2006-01-02"`
}
