package models

import "time"

type Animal struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Species      string    `json:"species"`
	Breed        string    `json:"breed"`
	Age          int       `json:"age"`
	Gender       string    `json:"gender"`
	HealthStatus string    `json:"health_status"`
	Description  string    `json:"description"`
	Adopted      bool      `json:"adopted"`
	ArrivalDate  time.Time `json:"arrival_date"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
