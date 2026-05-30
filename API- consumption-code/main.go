package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Product represents minimal product structure
type Product struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// ProductResponse represents API response
type ProductResponse struct {
	Products []Product `json:"products"`
}

// GetProductCount fetches total product count
func GetProductCount() (int, error) {
	url := "https://dummyjson.com/products"

	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to call API: %w", err)
	}
	defer resp.Body.Close()

	//check status
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	//decode json and convert json to struct
	var data ProductResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	return len(data.Products), nil //extract result and return and count products
}

// GetProductsByQuery fetches products based on search query
func GetProductsByQuery(query string) ([]Product, error) {
	url := fmt.Sprintf("https://dummyjson.com/products/search?q=%s", query)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var data ProductResponse // create a variable to hold data. Preparing a container(struct)to store decoded JSON
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return data.Products, nil
}

func main() {
	// 1. Get product count
	count, err := GetProductCount()
	if err != nil {
		log.Println("Error:", err)
		return
	}
	fmt.Println("Total Products:", count)

	// 2. Get laptops
	laptops, err := GetProductsByQuery("laptop")
	if err != nil {
		log.Println("Error:", err)
		return
	}

	fmt.Println("\nLaptops:")
	for _, product := range laptops {
		fmt.Printf("ID: %d, Name: %s\n", product.ID, product.Title)
	}
}
