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
    OrderID  string      `json:"orderId"`
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
    // Define a CloudEvent structure
    var cloudEvent struct {
        Data json.RawMessage `json:"data"` // CloudEvent "data" attribute
    }

    // Parse the CloudEvent from the request body
    if err := json.NewDecoder(r.Body).Decode(&cloudEvent); err != nil {
        http.Error(w, "Invalid CloudEvent payload", http.StatusBadRequest)
        log.Println("Failed to decode CloudEvent:", err)
        return
    }

    // Decode the "data" field into the OrderDTO struct
    var order OrderDTO
    if err := json.Unmarshal(cloudEvent.Data, &order); err != nil {
        http.Error(w, "Invalid order payload in data attribute", http.StatusBadRequest)
        log.Println("Failed to decode OrderDTO from data:", err)
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
    externalServiceURL := "http://payment-ms:5000/pay" // Replace with the actual URL

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

// Middleware to handle CORS
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	r := mux.NewRouter()

	// Apply CORS middleware
	r.Use(enableCORS)

	r.HandleFunc("/order", AddOrder).Methods("POST", "OPTIONS")
	r.HandleFunc("/orders", GetOrders).Methods("GET", "OPTIONS")
	r.HandleFunc("/payOrder", PayOrder).Methods("POST", "OPTIONS")

	log.Println("Server running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}
