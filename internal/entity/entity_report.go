package entity

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
