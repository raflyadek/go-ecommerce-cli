package repository

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestReportRepo_FindCompletedOrders(t *testing.T) {
	// Setup database mock dan repository
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewReportRepo(db)

	// Mock data pesanan yang sudah selesai
	rows := sqlmock.NewRows([]string{"id", "customer_name", "total_amount", "status_name", "completed_date", "address"}).
		AddRow(1, "John Doe", 100000.00, "completed", time.Now(), "Jakarta")

	mock.ExpectQuery("SELECT.*FROM orders o.*WHERE s.status_name = 'completed'").WillReturnRows(rows)

	// Test fungsi FindCompletedOrders
	orders, err := repo.FindCompletedOrders()

	// Verifikasi hasil
	assert.NoError(t, err)
	assert.Len(t, orders, 1)
	assert.Equal(t, "John Doe", orders[0].CustomerName)
	assert.Equal(t, float64(100000), orders[0].TotalAmount)
}

func TestReportRepo_FindUserOrders(t *testing.T) {
	// Setup database mock dan repository
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewReportRepo(db)

	// Mock data pesanan user tertentu
	rows := sqlmock.NewRows([]string{"id", "customer_name", "total_amount", "status_name", "completed_date", "address"}).
		AddRow(1, "John Doe", 50000.00, "pending", time.Now(), "Jakarta")

	mock.ExpectQuery("SELECT.*FROM orders o.*WHERE o.user_id = \\$1").WithArgs(1).WillReturnRows(rows)

	// Test fungsi FindUserOrders dengan user ID = 1
	orders, err := repo.FindUserOrders(1)

	// Verifikasi hasil
	assert.NoError(t, err)
	assert.Len(t, orders, 1)
	assert.Equal(t, "John Doe", orders[0].CustomerName)
}

func TestReportRepo_FindAllUsersWithOrders(t *testing.T) {
	// Setup database mock dan repository
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewReportRepo(db)

	// Mock data laporan user dengan pesanan
	rows := sqlmock.NewRows([]string{"id", "name", "order_count", "total_spent"}).
		AddRow(1, "John Doe", 2, 150000.00)

	mock.ExpectQuery("SELECT.*FROM users u.*WHERE u.role_id = 3").WillReturnRows(rows)

	// Test fungsi FindAllUsersWithOrders
	reports, err := repo.FindAllUsersWithOrders()

	// Verifikasi hasil
	assert.NoError(t, err)
	assert.Len(t, reports, 1)
	assert.Equal(t, "John Doe", reports[0].UserName)
	assert.Equal(t, 2, reports[0].OrderCount)
}

func TestReportRepo_GetStockReport(t *testing.T) {
	// Setup database mock dan repository
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewReportRepo(db)

	// Mock data laporan stok produk
	rows := sqlmock.NewRows([]string{"id", "name", "stock", "total_sold"}).
		AddRow(1, "T-Shirt", 50, 10)

	mock.ExpectQuery("SELECT.*FROM products p").WillReturnRows(rows)

	// Test fungsi GetStockReport
	reports, err := repo.GetStockReport()

	// Verifikasi hasil dan kalkulasi stok
	assert.NoError(t, err)
	assert.Len(t, reports, 1)
	assert.Equal(t, "T-Shirt", reports[0].ProductName)
	assert.Equal(t, 60, reports[0].StockIn) // 50 + 10
	assert.Equal(t, 10, reports[0].StockOut)
}

func TestReportRepo_GetBestSellingProducts(t *testing.T) {
	// Setup database mock dan repository
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewReportRepo(db)

	// Mock data produk terlaris dari view database
	rows := sqlmock.NewRows([]string{"product_id", "product_name", "total_sold", "total_revenue"}).
		AddRow(1, "T-Shirt", 20, 200000.00)

	mock.ExpectQuery("SELECT \\* FROM view_best_selling_products").WillReturnRows(rows)

	// Test fungsi GetBestSellingProducts
	products, err := repo.GetBestSellingProducts()

	// Verifikasi hasil
	assert.NoError(t, err)
	assert.Len(t, products, 1)
	assert.Equal(t, "T-Shirt", products[0].ProductName)
	assert.Equal(t, 20, products[0].TotalSold)
}

func TestReportRepo_FindCompletedOrders_Error(t *testing.T) {
	// Setup database mock dan repository
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewReportRepo(db)

	// Mock error koneksi database
	mock.ExpectQuery("SELECT.*FROM orders o.*WHERE s.status_name = 'completed'").WillReturnError(sql.ErrConnDone)

	// Test fungsi FindCompletedOrders dengan error
	orders, err := repo.FindCompletedOrders()

	// Verifikasi error handling
	assert.Error(t, err)
	assert.Nil(t, orders)
}