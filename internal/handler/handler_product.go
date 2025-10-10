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
	table.SetHeader([]string{"ID", "Name", "Price (Rp)", "Description", "Stock"})

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

// show product by id
func (h *ProductHandler) GetProductById(id int) {
	product, err := h.ProductRepo.ShowProductById(id)

	if err != nil {
		log.Println("Error fetching data from table product", err)
		return
	}
	fmt.Printf("Id: %d \nName: %s \nDescription: %s \nPrice: %.2f \nStock %d\n\n", product.ID, product.Name, product.Description, product.Price, product.Stock)
}


func (h *ProductHandler) AddProducts() {
	reader := bufio.NewReader(os.Stdin)

	// read product name input
	fmt.Print("Enter product name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	// read product description input
	fmt.Print("Enter product description: ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)
	// read price input
	fmt.Print("Enter price: ")
	var price float64
	if _, err := fmt.Scanln(&price); err != nil {
		fmt.Println("Please insert integer: ", err)
		return
	}
	// read stock input
	fmt.Print("Enter stock: ")
	var stock int
	if _, err := fmt.Scanln(&stock); err != nil {
		fmt.Println("Please insert integer: ", err)
		return
	}

	// execute addproduct to db
	err := h.ProductRepo.AddProduct(name, description, price, stock)
	if err != nil {
		log.Println("error add data to table product: ", err)
		return
	}
	fmt.Println("Success add product!")
}

func (h *ProductHandler) DeleteProducts() {
	reader := bufio.NewReader(os.Stdin)
	// read product id input
	fmt.Print("Enter product ID to delete: ")
	var id int 
	if _, err := fmt.Scanln(&id); err != nil {
		fmt.Println("Please insert integer: ", err)
		return
	}

	fmt.Println("\nSelected Product:")
	// get product info based on input id 
	h.GetProductById(id)
	// looping not close until user input between y or n
	for {
		fmt.Print("Are you sure you want to delete this product? (y/n): ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)
		switch choice {
		case "y":
			//go to statement to jump to the label
			goto delete
		case "n":
			fmt.Println("Not deleted anything")
			return
		default:
			fmt.Println("please insert y or n")
		} 

	}

	// label to continue from goto statement
	delete:

	// execute delete product by id
	err := h.ProductRepo.DeleteProduct(id)
	if err != nil {
		log.Println("Error execute from product table")
		return
	}

	fmt.Println("\nSuccess delete product!")
}

func (h *ProductHandler) UpdateProducts() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Product ID to update (or type '0' to cancel): ")
	var id int 
	/* loop if user input 0 then return to menu if other than 0 
	it execute the function update */
	for {
		if _, err := fmt.Scanln(&id); err != nil {
			fmt.Println("Please insert integer: ", err)
			return
		}
		if id == 0 {
			return
		} else {
			break
		}
	}
	fmt.Println("\nSelected Product:")
	// show product based on input id
	h.GetProductById(id)

	// read input name
	fmt.Print("Enter new name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	// read input description
	fmt.Print("Enter new description: ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	// read input price
	fmt.Print("Enter new price: ")
	var price float64 
	if _, err := fmt.Scanln(&price); err != nil {
		fmt.Println("Please insert integer: ", err)
		return
	}

	// read input stock
	fmt.Print("Enter new stock: ")
	var stock int 
	if _, err := fmt.Scanln(&stock); err != nil {
		fmt.Println("Please insert integer: ", err)
		return
	}

	// execute function UpdateProduct
	err := h.ProductRepo.UpdateProduct(name, description, price, stock, id)

	if err != nil {
		log.Println("error update data to table product: ", err)
		return
	}

	fmt.Println("Success update product!")
}

