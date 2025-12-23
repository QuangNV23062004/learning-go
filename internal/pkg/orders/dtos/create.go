package dtos

type CreateOrderDTO struct {
	ProductID string `json:"product_id" binding:"required,uuid"`
	Quantity  int    `json:"quantity" binding:"required,gt=0"`
}
