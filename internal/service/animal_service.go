package service

import (
	"errors"
	"pet-shelter/internal/models"
	"pet-shelter/internal/repository"
)

type AnimalService struct {
	repo *repository.AnimalRepository
}

func NewAnimalService(repo *repository.AnimalRepository) *AnimalService {
	return &AnimalService{repo: repo}
}

func (s *AnimalService) Create(animal *models.Animal) error {
	// Базовая валидация
	if animal.Name == "" {
		return errors.New("animal name is required")
	}
	if animal.Species == "" {
		return errors.New("animal species is required")
	}

	return s.repo.Create(animal)
}

func (s *AnimalService) GetAll() ([]models.Animal, error) {
	return s.repo.GetAll()
}

func (s *AnimalService) GetByID(id int) (*models.Animal, error) {
	return s.repo.GetByID(id)
}

func (s *AnimalService) Update(animal *models.Animal) error {
	// Проверяем, существует ли животное
	existing, err := s.repo.GetByID(animal.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("animal not found")
	}

	return s.repo.Update(animal)
}

func (s *AnimalService) Delete(id int) error {
	return s.repo.Delete(id)
}
