package entity

type UserReport struct {
	UserID     int
	UserName   string
	OrderCount int
	TotalSpent float64
}

type StockReport struct {
	ProductID    int
	ProductName  string
	StockIn      int
	StockOut     int
	CurrentStock int
}

type BestSellingProduct struct {
	ProductID    int
	ProductName  string
	TotalSold    int
	TotalRevenue float64
}

type ReportSummary struct {
	TotalUsers    int
	TotalOrders   int
	TotalRevenue  float64
	AvgOrderValue float64
}