package main

import (
	"log"
	"net/http"
	"os"
	"pet-shelter/internal/handlers"
)

func main() {
	// Определяем порт
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// Инициализируем обработчики
	handlers.Init()

	// Настройка статических файлов
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Регистрируем маршруты
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/admin", handlers.AdminHandler)
	http.HandleFunc("/api/docs", handlers.APIHandler)
	http.HandleFunc("/health", handlers.HealthHandler)

	// API маршруты
	http.HandleFunc("/api/animals", handlers.AnimalsHandler)
	http.HandleFunc("/api/animals/", handlers.AnimalHandler)
	http.HandleFunc("/api/catalog", handlers.CatalogHandler)
	http.HandleFunc("/api/login", handlers.LoginHandler)
	http.HandleFunc("/api/register", handlers.RegisterHandler)
	http.HandleFunc("/api/admin/animals", handlers.CreateAnimalHandler)
	http.HandleFunc("/api/admin/catalog", handlers.CreateCatalogItemHandler)

	// Запуск сервера
	log.Printf("🚀 Pet Shelter System запущен на порту %s", port)
	log.Printf("🌐 Откройте: http://localhost:%s", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("❌ Ошибка сервера:", err)
	}
}
