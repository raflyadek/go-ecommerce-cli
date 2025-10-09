package handler

import (
	"database/sql"
	"fmt"
	"go-ecommerce-cli/internal/repository"
	"os"

	"github.com/olekukonko/tablewriter"
)

type ReportHandler struct {
	ReportRepo repository.ReportRepository
}

func NewReportHandler(db *sql.DB) *ReportHandler {
	reportRepo := repository.NewReportRepo(db)
	return &ReportHandler{ReportRepo: reportRepo}
}

func (h *ReportHandler) ShowUserReport() {
	fmt.Printf("\n===== User Report ===== (\"Ctrl+C\" to return to dashboard)\n")
	fmt.Print("Search by ID: ")
	var userID int
	fmt.Scanln(&userID)

	orders, err := h.ReportRepo.FindUserOrders(userID)
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

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Order ID", "Status", "Total (Rp)"})

	var totalAmount float64
	for _, order := range orders {
		row := []string{
			fmt.Sprintf("%d", order.ID),
			order.StatusName,
			fmt.Sprintf("%.2f", float64(order.TotalAmount)),
		}
		table.Append(row)
		totalAmount += float64(order.TotalAmount)
	}

	table.Render()
	fmt.Printf("\n> Summary by User\n")
	fmt.Printf("%s: %d Orders (Total Rp%.2f)\n", userName, len(orders), totalAmount)

	fmt.Print("\nPress ENTER to continue...")
	fmt.Scanln()
}

func (h *ReportHandler) ShowCompletedOrders() {
	orders, err := h.ReportRepo.FindCompletedOrders()
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

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Order ID", "Customer", "Total (Rp)", "Status", "Completed Date"})

	totalOrders := 0
	var totalRevenue float64

	for _, order := range orders {
		dateStr := order.CompletedDate.Format("2006-01-02")
		row := []string{
			fmt.Sprintf("%d", order.ID),
			order.CustomerName,
			fmt.Sprintf("%.2f", float64(order.TotalAmount)),
			order.StatusName,
			dateStr,
		}
		table.Append(row)
		totalRevenue += float64(order.TotalAmount)
		totalOrders++
	}

	table.Render()
	fmt.Printf("\n> Total Completed Orders: %d\n", totalOrders)
	fmt.Printf("> Total Revenue: Rp%.2f\n", totalRevenue)

	fmt.Print("\nPress ENTER to continue...")
	fmt.Scanln()
}

func (h *ReportHandler) ShowStockReport() {
	fmt.Printf("\n===== Stock Report (Daily Summary) ===== (\"Ctrl+C\" to return to dashboard)\n")
	fmt.Print("Input Date: ")
	var dateInput string
	fmt.Scanln(&dateInput)

	stockReports, err := h.ReportRepo.GetStockReport()
	if err != nil {
		fmt.Println("Error retrieving stock report:", err)
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Product ID", "Product Name", "Stock In", "Stock Out", "Current Stock"})

	for _, stock := range stockReports {
		row := []string{
			fmt.Sprintf("%d", stock.ProductID),
			stock.ProductName,
			fmt.Sprintf("%d", stock.StockIn),
			fmt.Sprintf("%d", stock.StockOut),
			fmt.Sprintf("%d", stock.CurrentStock),
		}
		table.Append(row)
	}
	table.Render()

	fmt.Print("\nPress ENTER to continue...")
	fmt.Scanln()
}
