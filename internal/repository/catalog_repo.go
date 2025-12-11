package repository

import (
	"database/sql"
	"pet-shelter/internal/models"
)

type CatalogRepository struct {
	db *sql.DB
}

func NewCatalogRepository(db *sql.DB) *CatalogRepository {
	return &CatalogRepository{db: db}
}

func (r *CatalogRepository) Create(item *models.CatalogItem) error {
	query := `INSERT INTO catalog_items (name, type, category, price, 
              quantity, description, supplier) 
              VALUES ($1, $2, $3, $4, $5, $6, $7) 
              RETURNING id, created_at, updated_at`

	return r.db.QueryRow(query,
		item.Name, item.Type, item.Category, item.Price,
		item.Quantity, item.Description, item.Supplier,
	).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
}

func (r *CatalogRepository) GetAll() ([]models.CatalogItem, error) {
	query := `SELECT id, name, type, category, price, quantity, 
              description, supplier, created_at, updated_at 
              FROM catalog_items ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.CatalogItem
	for rows.Next() {
		var item models.CatalogItem
		err := rows.Scan(
			&item.ID, &item.Name, &item.Type, &item.Category,
			&item.Price, &item.Quantity, &item.Description,
			&item.Supplier, &item.CreatedAt, &item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *CatalogRepository) GetByID(id int) (*models.CatalogItem, error) {
	item := &models.CatalogItem{}
	query := `SELECT id, name, type, category, price, quantity, 
              description, supplier, created_at, updated_at 
              FROM catalog_items WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&item.ID, &item.Name, &item.Type, &item.Category,
		&item.Price, &item.Quantity, &item.Description,
		&item.Supplier, &item.CreatedAt, &item.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return item, nil
}

func (r *CatalogRepository) Update(item *models.CatalogItem) error {
	query := `UPDATE catalog_items SET 
              name = $1, type = $2, category = $3, price = $4, 
              quantity = $5, description = $6, supplier = $7 
              WHERE id = $8 RETURNING updated_at`

	return r.db.QueryRow(query,
		item.Name, item.Type, item.Category, item.Price,
		item.Quantity, item.Description, item.Supplier, item.ID,
	).Scan(&item.UpdatedAt)
}

func (r *CatalogRepository) Delete(id int) error {
	query := `DELETE FROM catalog_items WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
