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
	productHandler := handler.NewProductHandler(productRepo)
	userHandler := handler.NewUserHandler(userRepo, db, productHandler)

	fmt.Println("========================================================================")
	fmt.Println(`     
 \o       o/      o           o__ __o        o__ __o__/_   o__ __o               o__ __o    ____o__ __o____     o__ __o        o__ __o         o__ __o__/_ 
  v\     /v      <|>         <|     v\      <|    v       <|     v\             /v     v\    /   \   /   \     /v     v\      <|     v\       <|    v      
   <\   />       / \         / \     <\     < >           / \     <\           />       <\        \o/         />       <\     / \     <\      < >          
     \o/       o/   \o       \o/       \o    |            \o/     o/          _\o____              |        o/           \o   \o/     o/       |           
      |       <|__ __|>       |         |>   o__/_         |__  _<|                \_\__o__       < >      <|             |>   |__  _<|        o__/_       
     / \      /       \      / \       //    |             |       \                     \         |        \\           //    |       \       |           
     \o/    o/         \o    \o/      /     <o>           <o>       \o         \         /         o          \         /     <o>       \o    <o>          
      |    /v           v\    |      o       |             |         v\         o       o         <|           o       o       |         v\    |           
     / \  />             <\  / \  __/>      / \  _\o__/_  / \         <\        <\__ __/>         / \          <\__ __/>      / \         <\  / \  _\o__/_ 
                                                                                                                                                           
                                                                                                                                                           
                                                                                                                                                                                                `)
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
			// waitEnter()
			if userHandler.Login() {
				return
			}
		case 1:
			// waitEnter()
			userHandler.RegisterUserCLI()
		case 2:
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid option")
			// waitEnter()
		}
	}
}

func main() {
	db := config.ConnectDB()
	defer db.Close()

	LoginCLI(db)
}
