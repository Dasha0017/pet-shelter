package repository

import (
	"database/sql"
	"pet-shelter/internal/models"
)

type AnimalRepository struct {
	db *sql.DB
}

func NewAnimalRepository(db *sql.DB) *AnimalRepository {
	return &AnimalRepository{db: db}
}

func (r *AnimalRepository) Create(animal *models.Animal) error {
	query := `INSERT INTO animals (name, species, breed, age, gender, 
              health_status, description, adopted, arrival_date) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
              RETURNING id, created_at, updated_at`

	return r.db.QueryRow(query,
		animal.Name, animal.Species, animal.Breed, animal.Age,
		animal.Gender, animal.HealthStatus, animal.Description,
		animal.Adopted, animal.ArrivalDate,
	).Scan(&animal.ID, &animal.CreatedAt, &animal.UpdatedAt)
}

func (r *AnimalRepository) GetAll() ([]models.Animal, error) {
	query := `SELECT id, name, species, breed, age, gender, health_status, 
              description, adopted, arrival_date, created_at, updated_at 
              FROM animals ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var animals []models.Animal
	for rows.Next() {
		var animal models.Animal
		err := rows.Scan(
			&animal.ID, &animal.Name, &animal.Species, &animal.Breed,
			&animal.Age, &animal.Gender, &animal.HealthStatus,
			&animal.Description, &animal.Adopted, &animal.ArrivalDate,
			&animal.CreatedAt, &animal.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		animals = append(animals, animal)
	}

	return animals, nil
}

func (r *AnimalRepository) GetByID(id int) (*models.Animal, error) {
	animal := &models.Animal{}
	query := `SELECT id, name, species, breed, age, gender, health_status, 
              description, adopted, arrival_date, created_at, updated_at 
              FROM animals WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&animal.ID, &animal.Name, &animal.Species, &animal.Breed,
		&animal.Age, &animal.Gender, &animal.HealthStatus,
		&animal.Description, &animal.Adopted, &animal.ArrivalDate,
		&animal.CreatedAt, &animal.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return animal, nil
}

func (r *AnimalRepository) Update(animal *models.Animal) error {
	query := `UPDATE animals SET 
              name = $1, species = $2, breed = $3, age = $4, 
              gender = $5, health_status = $6, description = $7, 
              adopted = $8, arrival_date = $9 
              WHERE id = $10 RETURNING updated_at`

	return r.db.QueryRow(query,
		animal.Name, animal.Species, animal.Breed, animal.Age,
		animal.Gender, animal.HealthStatus, animal.Description,
		animal.Adopted, animal.ArrivalDate, animal.ID,
	).Scan(&animal.UpdatedAt)
}

func (r *AnimalRepository) Delete(id int) error {
	query := `DELETE FROM animals WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
