package handler

import (
	"bufio"
	"bytes"
	"errors"
	"go-ecommerce-cli/internal/entity"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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

// --- TEST CASE: SUCCESS ---
func TestRegisterUserCLI_Success(t *testing.T) {
	mockRepo := &mockUserRepo{}
	input := "John Doe\njohn@example.com\npassword123\n"
	reader := bufio.NewReader(strings.NewReader(input))

	var out bytes.Buffer
	handler := &UserHandler{
		UserRepo: mockRepo,
		Reader:   reader,
		Writer:   &out,
	}

	handler.RegisterUserCLI()

	output := out.String()
	assert.Contains(t, output, "Registration successful", "Output should indicate success")
	assert.NotContains(t, output, "error", "Output should not contain error message")
}

// --- TEST CASE: EMPTY INPUT ---
func TestRegisterUserCLI_EmptyInput(t *testing.T) {
	mockRepo := &mockUserRepo{}
	input := "\n\n\n"
	reader := bufio.NewReader(strings.NewReader(input))

	var out bytes.Buffer
	handler := &UserHandler{
		UserRepo: mockRepo,
		Reader:   reader,
		Writer:   &out,
	}

	handler.RegisterUserCLI()

	output := out.String()
	assert.Contains(t, output, "All fields are required.")
}

// --- TEST CASE: EMAIL SUDAH ADA ---
func TestRegisterUserCLI_EmailExists(t *testing.T) {
	mockRepo := &mockUserRepo{}
	input := "Jane Doe\nexists@example.com\npassword123\n"
	reader := bufio.NewReader(strings.NewReader(input))

	var out bytes.Buffer
	handler := &UserHandler{
		UserRepo: mockRepo,
		Reader:   reader,
		Writer:   &out,
	}

	handler.RegisterUserCLI()

	output := out.String()
	assert.Contains(t, output, "email already exists", "Should show error for existing email")
}
