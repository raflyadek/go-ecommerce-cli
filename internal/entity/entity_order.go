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

type UserReport struct {
	UserID     int
	UserName   string
	OrderCount int
	TotalSpent int
}

type StockReport struct {
	ProductID    int
	ProductName  string
	StockIn      int
	StockOut     int
	CurrentStock int
}

type OrderItem struct {
	ID          int
	OrderID     int
	ProductID   int
	ProductName string
	Quantity    int
	Price       int
}
