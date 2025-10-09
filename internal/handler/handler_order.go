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
	DB        *sql.DB
}

func NewOrderHandler(orderRepo repository.OrderRepository, db *sql.DB) *OrderHandler {
	return &OrderHandler{
		OrderRepo: orderRepo,
		DB:        db,
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

		item := entity.OrderItem{
			ProductID: product.ID,
			Quantity:  qty,
			Price:     product.Price * float64(qty),
		}
		items = append(items, item)
		totalAmount += item.Price

		fmt.Printf("%d %s successfully added with price is %.2f.\n", qty, product.Name, item.Price)

		fmt.Print("Do you want to make another order? (y/n): ")
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
	}

	orderID, err := h.OrderRepo.CreateOrder(order)
	if err != nil {
		fmt.Println("Failed to create order:", err)
		return
	}

	fmt.Printf("Your order ID %d has been created with total price is %.2f\n", orderID, totalAmount)
	fmt.Printf("Delivery Address: %s\n", address)
}

// findProductByIDOrName searches a product by ID or Name in the database
func (h *OrderHandler) findProductByIDOrName(input string) (entity.Product, error) {
	var product entity.Product

	id, err := strconv.Atoi(input)
	if err == nil {
		query := `SELECT id, name, price FROM products WHERE id = $1 LIMIT 1`
		row := h.DB.QueryRow(query, id)
		err := row.Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			return product, fmt.Errorf("product not found with ID %d", id)
		}
		return product, nil
	}

	query := `SELECT id, name, price FROM products WHERE name ILIKE $1 LIMIT 1`
	row := h.DB.QueryRow(query, "%"+input+"%")
	err = row.Scan(&product.ID, &product.Name, &product.Price)
	if err != nil {
		return product, fmt.Errorf("product not found with name '%s'", input)
	}

	return product, nil
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

    fmt.Print("View order details with ID (\"CTRL + C\" to return): ")
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
    fmt.Println("\n=== Update Order Status === (\"Ctrl + C\" to return to dashboard)")

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
