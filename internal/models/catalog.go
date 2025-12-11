package models

import (
	"time"
)

type CatalogItem struct {
	ID          int       `json:"id"`
	Name        string    `json:"name" binding:"required"`
	Type        string    `json:"type" binding:"required"` // корм, игрушка, лекарство и т.д.
	Category    string    `json:"category"`                // для собак, для кошек, универсальный
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
	Description string    `json:"description"`
	Supplier    string    `json:"supplier"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
