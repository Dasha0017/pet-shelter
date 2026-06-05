package service

import (
	"context"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"pet-shelter/internal/models"
	"pet-shelter/internal/repository"
)

func EnsureAdmin(
	ctx context.Context,
	users repository.UserRepository,
	username string,
	password string,
) error {

	_, err := users.FindByUsername(ctx, username)

	if err == nil {
		return nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	admin := &models.User{
		Username: username,
		Email:    username + "@local",
		Password: string(hash),
		Role:     "admin",
	}

	return users.Create(ctx, admin)
}
