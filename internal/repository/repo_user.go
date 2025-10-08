package repository

import (
	"database/sql"
	"errors"
	"go-ecommerce-cli/internal/entity"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) GetByEmail(email string) (entity.User, error) {
	var user entity.User
	query := `
	SELECT u.id, u.name, u.email, u.password, r.name AS role_name
	FROM users u
	JOIN roles r ON u.role_id = r.id
	WHERE u.email = $1
	`
	err := r.DB.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.RoleName)
	if err != nil {
		return user, errors.New("user not found")
	}
	return user, nil
}

func (r *UserRepository) RegisterUser(name, email, hashed string) error {
	query := `INSERT INTO users (name, email, password, role_id) VALUES ($1, $2, $3, 3)`
	_, err := r.DB.Exec(query, name, email, hashed)
	return err
}

func (r *UserRepository) CreateStaff(name, email, hashed string) error {
	query := `INSERT INTO users (name, email, password, role_id) VALUES ($1, $2, $3, 2)`
	_, err := r.DB.Exec(query, name, email, hashed)
	return err
}
