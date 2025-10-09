package repository

import (
	"database/sql"
	"fmt"
	"go-ecommerce-cli/internal/entity"
	"log"
)

type ProductRepository struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}

func (r *ProductRepository) GetAll() ([]entity.Product, error) {
	rows, err := r.DB.Query(`SELECT id, name, description, price, stock FROM products ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []entity.Product{}
	for rows.Next() {
		var p entity.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock)
		if err != nil {
			fmt.Println("Scan error:", err)
			continue
		}
		products = append(products, p)
	}

	return products, nil
}

func (r *ProductRepository) AddProduct(name, description string, price float64, stock int) error {
	_, err := r.DB.Exec(`
		INSERT INTO products (name, description, price, stock)
		VALUES ($1, $2, $3, $4)
	`, name, description, price, stock)

	if err != nil {
		log.Println("Failed to add product", err)
		return err
	}
	return err
}

func (r *ProductRepository) DeleteProduct(id int) error {
	_, err := r.DB.Exec(`DELETE FROM products WHERE id = $1`, id)

	if err != nil {
		log.Println("Failed to delete product", err)
		return err
	}
	return err
}

func (r *ProductRepository) UpdateProduct(name, description string, price float64, stock, id int) error {
	_, err := r.DB.Exec(`
		UPDATE products SET name = $1, description = $2, price = $3, stock = $4
		WHERE id = $5
	`, name, description, price, stock, id)

	if err != nil {
		log.Println("Failed to update product", err)
		return err
	}
	return err
}

func (r *ProductRepository) ShowProductById(id int) (entity.Product, error) {
	var product entity.Product
	rows, err := r.DB.Query(`
		SELECT p.id, p.name, p.description, p.price, p.stock FROM products WHERE id = $1
	`, id)

	if err != nil {
		log.Println("Error fetching data from table products:", err)
		return product, err
	}

	// var products []entity.Product

	for rows.Next() {

		if err := rows.Scan(
			&product.ID,
			&product.Name, 
			&product.Description, 
			&product.Price, 
			&product.Stock,
			); err != nil {
				log.Println("error scanning product:", err)
			}
	}

	return product, err
}