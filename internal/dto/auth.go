package dto

import "pet-shelter/internal/models"

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type LoginResponse struct {
	Status string      `json:"status"`
	Token  string      `json:"token,omitempty"`
	User   models.User `json:"user"`
}
