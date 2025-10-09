package handler

import (
	"bufio"
	"fmt"
	// "go-ecommerce-cli/internal/entity"
	"go-ecommerce-cli/internal/repository"
	"log"
	"os"
	"strings"

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
	table.SetHeader([]string{"ID", "Name", "Price (Rp)", "Description, Stock"})

	// Tambahkan data ke tabel
	for _, p := range products {
		row := []string{
			fmt.Sprintf("%d", p.ID),
			p.Name,
			fmt.Sprintf("%.2f", p.Price),
			p.Description,
			fmt.Sprintf("%d", p.Stock),
		}
		table.Append(row)
	}

	// Render tabel ke terminal
	table.Render()
}

func (h *ProductHandler) GetProductById(id int) {
	// fmt.Print("Masukkan id: ")
	// var id int 
	// fmt.Scanln(&id)
	product, err := h.ProductRepo.ShowProductById(id)

	if err != nil {
		log.Println("Error fetching data from table product", err)
		return
	}
	fmt.Printf("id: %d \nname: %s \ndescription: %s \nprice: %.2f \np.stock %d\n", product.ID, product.Name, product.Description, product.Price, product.Stock)
}

func (h *ProductHandler) AddProducts() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter product name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	
	fmt.Print("Enter product description: ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	fmt.Print("Enter price: ")
	var price float64
	if _, err := fmt.Scanln(&price); err != nil {
		fmt.Println("Please insert integer: ", err)
		return
	}

	fmt.Print("Enter stock: ")
	var stock int
	if _, err := fmt.Scanln(&stock); err != nil {
		fmt.Println("Please insert integer: ", err)
		return
	}

	err := h.ProductRepo.AddProduct(name, description, price, stock)
	if err != nil {
		log.Println("error add data to table product: ", err)
		return
	}

	fmt.Println("Success add product!")
}

