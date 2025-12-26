package orders

type Order struct {
	CustomerID int64
	TotalPrice int32
	Items      []OrderItem
}

type OrderItem struct {
	ProductID   int64
	ProductName string
	Quantity    int32
	Price       int32
}
