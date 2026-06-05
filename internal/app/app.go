package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"

	"pet-shelter/internal/config"
	httpserver "pet-shelter/internal/http"
	"pet-shelter/internal/http/handlers"
	"pet-shelter/internal/http/middleware"
	"pet-shelter/internal/repository/postgres"
	"pet-shelter/internal/service"
)

type App struct {
	cfg    config.Config
	db     *sql.DB
	server *http.Server
}

func New() (*App, error) {
	cfg := config.Load()
	if cfg.JWTSecret == "" {
		return nil, errors.New("JWT_SECRET is required")
	}

	db, err := postgres.NewDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}

	animalRepo := postgres.NewAnimalRepository(db)
	catalogRepo := postgres.NewCatalogRepository(db)
	userRepo := postgres.NewUserRepository(db)

	if err := service.EnsureAdmin(
		context.Background(),
		userRepo,
		cfg.AdminUsername,
		cfg.AdminPassword,
	); err != nil {
		return nil, fmt.Errorf("create admin: %w", err)
	}

	animalService := service.NewAnimalService(animalRepo)
	catalogService := service.NewCatalogService(catalogRepo)
	authService := service.NewAuthService(userRepo, service.AuthOptions{
		JWTSecret:     cfg.JWTSecret,
		AdminUsername: cfg.AdminUsername,
		AdminPassword: cfg.AdminPassword,
	})

	renderer := handlers.NewRenderer("web/templates")
	handler := handlers.New(renderer, animalService, catalogService, authService, cfg)
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)

	router := httpserver.NewRouter(handler, authMiddleware, cfg)
	server := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
	}

	return &App{cfg: cfg, db: db, server: server}, nil
}

func (a *App) Run() error {
	defer a.db.Close()

	host := localIP()
	log.Println("======================================")
	log.Printf("Pet Shelter запущен на порту %s", a.cfg.ServerPort)
	log.Printf("Локально: http://localhost:%s", a.cfg.ServerPort)
	log.Printf("В сети LAN: http://%s:%s", host, a.cfg.ServerPort)
	if a.cfg.IsDev() {
		log.Printf("Dev auto login: http://localhost:%s/dev-auto-login", a.cfg.ServerPort)
	}
	log.Println("Нажми Ctrl+C для остановки сервера")
	log.Println("======================================")

	if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func localIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "127.0.0.1"
	}
	defer conn.Close()

	addr, ok := conn.LocalAddr().(*net.UDPAddr)
	if !ok {
		return "127.0.0.1"
	}
	return addr.IP.String()
}
