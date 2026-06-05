package postgres

import (
	"context"
	"database/sql"

	"pet-shelter/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, username, email, password, role, created_at
		FROM users
		WHERE username = $1
	`, username)

	user, err := scanUser(row)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) List(ctx context.Context) ([]models.User, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, username, email, password, role, created_at
		FROM users
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`).Scan(&count)
	return count, err
}

type userScanner interface {
	Scan(dest ...any) error
}

func scanUser(scanner userScanner) (models.User, error) {
	var user models.User
	var email sql.NullString
	var role sql.NullString

	err := scanner.Scan(&user.ID, &user.Username, &email, &user.Password, &role, &user.CreatedAt)
	if err != nil {
		return user, err
	}

	user.Email = email.String
	if role.Valid {
		user.Role = role.String
	}

	return user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO users (
			username,
			email,
			password,
			role
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`,
		user.Username,
		user.Email,
		user.Password,
		user.Role,
	).Scan(&user.ID, &user.CreatedAt)
}
