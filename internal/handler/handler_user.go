package handler

import (
	"bufio"
	"database/sql"
	"fmt"
	"go-ecommerce-cli/internal/entity"
	"go-ecommerce-cli/internal/repository"
	"go-ecommerce-cli/pkg/utils"
	"io"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"

	"github.com/manifoldco/promptui"
)

type UserHandler struct {
    UserRepo repository.IUserRepository
    ProductHandler *ProductHandler
    OrderHandler *OrderHandler
    Reader         *bufio.Reader
    DB             *sql.DB
    Writer io.Writer
}

// Tambahkan orderHandler di constructor
func NewUserHandler(userRepo repository.IUserRepository, db *sql.DB, productHandler *ProductHandler, orderHandler *OrderHandler) *UserHandler {
    return &UserHandler{
        UserRepo:       userRepo,
        ProductHandler: productHandler,
        OrderHandler: orderHandler,
        Reader:         bufio.NewReader(os.Stdin),
        DB:             db,
        Writer: os.Stdout,
    }
}

// --- LOGIN ---
func (h *UserHandler) Login() *entity.User {
	fmt.Println("\n=== LOGIN ===")

	fmt.Print("Enter email: ")
	email, _ := h.Reader.ReadString('\n')
	email = strings.TrimSpace(email)

	fmt.Print("Enter password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		fmt.Println("Error reading password:", err)
		return nil
	}
	password := strings.TrimSpace(string(bytePassword))

    user, err := h.LoginLogic(email, password)
    if err != nil {
        fmt.Println("Error:", err)
        return nil
    }

    fmt.Printf("Welcome, %s! Role: %s\n", user.Name, user.RoleName)
    h.ProductHandler.ShowAllProducts()
    return user
}

// LoginLogic handles the core authentication logic without CLI interaction.
func (h *UserHandler) LoginLogic(email, password string) (*entity.User, error) {
    user, err := h.UserRepo.GetByEmail(email)
    if err != nil {
        return nil, fmt.Errorf("user not found")
    }

    if !utils.CheckPassword(password, user.Password) {
        return nil, fmt.Errorf("password incorrect")
    }

    return &user, nil
}

// --- REGISTER USER ---
func (h *UserHandler) RegisterUserCLI() {
	fmt.Fprintln(h.Writer, "\n=== REGISTER USER ===")

	fmt.Fprint(h.Writer, "Enter name: ")
	name, _ := h.Reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Fprint(h.Writer, "Enter email: ")
	email, _ := h.Reader.ReadString('\n')
	email = strings.TrimSpace(email)

	fmt.Fprint(h.Writer, "Enter password: ")
	password, _ := h.Reader.ReadString('\n')
	password = strings.TrimSpace(password)

	if name == "" || email == "" || password == "" {
		fmt.Fprintln(h.Writer, "All fields are required.")
		return
	}

	hashed := utils.HashPassword(password)
	err := h.UserRepo.RegisterUser(name, email, hashed)
	if err != nil {
		fmt.Fprintln(h.Writer, "Failed to register user:", err)
		return
	}

	fmt.Fprintln(h.Writer, "Registration successful")
}

// --- ADD STAFF (Admin Only) ---
func (h *UserHandler) AddStaff() { fmt.Println("\n=== ADD STAFF ===") 
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

    if err != nil { fmt.Println("Error creating staff:", err) 
        return 
    } 
    fmt.Printf("Staff '%s' (%s) added successfully!\n", name, email) 
}

// --- DASHBOARD MENU ---
func (h *UserHandler) ShowDashboard(user *entity.User) bool {
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
            return true
        }

        switch strings.ToLower(user.RoleName) {
        case "admin":
            switch i {
            case 0:
                SeeProductsMenu(h.ProductHandler)
            case 1:
                ManageOrdersMenu(h.OrderHandler)
            case 2:
                h.ReportMenu()
            case 3:
                h.AddStaff()
            case 4:
                return true
            }
        case "staff":
            switch i {
            case 0:
                SeeProductsMenu(h.ProductHandler)
            case 1:
                ManageOrdersMenu(h.OrderHandler)
            case 2:
                h.ReportMenu()
            case 3:
                fmt.Println("Logging out...")
                return true
            }
        case "user":
            switch i {
            case 0:
                h.OrderHandler.CreateOrderCLI(user.ID)
            case 1:
                h.OrderHandler.MyOrdersCLI(user.ID)
            case 2:
                h.OrderHandler.HistoryOrdersCLI(user.ID)
            case 3:
                fmt.Println("Logging out...")
                return true
            }
        default:
            fmt.Println("Invalid role.")
            return true
        }
    }
}

func (h *UserHandler) ReportMenu() {
    // Tampilkan best selling products dulu
    reportHandler := NewReportHandler(h.DB)
    reportHandler.ShowBestSellingProducts()
    
    for {
        fmt.Println("\n===== Report Menu =====")

        menu := []string{
            "User Report",
            "Order Report",
            "Stock Report",
            "Report Summary",
            "Back to Dashboard",
        }

        prompt := promptui.Select{
            Label: "Select Report Type",
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
            reportHandler := NewReportHandler(h.DB)
            reportHandler.ShowReportSummary()
        case 4:
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
func ManageOrdersMenu(orderHandler *OrderHandler) {
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
            orderHandler.AllOrdersCLI()
        case 1:
            orderHandler.ViewOrderDetailsCLI()
        case 2:
            orderHandler.AllOrdersCLI()
            orderHandler.UpdateOrderStatusCLI()
        case 3:
            return
        }
    }
}
