package dtos

type UpdateOrderDTO struct {
	ProductID string `json:"product_id" binding:"required,uuid"`
	Quantity  int    `json:"quantity" binding:"required,gt=0"`
}
