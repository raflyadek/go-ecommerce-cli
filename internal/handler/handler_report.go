package handler

import (
	"database/sql"
	"fmt"
	"go-ecommerce-cli/internal/entity"
	"go-ecommerce-cli/internal/repository"
	"os"
	"sync"
	"time"

	"github.com/manifoldco/promptui"
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

	searchOptions := []string{"Search by ID", "Search by Name"}
	prompt := promptui.Select{
		Label: "Choose Search Option",
		Items: searchOptions,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . | cyan | bold }}",
			Active:   "{{ . | green | bold }}",
			Inactive: "  {{ . | white }}",
			Selected: "{{ . | bold }}",
		},
	}

	choice, _, err := prompt.Run()
	if err != nil {
		fmt.Println("Prompt failed:", err)
		return
	}

	var orders []entity.Order
	var searchKey string

	switch choice {
	case 0:
		fmt.Print("Enter User ID: ")
		var userID int
		fmt.Scanln(&userID)
		orders, err = h.ReportRepo.FindUserOrders(userID)
		searchKey = fmt.Sprintf("ID: %d", userID)
	case 1:
		fmt.Print("Enter User Name: ")
		var userName string
		fmt.Scanln(&userName)
		orders, err = h.ReportRepo.FindUserOrdersByName(userName)
		searchKey = fmt.Sprintf("Name: %s", userName)
	}

	if err != nil {
		fmt.Println("Error retrieving user orders:", err)
		return
	}

	if len(orders) == 0 {
		fmt.Printf("No orders found for %s.\n", searchKey)
		fmt.Print("\nPress ENTER to continue...")
		fmt.Scanln()
		return
	}

	userName := orders[0].CustomerName
	fmt.Printf("\nSearch Result for %s - User: %s\n\n", searchKey, userName)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Order ID", "Status", "Total (Rp)", "Address"})

	var totalAmount float64
	for _, order := range orders {
		row := []string{
			fmt.Sprintf("%d", order.ID),
			order.StatusName,
			fmt.Sprintf("%.2f", order.TotalAmount),
			order.Address,
		}
		table.Append(row)
		totalAmount += order.TotalAmount
	}

	table.Render()
	fmt.Printf("\n> Summary\n")
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

	fmt.Println("\n===== Order Report (Completed Orders) =====")
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
			fmt.Sprintf("%.2f", order.TotalAmount),
			order.StatusName,
			dateStr,
		}
		table.Append(row)
		totalRevenue += order.TotalAmount
		totalOrders++
	}

	table.Render()
	fmt.Printf("\n> Total Completed Orders: %d\n", totalOrders)
	fmt.Printf("> Total Revenue: Rp%.2f\n", totalRevenue)

	fmt.Print("\nPress ENTER to continue...")
	fmt.Scanln()
}

func (h *ReportHandler) ShowStockReport() {
	fmt.Printf("\n===== Stock Report (Daily Summary) =====\n")
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
func (h *ReportHandler) ShowBestSellingProducts() {
	products, err := h.ReportRepo.GetBestSellingProducts()
	if err != nil {
		fmt.Println("Error retrieving best selling products:", err)
		return
	}

	if len(products) == 0 {
		fmt.Println("No sales data found.")
		return
	}

	fmt.Println("\n===== Best Selling Products =====")
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Product ID", "Product Name", "Total Sold", "Total Revenue"})

	for _, product := range products {
		row := []string{
			fmt.Sprintf("%d", product.ProductID),
			product.ProductName,
			fmt.Sprintf("%d", product.TotalSold),
			fmt.Sprintf("%.2f", product.TotalRevenue),
		}
		table.Append(row)
	}
	table.Render()
}

func (h *ReportHandler) ShowReportSummary() {
	fmt.Println("\n===== Report Summary ===== (Loading...)")
	start := time.Now()

	// Channel untuk menerima hasil dari goroutine
	type result struct {
		name  string
		value interface{}
		err   error
	}

	resultChan := make(chan result, 4)
	var wg sync.WaitGroup

	// Goroutine 1: Total Users
	wg.Add(1)
	go func() {
		defer wg.Done()
		count, err := h.ReportRepo.GetTotalUsers()
		resultChan <- result{"Total Users", count, err}
	}()

	// Goroutine 2: Total Orders
	wg.Add(1)
	go func() {
		defer wg.Done()
		count, err := h.ReportRepo.GetTotalOrders()
		resultChan <- result{"Total Orders", count, err}
	}()

	// Goroutine 3: Total Revenue
	wg.Add(1)
	go func() {
		defer wg.Done()
		revenue, err := h.ReportRepo.GetTotalRevenue()
		resultChan <- result{"Total Revenue", revenue, err}
	}()

	// Goroutine 4: Average Order Value
	wg.Add(1)
	go func() {
		defer wg.Done()
		avg, err := h.ReportRepo.GetAvgOrderValue()
		resultChan <- result{"Average Order Value", avg, err}
	}()

	// Tutup channel setelah semua goroutine selesai
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Kumpulkan hasil dari semua goroutine
	results := make(map[string]interface{})
	for res := range resultChan {
		if res.err != nil {
			fmt.Printf("Error retrieving %s: %v\n", res.name, res.err)
			return
		}
		results[res.name] = res.value
	}

	elapsed := time.Since(start)
	fmt.Printf("\n===== Report Summary ===== (Loaded in %v)\n", elapsed)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Metric", "Value"})

	// Format hasil sesuai tipe data
	rows := [][]string{
		{"Total Users", fmt.Sprintf("%d", results["Total Users"].(int))},
		{"Total Orders", fmt.Sprintf("%d", results["Total Orders"].(int))},
		{"Total Revenue", fmt.Sprintf("Rp %.2f", results["Total Revenue"].(float64))},
		{"Average Order Value", fmt.Sprintf("Rp %.2f", results["Average Order Value"].(float64))},
	}

	for _, row := range rows {
		table.Append(row)
	}
	table.Render()
}
