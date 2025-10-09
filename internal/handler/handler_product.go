package handler

import (
	"fmt"
	"go-ecommerce-cli/internal/repository"
	"os"

	"github.com/olekukonko/tablewriter"
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

	// Membuat tabel
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Price (Rp)", "Description"})

	// Tambahkan data ke tabel
	for _, p := range products {
		row := []string{
			fmt.Sprintf("%d", p.ID),
			p.Name,
			fmt.Sprintf("%.2f", p.Price),
			p.Description,
		}
		table.Append(row)
	}

	// Render tabel ke terminal
	table.Render()
}