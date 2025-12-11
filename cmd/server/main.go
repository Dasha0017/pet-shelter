package main

import (
	"log"
	"net/http"
	"pet-shelter/internal/config"
	"pet-shelter/internal/handlers"
	"pet-shelter/internal/middleware"
	"pet-shelter/internal/repository"
	"pet-shelter/internal/service"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.Load()

	// Подключаемся к базе данных
	db, err := repository.NewDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Инициализируем репозитории
	userRepo := repository.NewUserRepository(db)
	animalRepo := repository.NewAnimalRepository(db)
	catalogRepo := repository.NewCatalogRepository(db)

	// Инициализируем сервисы
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	animalService := service.NewAnimalService(animalRepo)
	catalogService := service.NewCatalogService(catalogRepo)

	// Инициализируем обработчики
	animalHandler := handlers.NewAnimalHandler(animalService)
	catalogHandler := handlers.NewCatalogHandler(catalogService)
	authHandler := handlers.NewAuthHandler(authService)

	// Создаем маршрутизатор
	r := mux.NewRouter()

	// Публичные маршруты
	r.HandleFunc("/api/login", authHandler.Login).Methods("POST")
	r.HandleFunc("/api/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/api/animals", animalHandler.GetAnimals).Methods("GET")
	r.HandleFunc("/api/animals/{id}", animalHandler.GetAnimal).Methods("GET")
	r.HandleFunc("/api/catalog", catalogHandler.GetItems).Methods("GET")
	r.HandleFunc("/api/catalog/{id}", catalogHandler.GetItem).Methods("GET")

	// Защищенные маршруты для администраторов
	admin := r.PathPrefix("/api/admin").Subrouter()
	admin.Use(middleware.AuthMiddleware(authService))
	admin.Use(middleware.AdminOnly)

	admin.HandleFunc("/animals", animalHandler.CreateAnimal).Methods("POST")
	admin.HandleFunc("/animals/{id}", animalHandler.UpdateAnimal).Methods("PUT")
	admin.HandleFunc("/animals/{id}", animalHandler.DeleteAnimal).Methods("DELETE")
	admin.HandleFunc("/catalog", catalogHandler.CreateItem).Methods("POST")
	admin.HandleFunc("/catalog/{id}", catalogHandler.UpdateItem).Methods("PUT")
	admin.HandleFunc("/catalog/{id}", catalogHandler.DeleteItem).Methods("DELETE")

	// Запускаем сервер
	log.Printf("Server starting on port %s", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, r))
}
