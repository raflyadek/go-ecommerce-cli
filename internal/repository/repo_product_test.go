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

	repo := NewProductRepository(db)

	//Mock add product
	mock.ExpectExec("INSERT INTO products").
		WithArgs("Baju", "Baju panjang", float64(10000), 10).
		WillReturnResult(sqlmock.NewResult(1, 1))
	
	err = repo.AddProduct("Baju", "Baju panjang", 10000, 10)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepoProduct_DeleteProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewProductRepository(db)
	mock.ExpectExec("DELETE FROM products WHERE id =").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	
	err = repo.DeleteProduct(1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestRepoProduct_UpdateProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewProductRepository(db)

	mock.ExpectExec("UPDATE products SET name =").
		WithArgs("Baju", "Baju panjang", float64(12000), 15, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.UpdateProduct("Baju", "Baju panjang", 12000, 15, 1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestRepoProduct_ShowProductById(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewProductRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "stock"}).
		AddRow(1, "Baju", "Baju panjang", float64(10000), 10)

	mock.ExpectQuery("SELECT id, name, description, price, stock FROM products WHERE id =").
		WithArgs(1).
		WillReturnRows(rows)

	product, err := repo.ShowProductById(1)
	assert.NoError(t, err)
	assert.Equal(t, 1, product.ID)
	assert.Equal(t, "Baju", product.Name)
	assert.Equal(t, "Baju panjang", product.Description)
	assert.Equal(t, float64(10000), product.Price)
	assert.Equal(t, 10, product.Stock)

}