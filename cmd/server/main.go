package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"pet-shelter/internal/handlers"
)

// Модели данных
type Animal = handlers.Animal

// Хранилище в памяти
var animals = []Animal{
	{ID: 1, Name: "Барсик", Species: "Кошка", Age: 3},
	{ID: 2, Name: "Шарик", Species: "Собака", Age: 2},
}

var catalog = []map[string]interface{}{
	{"id": 1, "name": "Сухой корм для собак", "price": 25.99, "quantity": 50},
	{"id": 2, "name": "Игрушка для кошек", "price": 5.99, "quantity": 100},
}

var users = []map[string]interface{}{
	{"id": 1, "username": "admin", "password": "admin123", "role": "admin"},
	{"id": 2, "username": "user", "password": "user123", "role": "user"},
}

func main() {
	// Определяем порт
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// Инициализируем шаблоны
	log.Println("Загрузка шаблонов...")
	if err := handlers.InitTemplates(); err != nil {
		log.Printf("Не удалось загрузить шаблоны: %v", err)
	}

	// Статика
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Страницы (HTML)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// данные для home.html
		data := handlers.HomeData{
			AnimalCount:  len(animals),
			CatalogCount: len(catalog),
			UserCount:    len(users),
			Year:         time.Now().Year(),
			Port:         port,
		}
		handlers.Home(w, r, data)
	})
	// новая страница со списком животных
	http.HandleFunc("/animals", func(w http.ResponseWriter, r *http.Request) {
		handlers.AnimalsPage(w, r, animals)
	})

	// новая страница с каталогом
	http.HandleFunc("/catalog", func(w http.ResponseWriter, r *http.Request) {
		handlers.CatalogPage(w, r, catalog)
	})

	http.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		// данные для admin.html
		handlers.Admin(w, r, users)
	})

	http.HandleFunc("/api/docs", func(w http.ResponseWriter, r *http.Request) {
		handlers.APIDocs(w, r, port)
	})

	// Health check
	http.HandleFunc("/health", healthHandler)

	// API
	http.HandleFunc("/api/animals", animalsHandler)
	http.HandleFunc("/api/animals/", animalHandler)
	http.HandleFunc("/api/catalog", catalogHandler)
	http.HandleFunc("/api/login", loginHandler)
	http.HandleFunc("/api/register", registerHandler)
	http.HandleFunc("/api/admin/animals", createAnimalHandler)
	http.HandleFunc("/api/admin/catalog", createCatalogItemHandler)

	log.Printf("Pet Shelter запущен на порту %s", port)
	go openBrowser(port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// Функция для открытия браузера (как у тебя)
func openBrowser(port string) {
	time.Sleep(2 * time.Second)
	fmt.Printf("\nСервер запущен: http://localhost:%s\n", port)
}

// Health check
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "ok",
		"time":    time.Now().Format(time.RFC3339),
		"version": "1.0.0",
	})
}

// --- API хендлеры ниже оставляем без изменений ---

func animalsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(animals)
}

func animalHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/animals/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for _, animal := range animals {
		if animal.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(animal)
			return
		}
	}
	http.NotFound(w, r)
}

func catalogHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(catalog)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	for _, user := range users {
		if user["username"] == data.Username && user["password"] == data.Password {
			response := map[string]interface{}{
				"token": "admin-token-" + fmt.Sprint(time.Now().Unix()),
				"user": map[string]interface{}{
					"id":       user["id"],
					"username": user["username"],
					"role":     user["role"],
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	http.Error(w, "Неверные учетные данные", http.StatusUnauthorized)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	for _, user := range users {
		if user["username"] == data.Username {
			http.Error(w, "Пользователь уже существует", http.StatusConflict)
			return
		}
	}

	newUser := map[string]interface{}{
		"id":       len(users) + 1,
		"username": data.Username,
		"password": data.Password,
		"role":     "user",
	}

	users = append(users, newUser)

	response := map[string]interface{}{
		"message": "Пользователь создан",
		"user": map[string]interface{}{
			"id":       newUser["id"],
			"username": newUser["username"],
			"role":     newUser["role"],
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func createAnimalHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Admin") != "true" {
		http.Error(w, "Требуются права администратора", http.StatusForbidden)
		return
	}

	var animal Animal
	if err := json.NewDecoder(r.Body).Decode(&animal); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	if animal.Name == "" || animal.Species == "" {
		http.Error(w, "Имя и вид обязательны", http.StatusBadRequest)
		return
	}

	animal.ID = len(animals) + 1
	animals = append(animals, animal)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(animal)
}

func createCatalogItemHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Admin") != "true" {
		http.Error(w, "Требуются права администратора", http.StatusForbidden)
		return
	}

	var item struct {
		Name     string  `json:"name"`
		Price    float64 `json:"price"`
		Quantity int     `json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	newItem := map[string]interface{}{
		"id":       len(catalog) + 1,
		"name":     item.Name,
		"price":    item.Price,
		"quantity": item.Quantity,
	}

	catalog = append(catalog, newItem)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newItem)
}
