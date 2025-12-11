package service

import (
	"errors"
	"pet-shelter/internal/models"
	"pet-shelter/internal/repository"
	"pet-shelter/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Register(user *models.User) error {
	// Проверяем, существует ли пользователь
	existingUser, err := s.userRepo.GetUserByUsername(user.Username)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("username already exists")
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// Устанавливаем роль по умолчанию
	if user.Role == "" {
		user.Role = "user"
	}

	// Создаем пользователя
	return s.userRepo.CreateUser(user)
}

func (s *AuthService) Login(username, password string) (*models.LoginResponse, error) {
	// Получаем пользователя
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	// Проверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Генерируем токен
	token, err := jwt.GenerateToken(user.ID, user.Username, user.Role, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *AuthService) ValidateToken(token string) (*jwt.Claims, error) {
	return jwt.ValidateToken(token, s.jwtSecret)
}
