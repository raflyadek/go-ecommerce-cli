package entity

import (
	"time"
)

type Order struct {
	Address 	string
	ID          int
	UserID      int
	StatusID    int
	TotalAmount float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Items       []OrderItem
}

type OrderItem struct {
	ProductName string
	OrderID   int
	ProductID int
	Quantity  int
	Price     float64
}

type OrderSummary struct {
	Status    string
	OrderDate string
	ID        int
	Total     float64
	ItemCount int
}

type OrderItemDetails struct {
	ProductName string
	Quantity    int
	Price       float64
}

type UserOrder struct {
	ID        int
	Status    string
	Total     float64
	ItemCount int
	OrderDate time.Time
}
