package dtos

type UpdateProductDTO struct {
	Name  string  `json:"name" binding:"ommitempty"`
	Price float64 `json:"price" binding:"omitempty,gt=0"`
	Stock int     `json:"stock" binding:"omitempty,gte=0"`
}
