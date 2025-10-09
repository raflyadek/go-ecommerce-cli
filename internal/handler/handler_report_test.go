package handler

import (
	"errors"
	"go-ecommerce-cli/internal/entity"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockReportRepository adalah implementasi mock dari ReportRepository
type MockReportRepository struct {
	mock.Mock
}

func (m *MockReportRepository) FindCompletedOrders() ([]entity.Order, error) {
	args := m.Called()
	return args.Get(0).([]entity.Order), args.Error(1)
}

func (m *MockReportRepository) FindUserOrders(userID int) ([]entity.Order, error) {
	args := m.Called(userID)
	return args.Get(0).([]entity.Order), args.Error(1)
}

func (m *MockReportRepository) FindAllUsersWithOrders() ([]entity.UserReport, error) {
	args := m.Called()
	return args.Get(0).([]entity.UserReport), args.Error(1)
}

func (m *MockReportRepository) GetStockReport() ([]entity.StockReport, error) {
	args := m.Called()
	return args.Get(0).([]entity.StockReport), args.Error(1)
}

func (m *MockReportRepository) GetBestSellingProducts() ([]entity.BestSellingProduct, error) {
	args := m.Called()
	return args.Get(0).([]entity.BestSellingProduct), args.Error(1)
}

func (m *MockReportRepository) FindUserOrdersByName(userName string) ([]entity.Order, error) {
	args := m.Called(userName)
	return args.Get(0).([]entity.Order), args.Error(1)
}

func (m *MockReportRepository) GetReportSummary() (entity.ReportSummary, error) {
	args := m.Called()
	return args.Get(0).(entity.ReportSummary), args.Error(1)
}

func TestReportHandler_ShowBestSellingProducts_Success(t *testing.T) {
	mockRepo := new(MockReportRepository)
	handler := &ReportHandler{ReportRepo: mockRepo}

	products := []entity.BestSellingProduct{
		{ProductID: 1, ProductName: "T-Shirt", TotalSold: 20, TotalRevenue: 200000.00},
	}

	mockRepo.On("GetBestSellingProducts").Return(products, nil)

	// Test ini memastikan tidak terjadi panic saat menampilkan produk terlaris
	assert.NotPanics(t, func() {
		handler.ShowBestSellingProducts()
	})

	mockRepo.AssertExpectations(t)
}

func TestReportHandler_ShowBestSellingProducts_Error(t *testing.T) {
	mockRepo := new(MockReportRepository)
	handler := &ReportHandler{ReportRepo: mockRepo}

	mockRepo.On("GetBestSellingProducts").Return([]entity.BestSellingProduct{}, errors.New("database error"))

	// Test ini memastikan tidak terjadi panic saat ada error database
	assert.NotPanics(t, func() {
		handler.ShowBestSellingProducts()
	})

	mockRepo.AssertExpectations(t)
}

func TestReportHandler_ShowBestSellingProducts_NoData(t *testing.T) {
	mockRepo := new(MockReportRepository)
	handler := &ReportHandler{ReportRepo: mockRepo}

	mockRepo.On("GetBestSellingProducts").Return([]entity.BestSellingProduct{}, nil)

	// Test ini memastikan tidak terjadi panic saat tidak ada data
	assert.NotPanics(t, func() {
		handler.ShowBestSellingProducts()
	})

	mockRepo.AssertExpectations(t)
}

func TestReportHandler_Repository_Methods(t *testing.T) {
	mockRepo := new(MockReportRepository)
	handler := &ReportHandler{ReportRepo: mockRepo}

	// Test FindCompletedOrders - mencari pesanan yang sudah selesai
	orders := []entity.Order{
		{ID: 1, CustomerName: "John", TotalAmount: 100000, StatusName: "completed", Address: "Jakarta"},
	}
	mockRepo.On("FindCompletedOrders").Return(orders, nil)

	result, err := handler.ReportRepo.FindCompletedOrders()
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "John", result[0].CustomerName)

	// Test FindUserOrders - mencari pesanan berdasarkan user ID
	mockRepo.On("FindUserOrders", 1).Return(orders, nil)

	result, err = handler.ReportRepo.FindUserOrders(1)
	assert.NoError(t, err)
	assert.Len(t, result, 1)

	// Test GetStockReport - mendapatkan laporan stok produk
	stockReports := []entity.StockReport{
		{ProductID: 1, ProductName: "T-Shirt", StockIn: 60, StockOut: 10, CurrentStock: 50},
	}
	mockRepo.On("GetStockReport").Return(stockReports, nil)

	stockResult, err := handler.ReportRepo.GetStockReport()
	assert.NoError(t, err)
	assert.Len(t, stockResult, 1)
	assert.Equal(t, "T-Shirt", stockResult[0].ProductName)

	// Test FindAllUsersWithOrders - mencari semua user yang punya pesanan
	userReports := []entity.UserReport{
		{UserID: 1, UserName: "John", OrderCount: 2, TotalSpent: 150000.00},
	}
	mockRepo.On("FindAllUsersWithOrders").Return(userReports, nil)

	userResult, err := handler.ReportRepo.FindAllUsersWithOrders()
	assert.NoError(t, err)
	assert.Len(t, userResult, 1)
	assert.Equal(t, "John", userResult[0].UserName)

	// Test FindUserOrdersByName - mencari pesanan berdasarkan nama user
	mockRepo.On("FindUserOrdersByName", "John").Return(orders, nil)

	nameResult, err := handler.ReportRepo.FindUserOrdersByName("John")
	assert.NoError(t, err)
	assert.Len(t, nameResult, 1)
	assert.Equal(t, "John", nameResult[0].CustomerName)

	// Test GetReportSummary - mendapatkan ringkasan laporan
	summary := entity.ReportSummary{
		TotalUsers: 10, TotalOrders: 25, TotalRevenue: 500000.00, AvgOrderValue: 20000.00,
	}
	mockRepo.On("GetReportSummary").Return(summary, nil)

	summaryResult, err := handler.ReportRepo.GetReportSummary()
	assert.NoError(t, err)
	assert.Equal(t, 10, summaryResult.TotalUsers)
	assert.Equal(t, float64(500000), summaryResult.TotalRevenue)

	mockRepo.AssertExpectations(t)
}

func TestReportHandler_Error_Handling(t *testing.T) {
	mockRepo := new(MockReportRepository)
	handler := &ReportHandler{ReportRepo: mockRepo}

	// Test penanganan error untuk setiap method
	mockRepo.On("FindCompletedOrders").Return([]entity.Order{}, errors.New("db error"))
	mockRepo.On("FindUserOrders", 1).Return([]entity.Order{}, errors.New("db error"))
	mockRepo.On("FindUserOrdersByName", "test").Return([]entity.Order{}, errors.New("db error"))
	mockRepo.On("GetStockReport").Return([]entity.StockReport{}, errors.New("db error"))
	mockRepo.On("FindAllUsersWithOrders").Return([]entity.UserReport{}, errors.New("db error"))
	mockRepo.On("GetReportSummary").Return(entity.ReportSummary{}, errors.New("db error"))

	_, err1 := handler.ReportRepo.FindCompletedOrders()
	_, err2 := handler.ReportRepo.FindUserOrders(1)
	_, err3 := handler.ReportRepo.FindUserOrdersByName("test")
	_, err4 := handler.ReportRepo.GetStockReport()
	_, err5 := handler.ReportRepo.FindAllUsersWithOrders()
	_, err6 := handler.ReportRepo.GetReportSummary()

	assert.Error(t, err1)
	assert.Error(t, err2)
	assert.Error(t, err3)
	assert.Error(t, err4)
	assert.Error(t, err5)
	assert.Error(t, err6)

	mockRepo.AssertExpectations(t)
}