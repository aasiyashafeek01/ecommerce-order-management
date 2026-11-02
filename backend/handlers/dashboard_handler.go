package handlers

import (
	"database/sql"
	"net/http"
)

type DashboardHandler struct {
	DB *sql.DB
}

func NewDashboardHandler(db *sql.DB) *DashboardHandler {
	return &DashboardHandler{DB: db}
}

func (h *DashboardHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	var totalOrders int
	var totalOrderValue float64
	var pendingOrders int
	var totalProducts int
	var totalCustomers int
	var lowStockProducts int
	var inventoryValue float64

	// Order statistics. Cancelled orders are excluded from order value.
	err := h.DB.QueryRowContext(
		r.Context(),
		`SELECT
			COUNT(*),
			COALESCE(SUM(CASE WHEN status != 'CANCELLED' THEN total_amount ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status IN ('PLACED', 'CONFIRMED', 'PACKED') THEN 1 ELSE 0 END), 0)
		FROM orders`,
	).Scan(&totalOrders, &totalOrderValue, &pendingOrders)

	if err != nil {
		http.Error(w, "Could not retrieve order statistics", http.StatusInternalServerError)
		return
	}

	// Product statistics.
	err = h.DB.QueryRowContext(
		r.Context(),
		"SELECT COUNT(*), COALESCE(SUM(CASE WHEN stock_quantity <= 5 THEN 1 ELSE 0 END), 0), COALESCE(SUM(price * stock_quantity), 0) FROM products",
	).Scan(&totalProducts, &lowStockProducts, &inventoryValue)

	if err != nil {
		http.Error(w, "Could not retrieve product statistics", http.StatusInternalServerError)
		return
	}

	// Customer statistics.
	err = h.DB.QueryRowContext(
		r.Context(),
		"SELECT COUNT(*) FROM customers",
	).Scan(&totalCustomers)

	if err != nil {
		http.Error(w, "Could not retrieve customer statistics", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"total_orders":       totalOrders,
		"total_order_value":  totalOrderValue,
		"pending_orders":     pendingOrders,
		"total_products":     totalProducts,
		"total_customers":    totalCustomers,
		"low_stock_products": lowStockProducts,
		"inventory_value":    inventoryValue,
	})
}
