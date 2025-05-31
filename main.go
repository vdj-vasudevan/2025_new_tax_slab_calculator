package main

import (
	handlers "2025_new_tax_slab_calculator/handler"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handlers.RenderForm)
	http.HandleFunc("/calculate-tax", handlers.CalculateTax)

	http.HandleFunc("/custom-tax", handlers.RenderCustomForm)
	http.HandleFunc("/calculate-custom-tax", handlers.CalculateCustomTax)

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
