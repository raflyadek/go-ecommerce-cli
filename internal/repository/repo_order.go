package repository

import (
	"database/sql"
	"fmt"
	"go-ecommerce-cli/internal/entity"
)

type OrderRepository interface {
	FindCompletedOrders() ([]entity.Order, error)
}

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
			s.status_name,
			o.updated_at AS completed_date
		FROM orders o
		JOIN users u ON o.user_id = u.id
		JOIN status_order s ON o.status_id = s.id
		WHERE s.status_name = 'completed'
		ORDER BY o.updated_at DESC
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
			&order.StatusName,
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

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return orders, nil
}
