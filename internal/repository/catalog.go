package repository

import (
	"database/sql"
)

type CatalogItem struct {
	ID       int
	Name     string
	Type     string
	Category string
	Price    float64
	Quantity int
}

func GetAllCatalogItems(db *sql.DB) ([]CatalogItem, error) {
	rows, err := db.Query(`SELECT id, name, type, category, price, quantity FROM catalog_items ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []CatalogItem
	for rows.Next() {
		var it CatalogItem
		if err := rows.Scan(&it.ID, &it.Name, &it.Type, &it.Category, &it.Price, &it.Quantity); err != nil {
			return nil, err
		}
		items = append(items, it)
	}
	return items, rows.Err()
}

func CreateCatalogItem(db *sql.DB, it *CatalogItem) error {
	return db.QueryRow(
		`INSERT INTO catalog_items (name, type, category, price, quantity)
         VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		it.Name, it.Type, it.Category, it.Price, it.Quantity,
	).Scan(&it.ID)
}

func DeleteCatalogItem(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM catalog_items WHERE id = $1`, id)
	return err
}
