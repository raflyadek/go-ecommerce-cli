package handler

import (
	"bufio"
	"database/sql"
	"fmt"
	"go-ecommerce-cli/internal/entity"
	"go-ecommerce-cli/internal/repository"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
)

type OrderHandler struct {
	OrderRepo repository.OrderRepository
	ProductRepo *repository.ProductRepository
	ProductHandler *ProductHandler
	DB        *sql.DB
}

func NewOrderHandler(orderRepo repository.OrderRepository, productRepo *repository.ProductRepository, productHandler *ProductHandler, db *sql.DB,) *OrderHandler {
	return &OrderHandler{
		OrderRepo:      orderRepo,
		ProductRepo:    productRepo,
		ProductHandler: productHandler,
		DB:             db,
	}
}

// CreateOrderCLI handles CLI input for creating orders
func (h *OrderHandler) CreateOrderCLI(userID int) {
	reader := bufio.NewReader(os.Stdin)
	var items []entity.OrderItem
	totalAmount := 0.0

	fmt.Println("\n=== Create Order ===")

	for {
		fmt.Print("Input Item ID or Name: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		product, err := h.findProductByIDOrName(input)
		if err != nil {
			fmt.Println("Product not found:", err)
			continue
		}

		fmt.Print("Input Quantity: ")
		qtyStr, _ := reader.ReadString('\n')
		qtyStr = strings.TrimSpace(qtyStr)
		qty, err := strconv.Atoi(qtyStr)
		if err != nil || qty <= 0 {
			fmt.Println("Invalid quantity.")
			continue
		}

		// Cek stock terlebih dahulu
		if product.Stock < qty {
			fmt.Printf("Insufficient stock for %s. Available: %d\n", product.Name, product.Stock)
			continue
		}

		item := entity.OrderItem{
			ProductID: product.ID,
			Quantity:  qty,
			Price:     product.Price * float64(qty),
		}
		items = append(items, item)
		totalAmount += item.Price

		fmt.Printf("%d %s successfully added with price %.2f.\n", qty, product.Name, item.Price)

		fmt.Print("Do you want to add another item? (y/n): ")
		again, _ := reader.ReadString('\n')
		again = strings.TrimSpace(again)
		if strings.ToLower(again) != "y" {
			break
		}
	}

	fmt.Print("Input your address: ")
	address, _ := reader.ReadString('\n')
	address = strings.TrimSpace(address)

	order := entity.Order{
		UserID:      userID,
		TotalAmount: totalAmount,
		Items:       items,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Address:     address,
	}

	// Insert order ke DB
	orderID, err := h.OrderRepo.CreateOrder(order)
	if err != nil {
		fmt.Println("Failed to create order:", err)
		return
	}

	// Kurangi stock tiap produk
	for _, item := range order.Items {
		err := h.ProductHandler.ProductRepo.ReduceStock(item.ProductID, item.Quantity)
		if err != nil {
			fmt.Println("Failed to reduce stock:", err)
			return
		}
	}

	fmt.Printf("Your order ID %d has been created with total price %.2f\n", orderID, totalAmount)
	fmt.Printf("Delivery Address: %s\n", address)
}

// findProductByIDOrName searches a product by ID or Name in the database
func (h *OrderHandler) findProductByIDOrName(input string) (entity.Product, error) {
	var product entity.Product

	id, err := strconv.Atoi(input)
	if err == nil {
		query := `SELECT id, name, description, price, stock FROM products WHERE id = $1 LIMIT 1`
		row := h.DB.QueryRow(query, id)
		err := row.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock)
		if err != nil {
			return product, fmt.Errorf("product not found with ID %d", id)
		}
		return product, nil
	}

	query := `SELECT id, name, description, price, stock FROM products WHERE name ILIKE $1 LIMIT 1`
	row := h.DB.QueryRow(query, "%"+input+"%")
	err = row.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock)
	if err != nil {
		return product, fmt.Errorf("product not found with name '%s'", input)
	}

	return product, nil
}

func (h *OrderHandler) AllOrdersCLI() {
    orders, err := h.OrderRepo.GetAllOrders()
    if err != nil {
        fmt.Println("Failed to fetch orders:", err)
        return
    }

    if len(orders) == 0 {
        fmt.Println("You have no orders yet.")
        return
    }

    fmt.Println("\n=== All Orders ===")
    table := tablewriter.NewWriter(os.Stdout)
    table.SetHeader([]string{"ID", "Status", "Total (Rp)", "Items", "Order Date"})

    for _, o := range orders {
        row := []string{
            fmt.Sprintf("%d", o.ID),
            o.Status,
            fmt.Sprintf("%.2f", o.Total),
            fmt.Sprintf("%d", o.ItemCount),
            o.OrderDate,
        }
        table.Append(row)
    }

    table.Render()
}

func (h *OrderHandler) HistoryOrdersCLI(userID int) {
    orders, err := h.OrderRepo.GetHistoryOrders(userID)
    if err != nil {
        fmt.Println("Failed to fetch orders:", err)
        return
    }

    if len(orders) == 0 {
        fmt.Println("You have no orders yet.")
        return
    }

    fmt.Println("\n=== Histories ===")
    table := tablewriter.NewWriter(os.Stdout)
    table.SetHeader([]string{"ID", "Status", "Total (Rp)", "Items", "Order Date"})

    for _, o := range orders {
        row := []string{
            fmt.Sprintf("%d", o.ID),
            o.Status,
            fmt.Sprintf("%.2f", o.Total),
            fmt.Sprintf("%d", o.ItemCount),
            o.OrderDate,
        }
        table.Append(row)
    }

    table.Render()

    fmt.Print("View order details with ID: ")
    var orderID int
    fmt.Scanln(&orderID)

    items, order, err := h.OrderRepo.GetOrderDetails(orderID)
    if err != nil {
        fmt.Println("Failed to fetch order details:", err)
        return
    }

    fmt.Printf("\nOrder ID: %d, Total: %.2f\n", order.ID, order.TotalAmount)
    fmt.Println("Items:")

    itemTable := tablewriter.NewWriter(os.Stdout)
    itemTable.SetHeader([]string{"Product Name", "Quantity", "Price (Rp)"})

    for _, item := range items {
        row := []string{
            item.ProductName,
            fmt.Sprintf("%d", item.Quantity),
            fmt.Sprintf("%.2f", item.Price),
        }
        itemTable.Append(row)
    }

    itemTable.Render()
}

func (h *OrderHandler) MyOrdersCLI(userID int) {
    orders, err := h.OrderRepo.GetUserOrders(userID)
    if err != nil {
        fmt.Println("Failed to fetch orders:", err)
        return
    }

    if len(orders) == 0 {
        fmt.Println("You have no orders yet.")
        return
    }

    fmt.Println("\n=== My Orders ===")
    table := tablewriter.NewWriter(os.Stdout)
    table.SetHeader([]string{"ID", "Status", "Total (Rp)", "Items", "Order Date"})

    for _, o := range orders {
        row := []string{
            fmt.Sprintf("%d", o.ID),
            o.Status,
            fmt.Sprintf("%.2f", o.Total),
            fmt.Sprintf("%d", o.ItemCount),
            o.OrderDate,
        }
        table.Append(row)
    }

    table.Render()

    fmt.Print("View order details with ID: ")
    var orderID int
    fmt.Scanln(&orderID)

    items, order, err := h.OrderRepo.GetOrderDetails(orderID)
    if err != nil {
        fmt.Println("Failed to fetch order details:", err)
        return
    }

    fmt.Printf("\nOrder ID: %d, Total: %.2f\n", order.ID, order.TotalAmount)
    fmt.Println("Items:")

    itemTable := tablewriter.NewWriter(os.Stdout)
    itemTable.SetHeader([]string{"Product Name", "Quantity", "Price (Rp)"})

    for _, item := range items {
        row := []string{
            item.ProductName,
            fmt.Sprintf("%d", item.Quantity),
            fmt.Sprintf("%.2f", item.Price),
        }
        itemTable.Append(row)
    }

    itemTable.Render()
}

// UpdateOrderStatusCLI allows admin/staff to update the status of an order
func (h *OrderHandler) UpdateOrderStatusCLI() {
    reader := bufio.NewReader(os.Stdin)
    fmt.Println("\n=== Update Order Status ===")

    fmt.Print("Enter Order ID: ")
    idStr, _ := reader.ReadString('\n')
    idStr = strings.TrimSpace(idStr)
    orderID, err := strconv.Atoi(idStr)
    if err != nil {
        fmt.Println("Invalid Order ID")
        return
    }

    fmt.Print("Enter new status (Pending/Shipped/Completed/Paid/Cancelled): ")
    status, _ := reader.ReadString('\n')
    status = strings.TrimSpace(status)

    err = h.OrderRepo.UpdateOrderStatus(orderID, status)
    if err != nil {
        fmt.Println("Failed to update order status:", err)
        return
    }

    fmt.Printf("Order #%d status updated successfully!\n", orderID)
}

func (h *OrderHandler) ViewOrderDetailsCLI() {
    reader := bufio.NewReader(os.Stdin)

    fmt.Println("\n=== View Order Details ===")
    fmt.Println("=== (“Ctrl + C” to return to dashboard) ===")

    fmt.Print("Input Order ID: ")
    idStr, _ := reader.ReadString('\n')
    idStr = strings.TrimSpace(idStr)
    orderID, err := strconv.Atoi(idStr)
    if err != nil {
        fmt.Println("Invalid Order ID")
        return
    }

    items, order, err := h.OrderRepo.GetOrderDetails(orderID)
    if err != nil {
        fmt.Println("Failed to fetch order details:", err)
        return
    }

    // Ambil nama customer (user) dari database
    var customerName string
    err = h.DB.QueryRow("SELECT name FROM users WHERE id = $1", order.UserID).Scan(&customerName)
    if err != nil {
        customerName = "Unknown"
    }

    fmt.Printf("\n=== Order #%d Details ===\n", order.ID)
    fmt.Printf("Customer     : %s\n", customerName)
    fmt.Printf("Status       : %s\n", h.getStatusName(order.StatusID))
    fmt.Printf("Order Date   : %s\n", order.CreatedAt.Format("2006-01-02 15:04:05"))
    fmt.Printf("Updated At   : %s\n", order.UpdatedAt.Format("2006-01-02 15:04:05"))
    fmt.Println("-------------------------------------------------")
    fmt.Printf("%-20s | %-3s | %-12s | %-12s\n", "Product Name", "Qty", "Price (Rp)", "Subtotal (Rp)")
    fmt.Println("-------------------------------------------------")

    total := 0.0
    for _, item := range items {
        subtotal := float64(item.Quantity) * item.Price
        fmt.Printf("%-20s | %-3d | %-12.2f | %-12.2f\n",
            item.ProductName, item.Quantity, item.Price, subtotal)
        total += subtotal
    }

    fmt.Println("-------------------------------------------------")
    fmt.Printf("Total Amount: Rp%.2f\n", total)
}

// getStatusName ambil status_name dari status_id
func (h *OrderHandler) getStatusName(statusID int) string {
    var statusName string
    err := h.DB.QueryRow("SELECT status_name FROM status_order WHERE id = $1", statusID).Scan(&statusName)
    if err != nil {
        return "Unknown"
    }
    return statusName
}
