package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"ecommerce-order-management/backend/models"
)

type OrderHandler struct {
	DB *sql.DB
}

func NewOrderHandler(db *sql.DB) *OrderHandler {
	return &OrderHandler{DB: db}
}

// CreateOrder creates an order and its items in one transaction.
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var request models.CreateOrderRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	if request.CustomerID <= 0 || len(request.Items) == 0 {
		http.Error(w, "A valid customer and at least one order item are required", http.StatusBadRequest)
		return
	}

	// Prevent the same product from appearing more than once.
	seenProducts := make(map[int]bool)

	for _, item := range request.Items {
		if item.ProductID <= 0 || item.Quantity <= 0 {
			http.Error(w, "Product IDs and quantities must be positive", http.StatusBadRequest)
			return
		}

		if seenProducts[item.ProductID] {
			http.Error(w, "Each product can appear only once in an order", http.StatusBadRequest)
			return
		}

		seenProducts[item.ProductID] = true
	}

	// Start a database transaction.
	tx, err := h.DB.BeginTx(r.Context(), nil)
	if err != nil {
		http.Error(w, "Could not start order transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Check whether the customer exists.
	var customerID int
	err = tx.QueryRowContext(
		r.Context(),
		"SELECT id FROM customers WHERE id = ?",
		request.CustomerID,
	).Scan(&customerID)

	if err == sql.ErrNoRows {
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Could not verify customer", http.StatusInternalServerError)
		return
	}

	// Generate a unique order number.
	orderNumber := fmt.Sprintf("ORD-%d", time.Now().UnixNano())

	// Create the order first. Its total will be updated after inserting items.
	result, err := tx.ExecContext(
		r.Context(),
		`INSERT INTO orders (order_number, customer_id, total_amount, status)
		 VALUES (?, ?, 0, 'PLACED')`,
		orderNumber,
		customerID,
	)
	if err != nil {
		http.Error(w, "Could not create order", http.StatusInternalServerError)
		return
	}

	orderID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Could not retrieve order ID", http.StatusInternalServerError)
		return
	}

	var totalAmount float64
	orderItems := make([]models.OrderItem, 0, len(request.Items))

	// Fetch each product's actual price and save its order item.
	for _, item := range request.Items {
		var price float64

		err := tx.QueryRowContext(
			r.Context(),
			"SELECT price FROM products WHERE id = ?",
			item.ProductID,
		).Scan(&price)

		if err == sql.ErrNoRows {
			http.Error(w, fmt.Sprintf("Product %d not found", item.ProductID), http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, "Could not retrieve product price", http.StatusInternalServerError)
			return
		}

		_, err = tx.ExecContext(
			r.Context(),
			`INSERT INTO order_items (order_id, product_id, quantity, unit_price)
			 VALUES (?, ?, ?, ?)`,
			orderID,
			item.ProductID,
			item.Quantity,
			price,
		)
		if err != nil {
			http.Error(w, "Could not save order item", http.StatusInternalServerError)
			return
		}

		subtotal := price * float64(item.Quantity)
		totalAmount += subtotal

		orderItems = append(orderItems, models.OrderItem{
			OrderID:   int(orderID),
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: price,
			Subtotal:  subtotal,
		})
	}

	// Update the order total using the prices retrieved from the database.
	_, err = tx.ExecContext(
		r.Context(),
		"UPDATE orders SET total_amount = ? WHERE id = ?",
		totalAmount,
		orderID,
	)
	if err != nil {
		http.Error(w, "Could not calculate order total", http.StatusInternalServerError)
		return
	}

	// Commit saves the order and all its items together.
	if err := tx.Commit(); err != nil {
		http.Error(w, "Could not complete order transaction", http.StatusInternalServerError)
		return
	}

	order := models.Order{
		ID:          int(orderID),
		OrderNumber: orderNumber,
		CustomerID:  customerID,
		TotalAmount: totalAmount,
		Status:      "PLACED",
		CreatedAt:   time.Now(),
		Items:       orderItems,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)

	_ = strconv.Itoa // Removed in the next cleanup step if unused.
}

// GetOrders returns all orders.
func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.QueryContext(
		r.Context(),
		`SELECT id, order_number, customer_id, total_amount, status, created_at
		 FROM orders
		 ORDER BY created_at DESC`,
	)
	if err != nil {
		http.Error(w, "Could not retrieve orders", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	orders := []models.Order{}

	for rows.Next() {
		var order models.Order

		err := rows.Scan(
			&order.ID,
			&order.OrderNumber,
			&order.CustomerID,
			&order.TotalAmount,
			&order.Status,
			&order.CreatedAt,
		)
		if err != nil {
			http.Error(w, "Could not read order data", http.StatusInternalServerError)
			return
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error while reading orders", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, orders)
}

// GetOrder returns one order with its items.
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	var order models.Order

	err = h.DB.QueryRowContext(
		r.Context(),
		`SELECT id, order_number, customer_id, total_amount, status, created_at
		 FROM orders
		 WHERE id = ?`,
		id,
	).Scan(
		&order.ID,
		&order.OrderNumber,
		&order.CustomerID,
		&order.TotalAmount,
		&order.Status,
		&order.CreatedAt,
	)

	if err == sql.ErrNoRows {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Could not retrieve order", http.StatusInternalServerError)
		return
	}

	rows, err := h.DB.QueryContext(
		r.Context(),
		`SELECT id, order_id, product_id, quantity, unit_price
		 FROM order_items
		 WHERE order_id = ?`,
		id,
	)
	if err != nil {
		http.Error(w, "Could not retrieve order items", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	order.Items = []models.OrderItem{}

	for rows.Next() {
		var item models.OrderItem

		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.Quantity,
			&item.UnitPrice,
		)
		if err != nil {
			http.Error(w, "Could not read order items", http.StatusInternalServerError)
			return
		}

		item.Subtotal = item.UnitPrice * float64(item.Quantity)
		order.Items = append(order.Items, item)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error while reading order items", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, order)
}

// ConfirmOrder checks stock and reserves inventory.
func (h *OrderHandler) ConfirmOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	tx, err := h.DB.BeginTx(r.Context(), nil)
	if err != nil {
		http.Error(w, "Could not start transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Lock the order while confirming it.
	var status string
	err = tx.QueryRowContext(
		r.Context(),
		"SELECT status FROM orders WHERE id = ? FOR UPDATE",
		id,
	).Scan(&status)

	if err == sql.ErrNoRows {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Could not retrieve order status", http.StatusInternalServerError)
		return
	}

	if status != "PLACED" {
		http.Error(w, "Only PLACED orders can be confirmed", http.StatusConflict)
		return
	}

	// Retrieve all products in this order.
	rows, err := tx.QueryContext(
		r.Context(),
		"SELECT product_id, quantity FROM order_items WHERE order_id = ?",
		id,
	)
	if err != nil {
		http.Error(w, "Could not retrieve order items", http.StatusInternalServerError)
		return
	}

	type orderItem struct {
		ProductID int
		Quantity  int
	}

	items := []orderItem{}

	for rows.Next() {
		var item orderItem
		if err := rows.Scan(&item.ProductID, &item.Quantity); err != nil {
			rows.Close()
			http.Error(w, "Could not read order items", http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		rows.Close()
		http.Error(w, "Could not read order items", http.StatusInternalServerError)
		return
	}
	rows.Close()

	if len(items) == 0 {
		http.Error(w, "Order has no items", http.StatusConflict)
		return
	}

	// Lock products in sorted ID order to reduce deadlock risk.
	// First, sort the items by product ID.
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			if items[i].ProductID > items[j].ProductID {
				items[i], items[j] = items[j], items[i]
			}
		}
	}

	// Check stock and reserve it.
	for _, item := range items {
		var stock int

		err := tx.QueryRowContext(
			r.Context(),
			"SELECT stock_quantity FROM products WHERE id = ? FOR UPDATE",
			item.ProductID,
		).Scan(&stock)

		if err == sql.ErrNoRows {
			http.Error(w, "Product no longer exists", http.StatusConflict)
			return
		}
		if err != nil {
			http.Error(w, "Could not check inventory", http.StatusInternalServerError)
			return
		}

		if stock < item.Quantity {
			http.Error(
				w,
				fmt.Sprintf("Insufficient stock for product ID %d. Available: %d, Required: %d", item.ProductID, stock, item.Quantity),
				http.StatusConflict,
			)
			return
		}

		_, err = tx.ExecContext(
			r.Context(),
			"UPDATE products SET stock_quantity = stock_quantity - ? WHERE id = ?",
			item.Quantity,
			item.ProductID,
		)
		if err != nil {
			http.Error(w, "Could not reserve inventory", http.StatusInternalServerError)
			return
		}
	}

	// Change order status after successful stock reservation.
	_, err = tx.ExecContext(
		r.Context(),
		"UPDATE orders SET status = 'CONFIRMED' WHERE id = ?",
		id,
	)
	if err != nil {
		http.Error(w, "Could not update order status", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Could not complete confirmation", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"message":  "Order confirmed and inventory reserved successfully",
		"order_id": id,
		"status":   "CONFIRMED",
	})
}

// UpdateOrderStatus changes an order's status and restores stock if cancelled.
func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	var request struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	tx, err := h.DB.BeginTx(r.Context(), nil)
	if err != nil {
		http.Error(w, "Could not start transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	var currentStatus string

	err = tx.QueryRowContext(
		r.Context(),
		"SELECT status FROM orders WHERE id = ? FOR UPDATE",
		id,
	).Scan(&currentStatus)

	if err == sql.ErrNoRows {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Could not retrieve order", http.StatusInternalServerError)
		return
	}

	// Define valid status transitions.
	validTransitions := map[string]string{
		"CONFIRMED": "PACKED",
		"PACKED":    "SHIPPED",
		"SHIPPED":   "DELIVERED",
	}

	isCancellation := request.Status == "CANCELLED"

	if !isCancellation && validTransitions[currentStatus] != request.Status {
		http.Error(
			w,
			fmt.Sprintf("Cannot change order status from %s to %s", currentStatus, request.Status),
			http.StatusConflict,
		)
		return
	}

	if isCancellation {
		if currentStatus == "SHIPPED" || currentStatus == "DELIVERED" || currentStatus == "CANCELLED" {
			http.Error(w, "This order cannot be cancelled", http.StatusConflict)
			return
		}

		// Restore reserved stock for confirmed or packed orders.
		if currentStatus == "CONFIRMED" || currentStatus == "PACKED" {
			rows, err := tx.QueryContext(
				r.Context(),
				"SELECT product_id, quantity FROM order_items WHERE order_id = ?",
				id,
			)
			if err != nil {
				http.Error(w, "Could not retrieve order items", http.StatusInternalServerError)
				return
			}

			type item struct {
				ProductID int
				Quantity  int
			}

			items := []item{}

			for rows.Next() {
				var current item
				if err := rows.Scan(&current.ProductID, &current.Quantity); err != nil {
					rows.Close()
					http.Error(w, "Could not read order items", http.StatusInternalServerError)
					return
				}
				items = append(items, current)
			}

			if err := rows.Err(); err != nil {
				rows.Close()
				http.Error(w, "Could not read order items", http.StatusInternalServerError)
				return
			}
			rows.Close()

			// Lock products in sorted order.
			for i := 0; i < len(items); i++ {
				for j := i + 1; j < len(items); j++ {
					if items[i].ProductID > items[j].ProductID {
						items[i], items[j] = items[j], items[i]
					}
				}
			}

			for _, current := range items {
				var productID int

				err := tx.QueryRowContext(
					r.Context(),
					"SELECT id FROM products WHERE id = ? FOR UPDATE",
					current.ProductID,
				).Scan(&productID)

				if err != nil {
					http.Error(w, "Could not lock product inventory", http.StatusInternalServerError)
					return
				}

				_, err = tx.ExecContext(
					r.Context(),
					"UPDATE products SET stock_quantity = stock_quantity + ? WHERE id = ?",
					current.Quantity,
					current.ProductID,
				)
				if err != nil {
					http.Error(w, "Could not restore inventory", http.StatusInternalServerError)
					return
				}
			}
		}
	}

	_, err = tx.ExecContext(
		r.Context(),
		"UPDATE orders SET status = ? WHERE id = ?",
		request.Status,
		id,
	)
	if err != nil {
		http.Error(w, "Could not update order status", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Could not complete status update", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"message":  "Order status updated successfully",
		"order_id": id,
		"status":   request.Status,
	})
}
