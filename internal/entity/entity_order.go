package entity

import (
	"time"
)

type Order struct {
	ID            int
	UserID        int
	TotalAmount   int
	StatusName    string
	CompletedDate time.Time

	CustomerName string
	Items        []OrderItem
}

type OrderItem struct {
	ID          int
	OrderID     int
	ProductID   int
	ProductName string
	Quantity    int
	Price       int
}
