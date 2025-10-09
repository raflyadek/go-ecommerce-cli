package entity

import (
	"time"
)

type Order struct {
	ID            int
	UserID        int
	TotalAmount   float64
	StatusName    string
	CompletedDate time.Time
	Address       string

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
