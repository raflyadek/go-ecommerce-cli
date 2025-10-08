package repository

import (
	"database/sql"
	"fmt"
	"go-ecommerce-cli/internal/entity"
)

type ProductRepository struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}

func (r *ProductRepository) GetAll() ([]entity.Product, error) {
	rows, err := r.DB.Query(`SELECT id, name, description, price FROM products ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []entity.Product{}
	for rows.Next() {
		var p entity.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price)
		if err != nil {
			fmt.Println("Scan error:", err)
			continue
		}
		products = append(products, p)
	}

	return products, nil
}
