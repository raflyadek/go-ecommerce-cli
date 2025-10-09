package repository

import (
	"database/sql"
	"fmt"
	"go-ecommerce-cli/internal/entity"
)

type ReportRepository interface {
	FindCompletedOrders() ([]entity.Order, error)
	FindUserOrders(userID int) ([]entity.Order, error)
	FindAllUsersWithOrders() ([]entity.UserReport, error)
	GetStockReport() ([]entity.StockReport, error)
}

type ReportRepo struct {
	DB *sql.DB
}

func NewReportRepo(db *sql.DB) *ReportRepo {
	return &ReportRepo{DB: db}
}

func (r *ReportRepo) FindCompletedOrders() ([]entity.Order, error) {
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

func (r *ReportRepo) FindUserOrders(userID int) ([]entity.Order, error) {
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
		WHERE o.user_id = $1
		ORDER BY o.created_at DESC
	`

	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying user orders: %w", err)
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

	return orders, nil
}

func (r *ReportRepo) FindAllUsersWithOrders() ([]entity.UserReport, error) {
	query := `
		SELECT 
			u.id,
			u.name,
			COUNT(o.id) as order_count,
			COALESCE(SUM(o.total_amount), 0) as total_spent
		FROM users u
		LEFT JOIN orders o ON u.id = o.user_id
		WHERE u.role_id = 3
		GROUP BY u.id, u.name
		HAVING COUNT(o.id) > 0
		ORDER BY total_spent DESC
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying user reports: %w", err)
	}
	defer rows.Close()

	var reports []entity.UserReport
	for rows.Next() {
		var report entity.UserReport
		err := rows.Scan(
			&report.UserID,
			&report.UserName,
			&report.OrderCount,
			&report.TotalSpent,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning user report row: %w", err)
		}
		reports = append(reports, report)
	}

	return reports, nil
}

func (r *ReportRepo) GetStockReport() ([]entity.StockReport, error) {
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

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying stock report: %w", err)
	}
	defer rows.Close()

	var reports []entity.StockReport
	for rows.Next() {
		var report entity.StockReport
		var totalSold int
		err := rows.Scan(
			&report.ProductID,
			&report.ProductName,
			&report.CurrentStock,
			&totalSold,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning stock report row: %w", err)
		}
		
		// Calculate Stock In = Current Stock + Total Sold
		report.StockIn = report.CurrentStock + totalSold
		report.StockOut = totalSold
		
		reports = append(reports, report)
	}

	return reports, nil
}