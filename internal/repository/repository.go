package repository

import (
	"context"

	"pet-shelter/internal/models"
)

type AnimalRepository interface {
	Create(ctx context.Context, animal *models.Animal) error
	List(ctx context.Context) ([]models.Animal, error)
	GetByID(ctx context.Context, id int) (*models.Animal, error)
	Delete(ctx context.Context, id int) error
	Count(ctx context.Context) (int, error)
}

type CatalogRepository interface {
	Create(ctx context.Context, item *models.CatalogItem) error
	List(ctx context.Context) ([]models.CatalogItem, error)
	GetByID(ctx context.Context, id int) (*models.CatalogItem, error)
	Delete(ctx context.Context, id int) error
	Count(ctx context.Context) (int, error)
}

type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	List(ctx context.Context) ([]models.User, error)
	Count(ctx context.Context) (int, error)
}
