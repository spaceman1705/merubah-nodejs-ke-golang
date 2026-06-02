package dto

type CreateOrderRequest struct {
	UserID          string      `json:"-"`
	UserAddressID   string      `json:"userAddressId" validate:"required"`
	StoreID         string      `json:"storeId" validate:"required"`
	Items           []OrderItem `json:"items" validate:"required,min=1,dive"`
	Subtotal        float64     `json:"subtotal" validate:"required,min=0"`
	ShippingCost    float64     `json:"shippingCost" validate:"required,min=0"`
	DiscountAmount  float64     `json:"discountAmount" validate:"min=0"`
	TotalAmount     float64     `json:"totalAmount" validate:"required,min=0"`
	VoucherCodeUsed *string     `json:"voucherCodeUsed"`
	PaymentMethod   string      `json:"paymentMethod" validate:"required,oneof=TRANSFER COD"`
}

type OrderItem struct {
	ProductID string  `json:"productId" validate:"required"`
	Quantity  int     `json:"quantity" validate:"required,min=1"`
	Price     float64 `json:"price" validate:"required,min=0"`
}
