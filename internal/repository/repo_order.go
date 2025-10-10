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
	GetAllOrders() ([]entity.OrderSummary, error)
	GetHistoryOrders(userID int) ([]entity.OrderSummary, error)
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

    //query with trigger to set status_id = 1 when insert
	query := `
        INSERT INTO orders (user_id, total_amount, address, created_at, updated_at)
        VALUES ($1, $2, $3, NOW(), NOW())
        RETURNING id
	`
	err = tx.QueryRow(query, order.UserID, order.TotalAmount, order.Address).Scan(&orderID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert order: %w", err)
	}

	// Insert or update order items
	for _, item := range order.Items {
		queryItem := `
			INSERT INTO order_items (order_id, product_id, quantity, price_at_order)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (order_id, product_id)
			DO UPDATE SET quantity = order_items.quantity + EXCLUDED.quantity
		`
		_, err = tx.Exec(queryItem, orderID, item.ProductID, item.Quantity, item.Price)
		if err != nil {
			return 0, fmt.Errorf("failed to insert or update order item: %w", err)
		}
	}

	return orderID, nil
}

func (r *OrderRepo) fetchOrders(query string, userID int) ([]entity.OrderSummary, error) {
    rows, err := r.DB.Query(query, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch orders: %w", err)
    }
    defer rows.Close()

    var orders []entity.OrderSummary
    for rows.Next() {
        var o entity.OrderSummary
        var createdAt time.Time
        if err := rows.Scan(&o.ID, &o.Status, &o.Total, &o.ItemCount, &createdAt); err != nil {
            return nil, fmt.Errorf("failed to scan order: %w", err)
        }
        o.OrderDate = createdAt.Format("2006-01-02 15:04:05")
        orders = append(orders, o)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("rows error: %w", err)
    }

    return orders, nil
}

func (r *OrderRepo) fetchAllOrders(query string) ([]entity.OrderSummary, error) {
    rows, err := r.DB.Query(query)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch orders: %w", err)
    }
    defer rows.Close()
    return r.scanOrders(rows)
}

func (r *OrderRepo) scanOrders(rows *sql.Rows) ([]entity.OrderSummary, error) {
    var orders []entity.OrderSummary
    for rows.Next() {
        var o entity.OrderSummary
        var createdAt time.Time
        if err := rows.Scan(&o.ID, &o.Status, &o.Total, &o.ItemCount, &createdAt); err != nil {
            return nil, fmt.Errorf("failed to scan order: %w", err)
        }
        o.OrderDate = createdAt.Format("2006-01-02 15:04:05")
        orders = append(orders, o)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("rows error: %w", err)
    }
    return orders, nil
}

func (r *OrderRepo) GetAllOrders() ([]entity.OrderSummary, error) {
    query := `
        SELECT 
            o.id,
            s.status_name,
            o.total_amount,
            COUNT(oi.product_id) AS item_count,
            o.created_at
        FROM orders o
        JOIN status_order s ON o.status_id = s.id
        LEFT JOIN order_items oi ON o.id = oi.order_id
        GROUP BY o.id, s.status_name, o.total_amount, o.created_at
        ORDER BY o.created_at DESC
    `
    return r.fetchAllOrders(query)
}

func (r *OrderRepo) GetUserOrders(userID int) ([]entity.OrderSummary, error) {
    query := `
        SELECT 
            o.id, 
            o.status_name, 
            o.total_amount, 
            o.item_count, 
            o.created_at
        FROM orders_active o
        WHERE o.user_id = $1
        ORDER BY o.created_at DESC
    `
    return r.fetchOrders(query, userID)
}

func (r *OrderRepo) GetHistoryOrders(userID int) ([]entity.OrderSummary, error) {
    query := `
        SELECT 
            o.id, 
            o.status_name, 
            o.total_amount, 
            o.item_count, 
            o.created_at
        FROM orders_history o
        WHERE o.user_id = $1
        ORDER BY o.created_at DESC
    `
    return r.fetchOrders(query, userID)
}

func (r *OrderRepo) GetOrderDetails(orderID int) ([]entity.OrderItemDetails, entity.Order, error) {
    // Ambil items
    queryItems := `
    SELECT p.name, oi.quantity, oi.price_at_order
    FROM order_items oi
    JOIN products p ON oi.product_id = p.id
    WHERE oi.order_id = $1
    `
    rows, err := r.DB.Query(queryItems, orderID)
    if err != nil {
        return nil, entity.Order{}, fmt.Errorf("failed to fetch order items: %w", err)
    }
    defer rows.Close()

    var items []entity.OrderItemDetails
    for rows.Next() {
        var i entity.OrderItemDetails
        if err := rows.Scan(&i.ProductName, &i.Quantity, &i.Price); err != nil {
            return nil, entity.Order{}, fmt.Errorf("failed to scan order item: %w", err)
        }
        items = append(items, i)
    }

    // Ambil info order + user_id
    var order entity.Order
    queryOrder := `
    SELECT id, user_id, total_amount, status_id, created_at, updated_at
    FROM orders
    WHERE id = $1
    `
    err = r.DB.QueryRow(queryOrder, orderID).Scan(&order.ID, &order.UserID, &order.TotalAmount, &order.StatusID, &order.CreatedAt, &order.UpdatedAt)
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
