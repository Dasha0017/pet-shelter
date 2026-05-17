package postgres

import (
	"context"
	"database/sql"

	"pet-shelter/internal/models"
)

type AnimalRepository struct {
	db *sql.DB
}

func NewAnimalRepository(db *sql.DB) *AnimalRepository {
	return &AnimalRepository{db: db}
}

func (r *AnimalRepository) Create(ctx context.Context, animal *models.Animal) error {
	query := `
		INSERT INTO animals (name, species, breed, age, gender, health_status, description, adopted)
		VALUES ($1, $2, NULLIF($3, ''), $4, NULLIF($5, ''), NULLIF($6, ''), NULLIF($7, ''), $8)
		RETURNING id, arrival_date, created_at, updated_at
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		animal.Name,
		animal.Species,
		animal.Breed,
		animal.Age,
		animal.Gender,
		animal.HealthStatus,
		animal.Description,
		animal.Adopted,
	).Scan(&animal.ID, &animal.ArrivalDate, &animal.CreatedAt, &animal.UpdatedAt)
}

func (r *AnimalRepository) List(ctx context.Context) ([]models.Animal, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, species, breed, age, gender, health_status, description,
		       adopted, arrival_date, created_at, updated_at
		FROM animals
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	animals := make([]models.Animal, 0)
	for rows.Next() {
		animal, err := scanAnimal(rows)
		if err != nil {
			return nil, err
		}
		animals = append(animals, animal)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return animals, nil
}

func (r *AnimalRepository) GetByID(ctx context.Context, id int) (*models.Animal, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, name, species, breed, age, gender, health_status, description,
		       adopted, arrival_date, created_at, updated_at
		FROM animals
		WHERE id = $1
	`, id)

	animal, err := scanAnimal(row)
	if err != nil {
		return nil, err
	}

	return &animal, nil
}

func (r *AnimalRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM animals WHERE id = $1`, id)
	return err
}

func (r *AnimalRepository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM animals`).Scan(&count)
	return count, err
}

type animalScanner interface {
	Scan(dest ...any) error
}

func scanAnimal(scanner animalScanner) (models.Animal, error) {
	var animal models.Animal
	var breed sql.NullString
	var age sql.NullInt64
	var gender sql.NullString
	var healthStatus sql.NullString
	var description sql.NullString

	err := scanner.Scan(
		&animal.ID,
		&animal.Name,
		&animal.Species,
		&breed,
		&age,
		&gender,
		&healthStatus,
		&description,
		&animal.Adopted,
		&animal.ArrivalDate,
		&animal.CreatedAt,
		&animal.UpdatedAt,
	)
	if err != nil {
		return animal, err
	}

	animal.Breed = breed.String
	if age.Valid {
		animal.Age = int(age.Int64)
	}
	animal.Gender = gender.String
	animal.HealthStatus = healthStatus.String
	animal.Description = description.String

	return animal, nil
}
