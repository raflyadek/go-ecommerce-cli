package repository

import (
	"database/sql"
	"fmt"
	"go-ecommerce-cli/internal/entity"
	"time"
)

type OrderRepository interface {
	CreateOrder(order entity.Order) (int, error)
	GetUserOrders(userID int) ([]entity.OrderSummary, error)
	GetOrderDetails(orderID int) ([]entity.OrderItemDetails, entity.Order, error)
	UpdateOrderStatus(orderID int, statusName string) error
}

type OrderRepo struct {
	DB *sql.DB
}

func NewOrderRepo(db *sql.DB) *OrderRepo {
	return &OrderRepo{DB: db}
}

func (r *OrderRepo) CreateOrder(order entity.Order) (int, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var orderID int
	query := `INSERT INTO orders (user_id, total_amount, status_id, address, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id`
	err = tx.QueryRow(query, order.UserID, order.TotalAmount, 1, order.Address).Scan(&orderID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert order: %w", err)
	}

	// insert order items
	for _, item := range order.Items {
		queryItem := `INSERT INTO order_items (order_id, product_id, quantity, price_at_order)
		              VALUES ($1, $2, $3, $4)`
		_, err = tx.Exec(queryItem, orderID, item.ProductID, item.Quantity, item.Price)
		if err != nil {
			return 0, fmt.Errorf("failed to insert order item: %w", err)
		}
	}
	return orderID, nil
}

func (r *OrderRepo) GetUserOrders(userID int) ([]entity.OrderSummary, error) {
    query := `
    SELECT o.id, s.status_name, o.total_amount, COUNT(oi.product_id) AS item_count, o.created_at
    FROM orders o
    JOIN status_order s ON o.status_id = s.id
    LEFT JOIN order_items oi ON o.id = oi.order_id
    WHERE o.user_id = $1
    GROUP BY o.id, s.status_name, o.total_amount, o.created_at
    ORDER BY o.created_at DESC
    `
    rows, err := r.DB.Query(query, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch user orders: %w", err)
    }
    defer rows.Close()

    var orders []entity.OrderSummary
    for rows.Next() {
        var o entity.OrderSummary
        var createdAt time.Time
        err := rows.Scan(&o.ID, &o.Status, &o.Total, &o.ItemCount, &createdAt)
        if err != nil {
            return nil, fmt.Errorf("failed to scan order: %w", err)
        }
        o.OrderDate = createdAt.Format("2006-01-02 15:04:05")
        orders = append(orders, o)
    }
    return orders, nil
}

func (r *OrderRepo) GetOrderDetails(orderID int) ([]entity.OrderItemDetails, entity.Order, error) {
    query := `
    SELECT p.name, oi.quantity, oi.price_at_order
    FROM order_items oi
    JOIN products p ON oi.product_id = p.id
    WHERE oi.order_id = $1
    `
    rows, err := r.DB.Query(query, orderID)
    if err != nil {
        return nil, entity.Order{}, fmt.Errorf("failed to fetch order items: %w", err)
    }
    defer rows.Close()

    var items []entity.OrderItemDetails
    for rows.Next() {
        var i entity.OrderItemDetails
        err := rows.Scan(&i.ProductName, &i.Quantity, &i.Price)
        if err != nil {
            return nil, entity.Order{}, fmt.Errorf("failed to scan order item: %w", err)
        }
        items = append(items, i)
    }

    var order entity.Order
    queryOrder := `SELECT id, total_amount, status_id, created_at FROM orders WHERE id = $1`
    err = r.DB.QueryRow(queryOrder, orderID).Scan(&order.ID, &order.TotalAmount, &order.StatusID, &order.CreatedAt)
    if err != nil {
        return nil, entity.Order{}, fmt.Errorf("failed to fetch order info: %w", err)
    }

    return items, order, nil
}

// UpdateOrderStatus updates the status of an order (for staff / admin)
func (r *OrderRepo) UpdateOrderStatus(orderID int, statusName string) error {
    var statusID int
    // Ambil status_id dari tabel status_order
    err := r.DB.QueryRow("SELECT id FROM status_order WHERE LOWER(status_name) = LOWER($1) LIMIT 1", statusName).Scan(&statusID)
    if err != nil {
        return fmt.Errorf("status '%s' not found", statusName)
    }

    // Update orders
    _, err = r.DB.Exec("UPDATE orders SET status_id = $1, updated_at = NOW() WHERE id = $2", statusID, orderID)
    if err != nil {
        return fmt.Errorf("failed to update order status: %v", err)
    }

    return nil
}
