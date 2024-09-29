package main

import (
	"net/http"
	"encoding/json"
	"log"
	"github.com/gorilla/mux"
)

type OrderDTO struct {
	OrderID  int    `json:"orderId"`
	Product  string `json:"product"`
	Quantity int    `json:"quantity"`
}

func AddOrder(w http.ResponseWriter, r *http.Request) {
	var order OrderDTO
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	

	response := map[string]interface{}{
		"message": "Order added successfully",
		"order":   order,
	}
	log.Println("order placed  in your basket")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/order", AddOrder).Methods("POST")

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
