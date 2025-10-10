package entity

type UserReport struct {
	UserID     int
	UserName   string
	OrderCount int
	TotalSpent float64
}

type UserOrderSummary struct {
	UserName        string
	TotalOrders     int
	PendingOrders   int
	CompletedOrders int
	CancelledOrders int
	TotalAmount     float64
	CompletedAmount float64
	PendingAmount   float64
	CancelledAmount float64
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