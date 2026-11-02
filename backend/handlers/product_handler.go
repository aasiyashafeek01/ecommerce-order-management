package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"ecommerce-order-management/backend/models"
)

type ProductHandler struct {
	DB *sql.DB
}

func NewProductHandler(db *sql.DB) *ProductHandler {
	return &ProductHandler{DB: db}
}

// GET /api/products
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(`
		SELECT id, name, sku, price, stock_quantity, created_at
		FROM products
		ORDER BY id DESC
	`)
	if err != nil {
		http.Error(w, "Failed to retrieve products", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	products := []models.Product{}

	for rows.Next() {
		var p models.Product

		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.SKU,
			&p.Price,
			&p.StockQuantity,
			&p.CreatedAt,
		)
		if err != nil {
			http.Error(w, "Failed to read products", http.StatusInternalServerError)
			return
		}

		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Failed to read products", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, products)
}

// GET /api/products/{id}
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id, ok := getProductID(w, r)
	if !ok {
		return
	}

	var p models.Product

	err := h.DB.QueryRow(`
		SELECT id, name, sku, price, stock_quantity, created_at
		FROM products
		WHERE id = ?
	`, id).Scan(
		&p.ID,
		&p.Name,
		&p.SKU,
		&p.Price,
		&p.StockQuantity,
		&p.CreatedAt,
	)

	if err == sql.ErrNoRows {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to retrieve product", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, p)
}

// POST /api/products
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var p models.Product

	if !decodeJSON(w, r, &p) {
		return
	}

	p.Name = strings.TrimSpace(p.Name)
	p.SKU = strings.TrimSpace(p.SKU)

	if p.Name == "" || p.SKU == "" || p.Price <= 0 || p.StockQuantity < 0 {
		http.Error(w, "Valid name, SKU, positive price and non-negative stock are required", http.StatusBadRequest)
		return
	}

	result, err := h.DB.Exec(`
		INSERT INTO products (name, sku, price, stock_quantity)
		VALUES (?, ?, ?, ?)
	`, p.Name, p.SKU, p.Price, p.StockQuantity)

	if err != nil {
		http.Error(w, "Failed to create product. Check that the SKU is unique.", http.StatusConflict)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Product created but its ID could not be retrieved", http.StatusInternalServerError)
		return
	}

	p.ID = int(id)

	err = h.DB.QueryRow(`
		SELECT created_at FROM products WHERE id = ?
	`, p.ID).Scan(&p.CreatedAt)

	if err != nil {
		http.Error(w, "Product created but its details could not be retrieved", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, p)
}

// PUT /api/products/{id}
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, ok := getProductID(w, r)
	if !ok {
		return
	}

	var p models.Product

	if !decodeJSON(w, r, &p) {
		return
	}

	p.Name = strings.TrimSpace(p.Name)
	p.SKU = strings.TrimSpace(p.SKU)

	if p.Name == "" || p.SKU == "" || p.Price <= 0 || p.StockQuantity < 0 {
		http.Error(w, "Valid name, SKU, positive price and non-negative stock are required", http.StatusBadRequest)
		return
	}

	result, err := h.DB.Exec(`
		UPDATE products
		SET name = ?, sku = ?, price = ?, stock_quantity = ?
		WHERE id = ?
	`, p.Name, p.SKU, p.Price, p.StockQuantity, id)

	if err != nil {
		http.Error(w, "Failed to update product. Check that the SKU is unique.", http.StatusConflict)
		return
	}

	affected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to verify product update", http.StatusInternalServerError)
		return
	}

	if affected == 0 {
		var exists bool
		err := h.DB.QueryRow(
			"SELECT EXISTS(SELECT 1 FROM products WHERE id = ?)", id,
		).Scan(&exists)

		if err != nil {
			http.Error(w, "Failed to verify product", http.StatusInternalServerError)
			return
		}
		if !exists {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Product updated successfully",
	})
}

// DELETE /api/products/{id}
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, ok := getProductID(w, r)
	if !ok {
		return
	}

	result, err := h.DB.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Cannot delete product. It may be used in an order.", http.StatusConflict)
		return
	}

	affected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to verify product deletion", http.StatusInternalServerError)
		return
	}

	if affected == 0 {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Product deleted successfully",
	})
}

// Helper functions

func getProductID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return 0, false
	}
	return id, true
}

func decodeJSON(w http.ResponseWriter, r *http.Request, destination any) bool {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(destination); err != nil {
		http.Error(w, "Invalid JSON request body", http.StatusBadRequest)
		return false
	}

	return true
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
