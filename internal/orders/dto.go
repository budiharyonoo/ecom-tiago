package orders

type CreateOrderRequest struct {
	CustomerID int64          `json:"customer_id" validate:"required,numeric"`
	TotalPrice int32          `json:"total_price" validate:"required,numeric"`
	Items      []OrderItemDTO `json:"items" validate:"required,dive"`
}

type OrderItemDTO struct {
	ProductID   int64  `json:"product_id" validate:"required,numeric"`
	ProductName string `json:"product_name" validate:"required,alphanumspace,max=255"`
	Quantity    int32  `json:"quantity" validate:"required,numeric"`
	Price       int32  `json:"price" validate:"required,numeric"`
}
