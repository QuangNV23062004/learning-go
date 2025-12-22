package dtos

type CreateProductDTO struct {
	Name   string  `json:"name" binding:"required"`
	Price  float64 `json:"price" binding:"required,gt=0"`
	Stock  int     `json:"stock" binding:"required,gte=0"`
	UserID string  `json:"user_id" binding:"required,uuid"`
}
