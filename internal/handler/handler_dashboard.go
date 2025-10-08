package handler

import (
	"fmt"
	"go-ecommerce-cli/internal/entity"
	"strings"

	"github.com/manifoldco/promptui"
)

type DashboardHandler struct{}

func (h *DashboardHandler) ShowDashboard(user *entity.User) {
	for {
		var menu []string
		switch user.RoleName {
		case "admin":
			menu = []string{
				"See Products",
				"Manage Orders",
				"Report",
				"Add Staff",
				"Logout",
			}
		case "staff":
			menu = []string{
				"See Products",
				"Manage Orders",
				"Report",
				"Logout",
			}
		case "user":
			menu = []string{
				"Create Order",
				"My Orders",
				"History",
				"Logout",
			}
		default:
			fmt.Println("Unknown role")
			return
		}

		prompt := promptui.Select{
			Label: "=== Dashboard ===",
			Items: menu,
		}

		i, _, err := prompt.Run()
		if err != nil {
			fmt.Println("Prompt failed:", err)
			return
		}

		choice := strings.TrimSpace(fmt.Sprint(i + 1))
		if h.handleChoice(user, choice) {
			break
		}
	}
}

func (h *DashboardHandler) handleChoice(user *entity.User, choice string) bool {
	switch user.RoleName {
	case "admin":
		return h.handleAdminMenu(choice)
	case "staff":
		return h.handleStaffMenu(choice)
	case "user":
		return h.handleUserMenu(choice)
	default:
		fmt.Println("Unknown role")
		return false
	}
}

func (h *DashboardHandler) handleAdminMenu(choice string) bool {
	switch choice {
	case "1":
		fmt.Println("Showing all products...")
	case "2":
		fmt.Println("Managing orders...")
	case "3":
		fmt.Println("Opening report menu...")
	case "4":
		fmt.Println("Adding staff...")
	case "5":
		fmt.Println("Logged out.")
		return true
	default:
		fmt.Println("Invalid choice.")
	}
	return false
}

func (h *DashboardHandler) handleStaffMenu(choice string) bool {
	switch choice {
	case "1":
		fmt.Println("Showing all products...")
	case "2":
		fmt.Println("Managing orders...")
	case "3":
		fmt.Println("Opening report menu...")
	case "4":
		fmt.Println("Logged out.")
		return true
	default:
		fmt.Println("Invalid choice.")
	}
	return false
}

func (h *DashboardHandler) handleUserMenu(choice string) bool {
	switch choice {
	case "1":
		fmt.Println("Creating new order...")
	case "2":
		fmt.Println("Showing my orders...")
	case "3":
		fmt.Println("Showing order history...")
	case "4":
		fmt.Println("Logged out.")
		return true
	default:
		fmt.Println("Invalid choice.")
	}
	return false
}
