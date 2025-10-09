package repository

import (
	"database/sql"
)

type OrderRepository interface {
	// Kosongkan untuk teman yang mengerjakan order
	// Semua fungsi report sudah dipindah ke repo_report.go
}

type OrderRepo struct {
	DB *sql.DB
}

func NewOrderRepo(db *sql.DB) *OrderRepo {
	return &OrderRepo{DB: db}
}

// Kosongkan untuk teman yang mengerjakan order
// Semua fungsi report sudah dipindah ke repo_report.go