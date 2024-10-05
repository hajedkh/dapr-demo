package main

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "sync"

    "github.com/gorilla/mux"
)


// Entity
type OrderDTO struct {
    OrderID  int      `json:"orderId"`
    Articles []string `json:"articleIds"` 
    Quantity int      `json:"quantity"`
}

var (
    orders []OrderDTO
    mu     sync.Mutex
)

func SaveOrders() error {
    mu.Lock()
    defer mu.Unlock()

    data, err := json.Marshal(orders)
    if err != nil {
        return err
    }

    return ioutil.WriteFile("orders.json", data, 0644)
}

func LoadOrders() error {
    mu.Lock()
    defer mu.Unlock()

    file, err := ioutil.ReadFile("orders.json")
    if err != nil {
        if os.IsNotExist(err) {
            return nil 
        }
        return err
    }

    return json.Unmarshal(file, &orders)
}


func AddOrder(w http.ResponseWriter, r *http.Request) {
    var order OrderDTO
    err := json.NewDecoder(r.Body).Decode(&order)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }


    orders = append(orders, order)

    if err := SaveOrders(); err != nil {
        http.Error(w, "Failed to save order", http.StatusInternalServerError)
        return
    }

    response := map[string]interface{}{
        "message": "Order added successfully",
        "order":   order,
    }
    log.Println("Order placed in your basket:", order)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}


func GetOrders(w http.ResponseWriter, r *http.Request) {

    if err := LoadOrders(); err != nil {
			http.Error(w, "Failed to load orders", http.StatusInternalServerError)
			return
		}


    mu.Lock()
    defer mu.Unlock()

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(orders); err != nil {
        http.Error(w, "Failed to encode orders", http.StatusInternalServerError)
        return
    }
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/order", AddOrder).Methods("POST")
	r.HandleFunc("/orders", GetOrders).Methods("GET")

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
