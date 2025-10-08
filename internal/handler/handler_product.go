package handler

import (
	"fmt"
	"go-ecommerce-cli/internal/repository"
)

type ProductHandler struct {
	ProductRepo *repository.ProductRepository
}

func NewProductHandler(repo *repository.ProductRepository) *ProductHandler {
	return &ProductHandler{ProductRepo: repo}
}

func (h *ProductHandler) ShowAllProducts() {
	products, err := h.ProductRepo.GetAll()
	if err != nil {
		fmt.Println("Failed to fetch products:", err)
		return
	}

	if len(products) == 0 {
		fmt.Println("No products found.")
		return
	}

	fmt.Println("\n=== PRODUCTS ===")
	fmt.Printf("%-3s | %-20s | %-10s | %-30s\n", "ID", "Name", "Price", "Description")
	fmt.Println("----------------------------------------------------------------------------------------")
	for _, p := range products {
		fmt.Printf("%-3d | %-20s | %-10.2f | %-30s\n", p.ID, p.Name, p.Price, p.Description)
	}
	fmt.Println()
}
