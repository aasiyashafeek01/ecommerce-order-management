package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"ecommerce-order-management/backend/models"
)

type CustomerHandler struct {
	DB *sql.DB
}

func NewCustomerHandler(db *sql.DB) *CustomerHandler {
	return &CustomerHandler{DB: db}
}

// GET /api/customers
func (h *CustomerHandler) GetCustomers(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(`
		SELECT id, name, email, phone, created_at
		FROM customers
		ORDER BY id DESC
	`)
	if err != nil {
		http.Error(w, "Failed to retrieve customers", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	customers := []models.Customer{}

	for rows.Next() {
		var c models.Customer
		var phone sql.NullString

		err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Email,
			&phone,
			&c.CreatedAt,
		)
		if err != nil {
			http.Error(w, "Failed to read customers", http.StatusInternalServerError)
			return
		}

		if phone.Valid {
			c.Phone = &phone.String
		}

		customers = append(customers, c)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Failed to read customers", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, customers)
}

// GET /api/customers/{id}
func (h *CustomerHandler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	id, ok := getProductID(w, r)
	if !ok {
		return
	}

	var c models.Customer
	var phone sql.NullString

	err := h.DB.QueryRow(`
		SELECT id, name, email, phone, created_at
		FROM customers
		WHERE id = ?
	`, id).Scan(
		&c.ID,
		&c.Name,
		&c.Email,
		&phone,
		&c.CreatedAt,
	)

	if err == sql.ErrNoRows {
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to retrieve customer", http.StatusInternalServerError)
		return
	}

	if phone.Valid {
		c.Phone = &phone.String
	}

	writeJSON(w, http.StatusOK, c)
}

// POST /api/customers
func (h *CustomerHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var c models.Customer

	if !decodeJSON(w, r, &c) {
		return
	}

	c.Name = strings.TrimSpace(c.Name)
	c.Email = strings.TrimSpace(c.Email)

	if c.Name == "" || c.Email == "" {
		http.Error(w, "Customer name and email are required", http.StatusBadRequest)
		return
	}

	if c.Phone != nil {
		phone := strings.TrimSpace(*c.Phone)
		c.Phone = &phone
	}

	result, err := h.DB.Exec(`
		INSERT INTO customers (name, email, phone)
		VALUES (?, ?, ?)
	`, c.Name, c.Email, c.Phone)

	if err != nil {
		http.Error(w, "Failed to create customer. Check whether the email already exists.", http.StatusConflict)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Customer created but ID could not be retrieved", http.StatusInternalServerError)
		return
	}

	c.ID = int(id)

	var phone sql.NullString
	err = h.DB.QueryRow(`
		SELECT phone, created_at FROM customers WHERE id = ?
	`, c.ID).Scan(&phone, &c.CreatedAt)

	if err != nil {
		http.Error(w, "Customer created but details could not be retrieved", http.StatusInternalServerError)
		return
	}

	if phone.Valid {
		c.Phone = &phone.String
	} else {
		c.Phone = nil
	}

	writeJSON(w, http.StatusCreated, c)
}

// PUT /api/customers/{id}
func (h *CustomerHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	id, ok := getProductID(w, r)
	if !ok {
		return
	}

	var c models.Customer

	if !decodeJSON(w, r, &c) {
		return
	}

	c.Name = strings.TrimSpace(c.Name)
	c.Email = strings.TrimSpace(c.Email)

	if c.Name == "" || c.Email == "" {
		http.Error(w, "Customer name and email are required", http.StatusBadRequest)
		return
	}

	if c.Phone != nil {
		phone := strings.TrimSpace(*c.Phone)
		c.Phone = &phone
	}

	result, err := h.DB.Exec(`
		UPDATE customers
		SET name = ?, email = ?, phone = ?
		WHERE id = ?
	`, c.Name, c.Email, c.Phone, id)

	if err != nil {
		http.Error(w, "Failed to update customer. Check whether the email already exists.", http.StatusConflict)
		return
	}

	affected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to verify customer update", http.StatusInternalServerError)
		return
	}

	if affected == 0 {
		var exists bool
		err := h.DB.QueryRow(
			"SELECT EXISTS(SELECT 1 FROM customers WHERE id = ?)", id,
		).Scan(&exists)

		if err != nil {
			http.Error(w, "Failed to verify customer", http.StatusInternalServerError)
			return
		}
		if !exists {
			http.Error(w, "Customer not found", http.StatusNotFound)
			return
		}
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Customer updated successfully",
	})
}

// DELETE /api/customers/{id}
func (h *CustomerHandler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	id, ok := getProductID(w, r)
	if !ok {
		return
	}

	result, err := h.DB.Exec("DELETE FROM customers WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Cannot delete customer. The customer may have existing orders.", http.StatusConflict)
		return
	}

	affected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to verify customer deletion", http.StatusInternalServerError)
		return
	}

	if affected == 0 {
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Customer deleted successfully",
	})
}
