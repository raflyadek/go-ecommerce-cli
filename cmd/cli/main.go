package main

import (
	"go-ecommerce-cli/config"
)

func main() {
	db := config.ConnectDB()
	defer db.Close()
}
