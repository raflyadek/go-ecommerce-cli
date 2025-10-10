package repository

import (
	"database/sql"
	"fmt"
	"go-ecommerce-cli/internal/entity"
)

// ReportRepository interface untuk operasi database laporan
type ReportRepository interface {
	FindCompletedOrders() ([]entity.Order, error)
	FindUserOrders(userID int) ([]entity.Order, error)
	FindUserOrdersByName(userName string) ([]entity.Order, error)
	FindAllUsersWithOrders() ([]entity.UserReport, error)
	GetStockReport() ([]entity.StockReport, error)
	GetBestSellingProducts() ([]entity.BestSellingProduct, error)
	GetReportSummary() (entity.ReportSummary, error)
	GetTotalUsers() (int, error)
	GetTotalOrders() (int, error)
	GetTotalRevenue() (float64, error)
	GetAvgOrderValue() (float64, error)
}

// ReportRepo implementasi dari ReportRepository
type ReportRepo struct {
	DB *sql.DB
}

// NewReportRepo membuat instance baru ReportRepo
func NewReportRepo(db *sql.DB) *ReportRepo {
	return &ReportRepo{DB: db}
}

// FindCompletedOrders mencari semua pesanan dengan status completed
func (r *ReportRepo) FindCompletedOrders() ([]entity.Order, error) {
	query := `
		SELECT 
			o.id,
			u.name AS customer_name,
			o.total_amount,
			s.status_name,
			o.updated_at AS completed_date,
			o.address
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
			&order.Address,
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

// FindUserOrders mencari pesanan berdasarkan user ID
func (r *ReportRepo) FindUserOrders(userID int) ([]entity.Order, error) {
	query := `
		SELECT 
			o.id,
			u.name AS customer_name,
			o.total_amount,
			s.status_name,
			o.updated_at AS completed_date,
			o.address
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
			&order.Address,
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

// FindAllUsersWithOrders mencari semua user yang memiliki pesanan
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

// GetStockReport mendapatkan laporan stok semua produk
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
// GetBestSellingProducts mendapatkan produk terlaris dari view database
func (r *ReportRepo) GetBestSellingProducts() ([]entity.BestSellingProduct, error) {
	query := `SELECT * FROM view_best_selling_products`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying best selling products: %w", err)
	}
	defer rows.Close()

	var products []entity.BestSellingProduct
	for rows.Next() {
		var product entity.BestSellingProduct
		err := rows.Scan(
			&product.ProductID,
			&product.ProductName,
			&product.TotalSold,
			&product.TotalRevenue,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning best selling product: %w", err)
		}
		products = append(products, product)
	}

	return products, nil
}

// FindUserOrdersByName mencari pesanan berdasarkan nama user (case insensitive)
func (r *ReportRepo) FindUserOrdersByName(userName string) ([]entity.Order, error) {
	query := `
		SELECT 
			o.id,
			u.name AS customer_name,
			o.total_amount,
			s.status_name,
			o.updated_at AS completed_date,
			o.address
		FROM orders o
		JOIN users u ON o.user_id = u.id
		JOIN status_order s ON o.status_id = s.id
		WHERE LOWER(u.name) LIKE LOWER($1)
		ORDER BY o.created_at DESC
	`

	rows, err := r.DB.Query(query, "%"+userName+"%")
	if err != nil {
		return nil, fmt.Errorf("error querying user orders by name: %w", err)
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
			&order.Address,
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

// GetReportSummary mendapatkan ringkasan laporan dalam 1 query (backup method)
func (r *ReportRepo) GetReportSummary() (entity.ReportSummary, error) {
	query := `
		SELECT 
			(SELECT COUNT(*) FROM users WHERE role_id = 3) as total_users,
			(SELECT COUNT(*) FROM orders) as total_orders,
			(SELECT COALESCE(SUM(total_amount), 0) FROM orders WHERE status_id IN (2,3,4)) as total_revenue,
			(SELECT COALESCE(AVG(total_amount), 0) FROM orders WHERE status_id IN (2,3,4)) as avg_order_value
	`

	var summary entity.ReportSummary
	err := r.DB.QueryRow(query).Scan(
		&summary.TotalUsers,
		&summary.TotalOrders,
		&summary.TotalRevenue,
		&summary.AvgOrderValue,
	)
	if err != nil {
		return summary, fmt.Errorf("error querying report summary: %w", err)
	}

	return summary, nil
}

// GetTotalUsers menghitung total user untuk goroutine
func (r *ReportRepo) GetTotalUsers() (int, error) {
	var count int
	err := r.DB.QueryRow("SELECT COUNT(*) FROM users WHERE role_id = 3").Scan(&count)
	return count, err
}

// GetTotalOrders menghitung total pesanan untuk goroutine
func (r *ReportRepo) GetTotalOrders() (int, error) {
	var count int
	err := r.DB.QueryRow("SELECT COUNT(*) FROM orders").Scan(&count)
	return count, err
}

// GetTotalRevenue menghitung total pendapatan untuk goroutine
func (r *ReportRepo) GetTotalRevenue() (float64, error) {
	var revenue float64
	err := r.DB.QueryRow("SELECT COALESCE(SUM(total_amount), 0) FROM orders WHERE status_id IN (2,3,4)").Scan(&revenue)
	return revenue, err
}

// GetAvgOrderValue menghitung rata-rata nilai pesanan untuk goroutine
func (r *ReportRepo) GetAvgOrderValue() (float64, error) {
	var avg float64
	err := r.DB.QueryRow("SELECT COALESCE(AVG(total_amount), 0) FROM orders WHERE status_id IN (2,3,4)").Scan(&avg)
	return avg, err
}