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
		totalRevenue += order.TotalAmount
		totalCompletedOrders++
	}

	fmt.Println("|---------|----------|------------|----------|----------------|")

	fmt.Printf("\n*** Total Completed Orders: %d\n", totalCompletedOrders)
	fmt.Printf("*** Total Revenue: Rp%.2f\n", totalRevenue)

	fmt.Print("\nPress ENTER to continue...")
	fmt.Scanln()
}
