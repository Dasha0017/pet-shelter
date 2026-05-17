package handlers

import (
	"pet-shelter/internal/config"
	"pet-shelter/internal/service"
)

type Handler struct {
	renderer *Renderer
	animals  *service.AnimalService
	catalog  *service.CatalogService
	auth     *service.AuthService
	cfg      config.Config
}

func New(
	renderer *Renderer,
	animals *service.AnimalService,
	catalog *service.CatalogService,
	auth *service.AuthService,
	cfg config.Config,
) *Handler {
	return &Handler{
		renderer: renderer,
		animals:  animals,
		catalog:  catalog,
		auth:     auth,
		cfg:      cfg,
	}
}
