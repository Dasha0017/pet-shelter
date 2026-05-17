package postgres

import (
	"context"
	"database/sql"

	"pet-shelter/internal/models"
)

type CatalogRepository struct {
	db *sql.DB
}

func NewCatalogRepository(db *sql.DB) *CatalogRepository {
	return &CatalogRepository{db: db}
}

func (r *CatalogRepository) Create(ctx context.Context, item *models.CatalogItem) error {
	query := `
		INSERT INTO catalog_items (name, type, category, price, quantity, description, supplier)
		VALUES ($1, $2, NULLIF($3, ''), $4, $5, NULLIF($6, ''), NULLIF($7, ''))
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		item.Name,
		item.Type,
		item.Category,
		item.Price,
		item.Quantity,
		item.Description,
		item.Supplier,
	).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
}

func (r *CatalogRepository) List(ctx context.Context) ([]models.CatalogItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, type, category, price, quantity, description, supplier, created_at, updated_at
		FROM catalog_items
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.CatalogItem, 0)
	for rows.Next() {
		item, err := scanCatalogItem(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *CatalogRepository) GetByID(ctx context.Context, id int) (*models.CatalogItem, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, name, type, category, price, quantity, description, supplier, created_at, updated_at
		FROM catalog_items
		WHERE id = $1
	`, id)

	item, err := scanCatalogItem(row)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *CatalogRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM catalog_items WHERE id = $1`, id)
	return err
}

func (r *CatalogRepository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM catalog_items`).Scan(&count)
	return count, err
}

type catalogScanner interface {
	Scan(dest ...any) error
}

func scanCatalogItem(scanner catalogScanner) (models.CatalogItem, error) {
	var item models.CatalogItem
	var category sql.NullString
	var price sql.NullFloat64
	var quantity sql.NullInt64
	var description sql.NullString
	var supplier sql.NullString

	err := scanner.Scan(
		&item.ID,
		&item.Name,
		&item.Type,
		&category,
		&price,
		&quantity,
		&description,
		&supplier,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		return item, err
	}

	item.Category = category.String
	if price.Valid {
		item.Price = price.Float64
	}
	if quantity.Valid {
		item.Quantity = int(quantity.Int64)
	}
	item.Description = description.String
	item.Supplier = supplier.String

	return item, nil
}
