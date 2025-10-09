package handler

import (
	"fmt"
	"go-ecommerce-cli/internal/repository"
)

type OrderHandler struct {
	OrderRepo repository.OrderRepository
}

func NewOrderHandler(orderRepo repository.OrderRepository) *OrderHandler {
	return &OrderHandler{
		OrderRepo: orderRepo,
	}
}

func (h *OrderHandler) ShowCompletedOrders() {
	orders, err := h.OrderRepo.FindCompletedOrders()
	if err != nil {
		fmt.Println("Error retrieving completed orders:", err)
		return
	}

	if len(orders) == 0 {
		fmt.Println("No completed orders found.")
		return
	}

	fmt.Println("\n===== Order Report (Completed Orders) ===== (\"Ctrl+C\" to return to dashboard)")
	fmt.Println()
	fmt.Printf("%-8s %-12s %-12s %-10s %-14s\n", "Order ID", "Customer", "Total (Rp)", "Status", "Completed Date")

	totalOrders := 0
	var totalRevenue float64

	for _, order := range orders {
		dateStr := order.CompletedDate.Format("2006-01-02")
		fmt.Printf("%-8d %-12s %-12.2f %-10s %-14s\n",
			order.ID,
			order.CustomerName,
			float64(order.TotalAmount),
			order.StatusName,
			dateStr,
		)
		totalRevenue += float64(order.TotalAmount)
		totalOrders++
	}

	fmt.Printf("\n> Total Completed Orders: %d\n", totalOrders)
	fmt.Printf("> Total Revenue: Rp%.2f\n", totalRevenue)

	fmt.Print("\nPress ENTER to continue...")
	fmt.Scanln()
}
func (h *OrderHandler) ShowUserReport() {
	fmt.Printf("\n===== User Report ===== (\"Ctrl+C\" to return to dashboard)\n")
	fmt.Print("Search by ID: ")
	var userID int
	fmt.Scanln(&userID)

	orders, err := h.OrderRepo.FindUserOrders(userID)
	if err != nil {
		fmt.Println("Error retrieving user orders:", err)
		return
	}

	if len(orders) == 0 {
		fmt.Println("No orders found for this user.")
		fmt.Print("\nPress ENTER to continue...")
		fmt.Scanln()
		return
	}

	userName := orders[0].CustomerName
	fmt.Printf("User ID: %d Name: %s\n\n", userID, userName)
	fmt.Printf("%-8s %-10s %-12s\n", "Order ID", "Status", "Total (Rp)")

	var totalAmount float64
	for _, order := range orders {
		fmt.Printf("%-8d %-10s %-12.2f\n",
			order.ID,
			order.StatusName,
			float64(order.TotalAmount),
		)
		totalAmount += float64(order.TotalAmount)
	}

	fmt.Printf("\n> Summary by User\n")
	fmt.Printf("%s: %d Orders (Total Rp%.2f)\n", userName, len(orders), totalAmount)

	fmt.Print("\nPress ENTER to continue...")
	fmt.Scanln()
}

func (h *OrderHandler) ShowStockReport() {
	fmt.Printf("\n===== Stock Report (Daily Summary) ===== (\"Ctrl+C\" to return to dashboard)\n")
	fmt.Print("Input Date: ")
	var dateInput string
	fmt.Scanln(&dateInput)

	stockReports, err := h.OrderRepo.GetStockReport()
	if err != nil {
		fmt.Println("Error retrieving stock report:", err)
		return
	}

	fmt.Printf("%-10s %-20s %-9s %-9s %-13s\n", "Product ID", "Product Name", "Stock In", "Stock Out", "Current Stock")

	for _, stock := range stockReports {
		fmt.Printf("%-10d %-20s %-9d %-9d %-13d\n",
			stock.ProductID,
			stock.ProductName,
			0, // Stock In - tidak ada data movement
			0, // Stock Out - tidak ada data movement
			stock.CurrentStock,
		)
	}

	fmt.Print("\nPress ENTER to continue...")
	fmt.Scanln()
}
