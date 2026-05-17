package service

import (
	"context"
	"errors"
	"strings"

	"pet-shelter/internal/models"
	"pet-shelter/internal/repository"
)

type CatalogService struct {
	repo repository.CatalogRepository
}

func NewCatalogService(repo repository.CatalogRepository) *CatalogService {
	return &CatalogService{repo: repo}
}

func (s *CatalogService) Create(ctx context.Context, item *models.CatalogItem) error {
	item.Name = strings.TrimSpace(item.Name)
	item.Type = strings.TrimSpace(item.Type)
	item.Category = strings.TrimSpace(item.Category)
	item.Description = strings.TrimSpace(item.Description)
	item.Supplier = strings.TrimSpace(item.Supplier)

	if item.Name == "" {
		return errors.New("название товара обязательно")
	}
	if item.Type == "" {
		return errors.New("тип товара обязателен")
	}
	if item.Price < 0 {
		return errors.New("цена не может быть отрицательной")
	}
	if item.Quantity < 0 {
		return errors.New("количество не может быть отрицательным")
	}

	return s.repo.Create(ctx, item)
}

func (s *CatalogService) List(ctx context.Context) ([]models.CatalogItem, error) {
	return s.repo.List(ctx)
}

func (s *CatalogService) GetByID(ctx context.Context, id int) (*models.CatalogItem, error) {
	if id <= 0 {
		return nil, errors.New("некорректный id товара")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *CatalogService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("некорректный id товара")
	}
	return s.repo.Delete(ctx, id)
}

func (s *CatalogService) Count(ctx context.Context) (int, error) {
	return s.repo.Count(ctx)
}
