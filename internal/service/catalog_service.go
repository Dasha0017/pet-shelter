package service

import (
	"errors"
	"pet-shelter/internal/models"
	"pet-shelter/internal/repository"
)

type CatalogService struct {
	repo *repository.CatalogRepository
}

func NewCatalogService(repo *repository.CatalogRepository) *CatalogService {
	return &CatalogService{repo: repo}
}

func (s *CatalogService) Create(item *models.CatalogItem) error {
	if item.Name == "" {
		return errors.New("item name is required")
	}
	if item.Type == "" {
		return errors.New("item type is required")
	}
	if item.Price < 0 {
		return errors.New("price cannot be negative")
	}

	return s.repo.Create(item)
}

func (s *CatalogService) GetAll() ([]models.CatalogItem, error) {
	return s.repo.GetAll()
}

func (s *CatalogService) GetByID(id int) (*models.CatalogItem, error) {
	return s.repo.GetByID(id)
}

func (s *CatalogService) Update(item *models.CatalogItem) error {
	existing, err := s.repo.GetByID(item.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("item not found")
	}

	return s.repo.Update(item)
}

func (s *CatalogService) Delete(id int) error {
	return s.repo.Delete(id)
}
