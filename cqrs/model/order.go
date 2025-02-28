package model

// Order
type Order struct {
	ID     string
	Status string
	Items  []OrderItem
}

// OrderItem
type OrderItem struct {
	ID       string
	Name     string
	Quantity int
	Price    float64
}
