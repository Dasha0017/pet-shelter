package models

import (
	"time"
)

type Animal struct {
	ID           int       `json:"id"`
	Name         string    `json:"name" binding:"required"`
	Species      string    `json:"species" binding:"required"` // собака, кошка и т.д.
	Breed        string    `json:"breed"`
	Age          int       `json:"age"`
	Gender       string    `json:"gender"` // male/female
	HealthStatus string    `json:"health_status"`
	Description  string    `json:"description"`
	Adopted      bool      `json:"adopted"`
	ArrivalDate  time.Time `json:"arrival_date"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
