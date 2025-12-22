package dtos

type CreateOrderDTO struct {
	ProductID string  `json:"product_id" binding:"required,uuid"`
	Quantity  int     `json:"quantity" binding:"required,gt=0"`
	UserID    string  `json:"user_id" binding:"required,uuid"`
	Total     float64 `json:"total" binding:"required,gt=0"`
}
