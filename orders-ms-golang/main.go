package main

import (
	"bytes"
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
    OrderID  int      `json:"order_id"`
    Articles []string `json:"article_ids"` 
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
func PayOrder(w http.ResponseWriter, r *http.Request) {
    var order OrderDTO
    err := json.NewDecoder(r.Body).Decode(&order)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Assuming the external payment service URL is defined here
    externalServiceURL := "http://localhost:5001/pay" // Replace with the actual URL

    // Convert the order to JSON for the request body
    orderData, err := json.Marshal(order)
    if err != nil {
        http.Error(w, "Failed to marshal order data", http.StatusInternalServerError)
        return
    }

    // Create a POST request to the external service
    resp, err := http.Post(externalServiceURL, "application/json", bytes.NewBuffer(orderData))
    if err != nil {
        http.Error(w, "Failed to communicate with the payment service", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    // Read the response from the external service
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        http.Error(w, "Failed to read payment service response", http.StatusInternalServerError)
        return
    }

    // Return the response from the payment service to the client
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(resp.StatusCode)
    w.Write(body)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/order", AddOrder).Methods("POST")
	r.HandleFunc("/orders", GetOrders).Methods("GET")
	r.HandleFunc("/payOrder", PayOrder).Methods("POST") 

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
