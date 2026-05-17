package service

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"pet-shelter/internal/models"
	"pet-shelter/internal/repository"
	jwtutil "pet-shelter/pkg/jwt"
)

type AuthOptions struct {
	JWTSecret     string
	AdminUsername string
	AdminPassword string
}

type AuthService struct {
	users repository.UserRepository
	opts  AuthOptions
}

func NewAuthService(users repository.UserRepository, opts AuthOptions) *AuthService {
	return &AuthService{users: users, opts: opts}
}

func (s *AuthService) Login(ctx context.Context, username, password, requestedRole string) (string, models.User, error) {
	username = strings.TrimSpace(username)
	requestedRole = strings.TrimSpace(requestedRole)

	if requestedRole == "" {
		requestedRole = "admin"
	}

	if username == "" {
		return "", models.User{}, errors.New("username обязателен")
	}

	if requestedRole == "user" {
		user := models.User{
			ID:       2,
			Username: username,
			Role:     "user",
		}

		token, err := jwtutil.GenerateToken(
			user.ID,
			user.Username,
			user.Role,
			s.opts.JWTSecret,
		)

		return token, user, err
	}

	user, err := s.users.FindByUsername(ctx, username)
	if err != nil {
		return "", models.User{}, errors.New("неверный логин или пароль")
	}

	if user.Role != "admin" {
		return "", models.User{}, errors.New("доступ разрешён только администратору")
	}

	// Сравнение хеша из БД с обычным паролем
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	); err != nil {
		return "", models.User{}, errors.New("неверный логин или пароль")
	}

	token, err := jwtutil.GenerateToken(
		user.ID,
		user.Username,
		user.Role,
		s.opts.JWTSecret,
	)

	return token, *user, err
}

func (s *AuthService) DevAdminToken() (string, error) {
	return jwtutil.GenerateToken(1, s.opts.AdminUsername, "admin", s.opts.JWTSecret)
}

func (s *AuthService) ListUsers(ctx context.Context) ([]models.User, error) {
	return s.users.List(ctx)
}

func (s *AuthService) UserCount(ctx context.Context) (int, error) {
	return s.users.Count(ctx)
}
