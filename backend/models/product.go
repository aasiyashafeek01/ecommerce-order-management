package models

import "time"

type Product struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	SKU           string    `json:"sku"`
	Price         float64   `json:"price"`
	StockQuantity int       `json:"stock_quantity"`
	CreatedAt     time.Time `json:"created_at"`
}
