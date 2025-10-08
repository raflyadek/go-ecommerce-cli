package handler

import (
	"fmt"
	"go-ecommerce-cli/internal/repository"

	"github.com/manifoldco/promptui"
)

type ReportHandler struct {
	OrderRepo repository.OrderRepository
}

func NewReportHandler(orderRepo repository.OrderRepository) *ReportHandler {
	return &ReportHandler{
		OrderRepo: orderRepo,
	}
}

func (h *ReportHandler) ShowReportMenu() {
	menu := []string{
		"Order Report (Completed Orders)",
		"Back to Dashboard",
	}

	for {

		prompt := promptui.Select{
			Label: "=== Report Menu (Minimal) ===",
			Items: menu,
			Templates: &promptui.SelectTemplates{
				Label:    "{{ . | bold }}",
				Active:   " {{ . | cyan | bold }}",
				Inactive: "  {{ . | white }}",
				Selected: " {{ . | green | bold }}",
			},
		}

		i, _, err := prompt.Run()
		if err != nil {
			fmt.Println("Prompt failed:", err)
			return
		}

		switch i {
		case 0:
			h.PrintCompletedOrderReport()
		case 1:
			return
		}
	}
}

func (h *ReportHandler) PrintCompletedOrderReport() {
	orders, err := h.OrderRepo.FindCompletedOrders()
	if err != nil {
		fmt.Println("Error retrieving completed orders:", err)
		return
	}

	fmt.Println("\n--- Order Report (Completed Orders) --- (Ctrl + C to return to dashboard)")

	fmt.Println("|---------|----------|------------|----------|----------------|")
	fmt.Printf("| %-7s | %-8s | %-10s | %-8s | %-14s |\n", "Order ID", "Customer", "Total (Rp)", "Status", "Order Date")
	fmt.Println("|---------|----------|------------|----------|----------------|")

	totalCompletedOrders := 0
	totalRevenue := 0.0

	for _, order := range orders {
		dateStr := order.CompletedDate.Format("2006-01-02")

		fmt.Printf("| %-7d | %-8s | %-10.2f | %-8s | %-14s |\n",
			order.ID,
			order.CustomerName,
			order.TotalAmount,
			order.StatusName, // Menggunakan StatusName
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
