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
	userHandler := handler.NewUserHandler(userRepo)
	productRepo := repository.NewProductRepository(db)
	productHandler := handler.NewProductHandler(productRepo)
	productHandler.ShowAllProducts()

	fmt.Println("===========================================================================================")
	fmt.Println("                                  WELCOME TO RYD STORE                                     ")
	fmt.Println("===========================================================================================")
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
			// waitEnter()
			if userHandler.Login() {
				return
			}
		case 1:
			waitEnter() 
			userHandler.RegisterUserCLI()
			
		case 2: // Exit
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid option")
			waitEnter()
		}
	}
}

func ShowProductCLI(db *sql.DB){

}

func main() {
	db := config.ConnectDB()
	defer db.Close()
	LoginCLI(db)
}
