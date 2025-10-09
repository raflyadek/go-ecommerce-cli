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
	UserRepo       *repository.UserRepository
	ProductHandler *ProductHandler

	Reader         *bufio.Reader
	DB             *sql.DB
}

// Tambahkan orderHandler di constructor
func NewUserHandler(userRepo *repository.UserRepository, db *sql.DB, productHandler *ProductHandler) *UserHandler {
	return &UserHandler{
		UserRepo:       userRepo,
		ProductHandler: productHandler,
		Reader:         bufio.NewReader(os.Stdin),
		DB:             db,
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
	h.ProductHandler.ShowAllProducts()
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

// --- DASHBOARD MENU ---
func (h *UserHandler) ShowDashboard(user *entity.User) {
	for {
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

		switch strings.ToLower(user.RoleName) {
		case "admin":
			switch i {
			case 0:
				SeeProductsMenu(h.ProductHandler)
			case 1:
				ManageOrdersMenu()
			case 2:
				h.ReportMenu()
			case 3:
				h.AddStaff()
			case 4:
				fmt.Println("Logging out...")
				return
			}
		case "staff":
			switch i {
			case 0:
				SeeProductsMenu(h.ProductHandler)
			case 1:
				ManageOrdersMenu()
			case 2:
				h.ReportMenu()
			case 3:
				fmt.Println("Logging out...")
				return
			}
		case "user":
			switch i {
			case 0:
				// h.OrderHandler.CreateOrderCLI(user.ID)
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

func (h *UserHandler) ReportMenu() {
	for {
		fmt.Println("\n===== Report Menu =====")

		menu := []string{
			"User Report",
			"Order Report",
			"Stock Report",
			"Back to Dashboard",
		}

		prompt := promptui.Select{
			Label: "Select Report Type",
			Items: menu,
		}

		i, _, err := prompt.Run()
		if err != nil {
			fmt.Println("Prompt failed:", err)
			return
		}

		switch i {
		case 0:
			reportHandler := NewReportHandler(h.DB)
			reportHandler.ShowUserReport()
		case 1:
			reportHandler := NewReportHandler(h.DB)
			reportHandler.ShowCompletedOrders()
		case 2:
			reportHandler := NewReportHandler(h.DB)
			reportHandler.ShowStockReport()
		case 3:
			return
		}
	}
}

// --- PRODUCT MENU ---
func SeeProductsMenu(ProductHandler *ProductHandler) {
	for {
		fmt.Println("\n====================================================")
		ProductHandler.ShowAllProducts()
		fmt.Println("====================================================")

		menu := []string{
			"Add Product",
			"Update Product",
			"Delete Product",
			"Back to Dashboard",
		}

		prompt := promptui.Select{
			Label: "Select Action",
			Items: menu,
			Templates: &promptui.SelectTemplates{
				Label:    "{{ . | cyan | bold }}",
				Active:   "> {{ . | green | bold }}",
				Inactive: "  {{ . | white }}",
				Selected: "{{ . | bold }}",
			},
		}

		i, _, err := prompt.Run()
		if err != nil {
			fmt.Println("Prompt failed:", err)
			return
		}

		switch i {
		case 0:
			ProductHandler.AddProducts()
		case 1:
			ProductHandler.UpdateProducts()
		case 2:
			ProductHandler.DeleteProducts()
		case 3:
			return
		}
	}
}

// --- MANAGE ORDERS MENU ---
func ManageOrdersMenu() {
	for {
		fmt.Println("\n=== Manage Orders ===")

		menu := []string{
			"View All Orders",
			"View Order Details",
			"Update Order Status",
			"Back to Dashboard",
		}

		prompt := promptui.Select{
			Label: "Select Action",
			Items: menu,
			Templates: &promptui.SelectTemplates{
				Label:    "{{ . | cyan | bold }}",
				Active:   "> {{ . | green | bold }}",
				Inactive: "  {{ . | white }}",
				Selected: "{{ . | bold }}",
			},
		}

		i, _, err := prompt.Run()
		if err != nil {
			fmt.Println("Prompt failed:", err)
			return
		}

		switch i {
		case 0:
			fmt.Println("View All Orders - TBD")
		case 1:
			fmt.Println("View Order Details - TBD")
		case 2:
			fmt.Println("Update Order Status - TBD")
		case 3:
			return
		}
	}
}
