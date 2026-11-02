package models

import "time"

type Order struct {
	ID          int         `json:"id"`
	OrderNumber string      `json:"order_number"`
	CustomerID  int         `json:"customer_id"`
	TotalAmount float64     `json:"total_amount"`
	Status      string      `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
	Items       []OrderItem `json:"items,omitempty"`
}

type OrderItem struct {
	ID        int     `json:"id"`
	OrderID   int     `json:"order_id"`
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
	Subtotal  float64 `json:"subtotal"`
}

type CreateOrderRequest struct {
	CustomerID int               `json:"customer_id"`
	Items      []CreateOrderItem `json:"items"`
}

type CreateOrderItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}
