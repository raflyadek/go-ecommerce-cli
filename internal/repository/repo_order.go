package repository

import (
	"database/sql"
	"fmt"
	"go-ecommerce-cli/internal/entity"
)

// OrderRepository: Kontrak untuk akses data Order
type OrderRepository interface {
	FindCompletedOrders() ([]entity.Order, error)
}

// ProductRepository: menunggu Interface di product repo di buat.

type OrderRepo struct {
	DB *sql.DB
}

func NewOrderRepo(db *sql.DB) *OrderRepo {
	return &OrderRepo{DB: db}
}

func (r *OrderRepo) FindCompletedOrders() ([]entity.Order, error) {
	query := `
		SELECT 
			o.id, 
			u.name AS customer_name, 
			o.total_amount, 
			o.status, 
			o.completed_at
		FROM orders o
		JOIN users u ON o.user_id = u.id -- Asumsi tabel users ada
		WHERE o.status = 'Completed'
		ORDER BY o.completed_at DESC
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying completed orders: %w", err)
	}
	defer rows.Close()

	var orders []entity.Order
	for rows.Next() {
		var order entity.Order
		var completedAt sql.NullTime

		err := rows.Scan(
			&order.ID,
			&order.CustomerName,
			&order.TotalAmount,
			&order.Status,
			&completedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning order row: %w", err)
		}

		if completedAt.Valid {
			order.CompletedDate = completedAt.Time
		}

		orders = append(orders, order)
	}

	return orders, rows.Err()
}
