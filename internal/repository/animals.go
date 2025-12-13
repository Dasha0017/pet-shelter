package repository

import (
	"database/sql"

	"pet-shelter/internal/handlers"
)

// Получить всех животных
func GetAllAnimals(db *sql.DB) ([]handlers.Animal, error) {
	rows, err := db.Query(`SELECT id, name, species, age FROM animals ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var animals []handlers.Animal
	for rows.Next() {
		var a handlers.Animal
		if err := rows.Scan(&a.ID, &a.Name, &a.Species, &a.Age); err != nil {
			return nil, err
		}
		animals = append(animals, a)
	}
	return animals, rows.Err()
}

// Добавить животное, вернуть новый id
func CreateAnimal(db *sql.DB, a *handlers.Animal) error {
	return db.QueryRow(
		`INSERT INTO animals (name, species, age) VALUES ($1, $2, $3) RETURNING id`,
		a.Name, a.Species, a.Age,
	).Scan(&a.ID)
}

// Удалить животное по id
func DeleteAnimal(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM animals WHERE id = $1`, id)
	return err
}
