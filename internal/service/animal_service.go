package service

import (
	"context"
	"errors"
	"strings"

	"pet-shelter/internal/models"
	"pet-shelter/internal/repository"
)

type AnimalService struct {
	repo repository.AnimalRepository
}

func NewAnimalService(repo repository.AnimalRepository) *AnimalService {
	return &AnimalService{repo: repo}
}

func (s *AnimalService) Create(ctx context.Context, animal *models.Animal) error {
	animal.Name = strings.TrimSpace(animal.Name)
	animal.Species = strings.TrimSpace(animal.Species)
	animal.Breed = strings.TrimSpace(animal.Breed)
	animal.Gender = strings.TrimSpace(animal.Gender)
	animal.HealthStatus = strings.TrimSpace(animal.HealthStatus)
	animal.Description = strings.TrimSpace(animal.Description)

	if animal.Name == "" {
		return errors.New("имя животного обязательно")
	}
	if animal.Species == "" {
		return errors.New("вид животного обязателен")
	}
	if animal.Age < 0 {
		return errors.New("возраст не может быть отрицательным")
	}

	return s.repo.Create(ctx, animal)
}

func (s *AnimalService) List(ctx context.Context) ([]models.Animal, error) {
	return s.repo.List(ctx)
}

func (s *AnimalService) GetByID(ctx context.Context, id int) (*models.Animal, error) {
	if id <= 0 {
		return nil, errors.New("некорректный id животного")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *AnimalService) Update(ctx context.Context, animal *models.Animal) error {
	if animal.ID <= 0 {
		return errors.New("некорректный id")
	}

	animal.Name = strings.TrimSpace(animal.Name)
	animal.Species = strings.TrimSpace(animal.Species)

	if animal.Name == "" {
		return errors.New("имя обязательно")
	}

	if animal.Species == "" {
		return errors.New("вид обязателен")
	}

	return s.repo.Update(ctx, animal)
}

func (s *AnimalService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("некорректный id животного")
	}
	return s.repo.Delete(ctx, id)
}

func (s *AnimalService) Count(ctx context.Context) (int, error) {
	return s.repo.Count(ctx)
}
