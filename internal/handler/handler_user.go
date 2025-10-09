package handler

import (
	"bufio"
	"database/sql"
	"fmt"
	"go-ecommerce-cli/internal/entity"
	"go-ecommerce-cli/internal/repository"
	"go-ecommerce-cli/pkg/utils"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
)

type UserHandler struct {
	UserRepo *repository.UserRepository
	Reader   *bufio.Reader
	DB       *sql.DB
}

func NewUserHandler(repo *repository.UserRepository, db *sql.DB) *UserHandler {
	return &UserHandler{
		UserRepo: repo,
		Reader:   bufio.NewReader(os.Stdin),
		DB:       db,
	}
}

// --- LOGIN ---
func (h *UserHandler) Login() bool {
	fmt.Println("\n=== LOGIN ===")

	fmt.Print("Enter email: ")
	email, _ := h.Reader.ReadString('\n')
	email = strings.TrimSpace(email)

	fmt.Print("Enter password: ")
	password, _ := h.Reader.ReadString('\n')
	password = strings.TrimSpace(password)

	user, err := h.UserRepo.GetByEmail(email)
	if err != nil {
		if email == "" {
			fmt.Println("Email cannot be empty.")
		} else {
			fmt.Println("Error:", err)
		}
		return false
	}

	if !utils.CheckPassword(password, user.Password) {
		fmt.Println("Password incorrect")
		return false
	}

	fmt.Printf("Welcome, %s! Role: %s\n", user.Name, user.RoleName)
	h.ShowDashboard(&user)
	return true
}

// --- REGISTER USER ---
func (h *UserHandler) RegisterUserCLI() {
	fmt.Println("\n=== REGISTER USER ===")

	fmt.Print("Enter name: ")
	name, _ := h.Reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Enter email: ")
	email, _ := h.Reader.ReadString('\n')
	email = strings.TrimSpace(email)

	fmt.Print("Enter password: ")
	password, _ := h.Reader.ReadString('\n')
	password = strings.TrimSpace(password)

	if name == "" || email == "" || password == "" {
		fmt.Println("All fields are required.")
		return
	}

	hashed := utils.HashPassword(password)
	err := h.UserRepo.RegisterUser(name, email, hashed)
	if err != nil {
		fmt.Println("Failed to register user:", err)
		return
	}

	fmt.Println("User registered successfully!")
}

// --- ADD STAFF (Admin Only) ---
func (h *UserHandler) AddStaff() {
	fmt.Println("\n=== ADD STAFF ===")

	fmt.Print("Enter staff name: ")
	name, _ := h.Reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Enter staff email: ")
	email, _ := h.Reader.ReadString('\n')
	email = strings.TrimSpace(email)

	fmt.Print("Enter password: ")
	password, _ := h.Reader.ReadString('\n')
	password = strings.TrimSpace(password)

	hashed := utils.HashPassword(password)
	err := h.UserRepo.CreateStaff(name, email, hashed)
	if err != nil {
		fmt.Println("Error creating staff:", err)
		return
	}

	fmt.Printf("Staff '%s' (%s) added successfully!\n", name, email)
}

// --- DASHBOARD MENU (Merged version) ---
func (h *UserHandler) ShowDashboard(user *entity.User) {
	for {
		// ERROR: method getMenuItems tidak ada
		// menu := h.getMenuItems(user.RoleName)
		var menu []string
		switch user.RoleName {
		case "admin":
			menu = []string{"See Products", "Manage Orders", "Report", "Add Staff", "Logout"}
		case "staff":
			menu = []string{"See Products", "Manage Orders", "Report", "Logout"}
		case "user":
			menu = []string{"Create Order", "My Orders", "History", "Logout"}
		}

		prompt := promptui.Select{
			Label: fmt.Sprintf("=== DASHBOARD (%s) ===", strings.ToUpper(user.RoleName)),
			Items: menu,
			Templates: &promptui.SelectTemplates{
				Label:    "{{ . | cyan | bold }}",
				Active:   "{{ . | green | bold }}",
				Inactive: "  {{ . | white }}",
				Selected: "{{ . | bold }}",
			},
		}

		i, _, err := prompt.Run()
		if err != nil {
			fmt.Println("Prompt failed:", err)
			return
		}

		// ERROR: variable 'role' tidak ada, harusnya user.RoleName
		switch strings.ToLower(user.RoleName) {
		case "admin":
			switch i {
			case 0:
				fmt.Println("Admin: See Products")
			case 1:
				fmt.Println("Admin: Manage Orders")
			case 2:
				fmt.Println("Admin: Report")
			case 3:
				h.AddStaff()
			case 4:
				fmt.Println("Logging out...")
				return
			}
		case "staff":
			switch i {
			case 0:
				fmt.Println("Staff: See Products")
			case 1:
				fmt.Println("Staff: Manage Orders")
			case 2:
				fmt.Println("Staff: Report")
			case 3:
				fmt.Println("Logging out...")
				return
			}
		case "user":
			switch i {
			case 0:
				fmt.Println("User: Create Order")
			case 1:
				fmt.Println("User: My Orders")
			case 2:
				fmt.Println("User: History")
			case 3:
				fmt.Println("Logging out...")
				return
			}
		default:
			fmt.Println("Invalid role.")
			return
		}
	}
}
