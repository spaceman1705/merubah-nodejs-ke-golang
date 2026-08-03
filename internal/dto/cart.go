package dto

type AddToCartRequest struct {
	ProductID string `json:"productId" validate:"required"`
	StoreID   string `json:"storeId" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,min=1"`
}

type UpdateCartRequest struct {
	Quantity int `json:"quantity" validate:"required,min=1"`
}
