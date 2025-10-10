package handler

import (
	"bufio"
	"errors"
	"go-ecommerce-cli/internal/entity"
	"strings"
	"testing"
)

type mockUserRepo struct{}

func (m *mockUserRepo) GetByEmail(email string) (entity.User, error) {
    if email == "user@example.com" {
        return entity.User{
            ID:       1,
            Name:     "User",
            Email:    "user@example.com",
            Password: "$2a$10$somehashedpassword",
            RoleName: "customer",
        }, nil
    }
    return entity.User{}, errors.New("user not found")
}

func (m *mockUserRepo) RegisterUser(name, email, hashed string) error {
    if email == "exists@example.com" {
        return errors.New("email already exists")
    }
    return nil
}

func (m *mockUserRepo) CreateStaff(name, email, hashed string) error {
    if email == "" {
        return errors.New("email required")
    }
    return nil
}

func TestRegisterUserCLI_Success(t *testing.T) {
    mockRepo := &mockUserRepo{}
    input := "John Doe\njohn@example.com\npassword123\n"
    reader := bufio.NewReader(strings.NewReader(input))

    handler := &UserHandler{
        UserRepo: mockRepo,
        Reader:   reader,
    }

    handler.RegisterUserCLI()
}

func TestRegisterUserCLI_EmptyInput(t *testing.T) {
    mockRepo := &mockUserRepo{}
    input := "\n\n\n"
    reader := bufio.NewReader(strings.NewReader(input))

    handler := &UserHandler{
        UserRepo: mockRepo,
        Reader:   reader,
    }

    handler.RegisterUserCLI()
}
