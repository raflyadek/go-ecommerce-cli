package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"go-ecommerce-cli/config"
	"go-ecommerce-cli/internal/handler"
	"go-ecommerce-cli/internal/repository"
	"os"

	"github.com/manifoldco/promptui"
)

func waitEnter() {
	fmt.Println("\nPress Enter to continue...")
	bufio.NewReader(os.Stdin).ReadString('\n')
}

func LoginCLI(db *sql.DB) {
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepo(db) // buat order repository

	productHandler := handler.NewProductHandler(productRepo)
	orderHandler := handler.NewOrderHandler(orderRepo) // buat order handler

	// Tambahkan orderHandler ke UserHandler
	userHandler := handler.NewUserHandler(userRepo, db, productHandler, orderHandler)

	fmt.Println("========================================================================")
	fmt.Println("                           WELCOME TO RYD STORE                                     ")
	fmt.Println("========================================================================")

	for {
		menu := []string{
			"Already a member? Sign in!",
			"Not a member? Sign up now!",
			"Exit",
		}

		prompt := promptui.Select{
			Label: "=== MAIN MENU ===",
			Items: menu,
		}

		i, _, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed: %v\n", err)
			return
		}

		switch i {
		case 0:
			waitEnter()
			if userHandler.Login() {
				return
			}
		case 1:
			waitEnter()
			userHandler.RegisterUserCLI()
		case 2:
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid option")
			waitEnter()
		}
	}
}

func main() {
	db := config.ConnectDB()
	defer db.Close()

	LoginCLI(db)
}
