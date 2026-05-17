package models

import "time"

type CatalogItem struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Category    string    `json:"category"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
	Description string    `json:"description"`
	Supplier    string    `json:"supplier"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
