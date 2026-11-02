package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"ecommerce-order-management/backend/database"
	"ecommerce-order-management/backend/handlers"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Println("Connected to MySQL successfully!")

	mux := http.NewServeMux()

	productHandler := handlers.NewProductHandler(db)

	mux.HandleFunc("GET /api/products", productHandler.GetProducts)
	mux.HandleFunc("GET /api/products/{id}", productHandler.GetProduct)
	mux.HandleFunc("POST /api/products", productHandler.CreateProduct)
	mux.HandleFunc("PUT /api/products/{id}", productHandler.UpdateProduct)
	mux.HandleFunc("DELETE /api/products/{id}", productHandler.DeleteProduct)

	customerHandler := handlers.NewCustomerHandler(db)

	mux.HandleFunc("GET /api/customers", customerHandler.GetCustomers)
	mux.HandleFunc("GET /api/customers/{id}", customerHandler.GetCustomer)
	mux.HandleFunc("POST /api/customers", customerHandler.CreateCustomer)
	mux.HandleFunc("PUT /api/customers/{id}", customerHandler.UpdateCustomer)
	mux.HandleFunc("DELETE /api/customers/{id}", customerHandler.DeleteCustomer)

	orderHandler := handlers.NewOrderHandler(db)

	mux.HandleFunc("POST /api/orders", orderHandler.CreateOrder)
	mux.HandleFunc("GET /api/orders", orderHandler.GetOrders)
	mux.HandleFunc("GET /api/orders/{id}", orderHandler.GetOrder)
	mux.HandleFunc("POST /api/orders/{id}/confirm", orderHandler.ConfirmOrder)
	mux.HandleFunc("PATCH /api/orders/{id}/status", orderHandler.UpdateOrderStatus)

	dashboardHandler := handlers.NewDashboardHandler(db)
	mux.HandleFunc("GET /api/dashboard", dashboardHandler.GetDashboard)

	mux.HandleFunc("/api/health", healthHandler(db))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      corsMiddleware(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Server running at http://localhost:8080")
	log.Fatal(server.ListenAndServe())
}

func healthHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		w.Header().Set("Content-Type", "application/json")

		if err := db.PingContext(ctx); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]string{
				"status":   "error",
				"database": "disconnected",
			})
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"status":   "success",
			"message":  "E-commerce API is running",
			"database": "connected",
		})
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
