package repository

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestRepoProduct_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	repo := NewProductRepository(db)

	//Mock get all products
	rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "stock"}).
		AddRow(1, "Baju", "Baju panjang", 10000.00, 10)
	mock.ExpectQuery("SELECT id, name, description, price, stock FROM products ORDER BY id").WillReturnRows(rows)

	products, err := repo.GetAll()

	//hasil
	assert.NoError(t, err)
	assert.Len(t, products, 1)
	assert.Equal(t, "Baju", products[0].Name)
	assert.Equal(t, "Baju panjang", products[0].Description)
	assert.Equal(t, float64(10000), products[0].Price)
	assert.Equal(t, 10, products[0].Stock)

}

func TestRepoProduct_AddProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	//Mock add product
	mock.ExpectExec("INSERT INTO products (name, description, price, stock").
		WithArgs("Baju", "Baju panjang", float64(10000), 10).
		WillReturnResult(sqlmock.NewResult(1, 1))
}

func TestRepoProduct_DeleteProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("")
}
func TestRepoProduct_UpdateProduct(t *testing.T) {

}
func TestRepoProduct_ShowProductById(t *testing.T) {

}