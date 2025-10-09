package handler

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
)

type ReportHandler struct {
	DB *sql.DB
}

func NewReportHandler(db *sql.DB) *ReportHandler {
	return &ReportHandler{DB: db}
}

func (h *ReportHandler) ShowUserReport() {
	fmt.Printf("\n===== User Report ===== (\"Ctrl+C\" to return to dashboard)\n")
	fmt.Print("Search by ID: ")
	var userID int
	fmt.Scanln(&userID)

	query := `
		SELECT 
			o.id,
			u.name AS customer_name,
			o.total_amount,
			s.status_name
		FROM orders o
		JOIN users u ON o.user_id = u.id
		JOIN status_order s ON o.status_id = s.id
		WHERE o.user_id = $1
		ORDER BY o.created_at DESC
	`

	rows, err := h.DB.Query(query, userID)
	if err != nil {
		fmt.Println("Error retrieving user orders:", err)
		return
	}
	defer rows.Close()

	var orders []struct {
		ID           int
		CustomerName string
		TotalAmount  int
		StatusName   string
	}

	for rows.Next() {
		var order struct {
			ID           int
			CustomerName string
			TotalAmount  int
			StatusName   string
		}
		err := rows.Scan(&order.ID, &order.CustomerName, &order.TotalAmount, &order.StatusName)
		if err != nil {
			fmt.Println("Error scanning order:", err)
			return
		}
		orders = append(orders, order)
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
	query := `
		SELECT 
			o.id,
			u.name AS customer_name,
			o.total_amount,
			s.status_name,
			o.updated_at
		FROM orders o
		JOIN users u ON o.user_id = u.id
		JOIN status_order s ON o.status_id = s.id
		WHERE s.status_name = 'completed'
		ORDER BY o.updated_at DESC
	`

	rows, err := h.DB.Query(query)
	if err != nil {
		fmt.Println("Error retrieving completed orders:", err)
		return
	}
	defer rows.Close()

	fmt.Println("\n===== Order Report (Completed Orders) ===== (\"Ctrl+C\" to return to dashboard)")
	fmt.Println()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Order ID", "Customer", "Total (Rp)", "Status", "Completed Date"})

	totalOrders := 0
	var totalRevenue float64

	for rows.Next() {
		var orderID int
		var customerName string
		var totalAmount int
		var statusName string
		var updatedAt time.Time

		err := rows.Scan(&orderID, &customerName, &totalAmount, &statusName, &updatedAt)
		if err != nil {
			fmt.Println("Error scanning order:", err)
			return
		}

		dateStr := updatedAt.Format("2006-01-02")
		row := []string{
			fmt.Sprintf("%d", orderID),
			customerName,
			fmt.Sprintf("%.2f", float64(totalAmount)),
			statusName,
			dateStr,
		}
		table.Append(row)
		totalRevenue += float64(totalAmount)
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

	query := `
		SELECT 
			p.id,
			p.name,
			p.stock,
			COALESCE(SUM(oi.quantity), 0) as total_sold
		FROM products p
		LEFT JOIN order_items oi ON p.id = oi.product_id
		LEFT JOIN orders o ON oi.order_id = o.id
		WHERE o.status_id IN (2,3,4) OR o.status_id IS NULL
		GROUP BY p.id, p.name, p.stock
		ORDER BY p.id
	`

	rows, err := h.DB.Query(query)
	if err != nil {
		fmt.Println("Error retrieving stock report:", err)
		return
	}
	defer rows.Close()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Product ID", "Product Name", "Stock In", "Stock Out", "Current Stock"})

	for rows.Next() {
		var productID int
		var productName string
		var currentStock int
		var totalSold int

		err := rows.Scan(&productID, &productName, &currentStock, &totalSold)
		if err != nil {
			fmt.Println("Error scanning stock:", err)
			return
		}

		stockIn := currentStock + totalSold
		row := []string{
			fmt.Sprintf("%d", productID),
			productName,
			fmt.Sprintf("%d", stockIn),
			fmt.Sprintf("%d", totalSold),
			fmt.Sprintf("%d", currentStock),
		}
		table.Append(row)
	}
	table.Render()

	fmt.Print("\nPress ENTER to continue...")
	fmt.Scanln()
}
