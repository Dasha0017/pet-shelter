package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

// Структуры данных
type Animal struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Species string `json:"species"`
	Age     int    `json:"age"`
}

// Хранилище в памяти
var (
	animals = []Animal{
		{ID: 1, Name: "Барсик", Species: "Кошка", Age: 3},
		{ID: 2, Name: "Шарик", Species: "Собака", Age: 2},
	}

	catalog = []map[string]interface{}{
		{"id": 1, "name": "Сухой корм для собак", "price": 25.99, "quantity": 50},
		{"id": 2, "name": "Игрушка для кошек", "price": 5.99, "quantity": 100},
	}

	users = []map[string]interface{}{
		{"id": 1, "username": "admin", "password": "admin123", "role": "admin"},
		{"id": 2, "username": "user", "password": "user123", "role": "user"},
	}

	templates map[string]string
)

// TemplateData содержит данные для шаблонов
type TemplateData struct {
	Title        string
	Year         int
	Port         string
	AnimalCount  int
	CatalogCount int
	UserCount    int
	AdminCount   int
	Users        []map[string]interface{}
}

// Init инициализирует обработчики
func Init() {
	log.Println("🔄 Инициализация обработчиков...")
	loadTemplates()
}

// loadTemplates загружает HTML шаблоны из файлов
func loadTemplates() {
	templates = make(map[string]string)

	templateFiles := []string{"home", "admin", "api"}

	for _, name := range templateFiles {
		path := filepath.Join("web", "templates", name+".html")
		content, err := ioutil.ReadFile(path)
		if err != nil {
			log.Printf("⚠️  Файл %s не найден, создаем шаблон по умолчанию", path)
			templates[name] = getDefaultTemplate(name)
		} else {
			templates[name] = string(content)
			log.Printf("✅ Загружен шаблон: %s", name)
		}
	}
}

// getDefaultTemplate возвращает шаблон по умолчанию если файл не найден
func getDefaultTemplate(name string) string {
	switch name {
	case "home":
		return `<!DOCTYPE html>
<html>
<head>
    <title>Pet Shelter - Главная</title>
    <style>
        body { font-family: Arial; padding: 20px; }
        h1 { color: #4361ee; }
    </style>
</head>
<body>
    <h1>🐾 Pet Shelter Management System</h1>
    <p>Животных в приюте: {{.AnimalCount}}</p>
    <p>Товаров в каталоге: {{.CatalogCount}}</p>
    <p>Пользователей: {{.UserCount}}</p>
    <p><a href="/admin">Админ панель</a> | <a href="/api/docs">API Docs</a></p>
</body>
</html>`

	case "admin":
		return `<!DOCTYPE html>
<html>
<head>
    <title>Панель администратора</title>
</head>
<body>
    <h1>👑 Панель администратора</h1>
    <p>Всего пользователей: {{.UserCount}}</p>
    <p><a href="/">На главную</a></p>
</body>
</html>`

	case "api":
		return `<!DOCTYPE html>
<html>
<head>
    <title>Документация API</title>
</head>
<body>
    <h1>📚 Документация API</h1>
    <p>Порт сервера: {{.Port}}</p>
    <p><a href="/">На главную</a></p>
</body>
</html>`

	default:
		return `<h1>Шаблон</h1>`
	}
}

// renderTemplate рендерит HTML шаблон
func renderTemplate(w http.ResponseWriter, tmplName string, data TemplateData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, ok := templates[tmplName]
	if !ok {
		http.Error(w, "Шаблон не найден", http.StatusInternalServerError)
		return
	}

	// Используем Go шаблоны для подстановки данных
	t, err := template.New(tmplName).Parse(tmpl)
	if err != nil {
		log.Printf("Ошибка парсинга шаблона: %v", err)
		http.Error(w, "Ошибка шаблона", http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		log.Printf("Ошибка рендеринга: %v", err)
	}
}

// ============ ВЕБ-ОБРАБОТЧИКИ ============

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	data := TemplateData{
		Title:        "Главная",
		Year:         time.Now().Year(),
		Port:         "8081",
		AnimalCount:  len(animals),
		CatalogCount: len(catalog),
		UserCount:    len(users),
	}

	renderTemplate(w, "home", data)
}

func AdminHandler(w http.ResponseWriter, r *http.Request) {
	adminCount := 0
	for _, user := range users {
		if user["role"] == "admin" {
			adminCount++
		}
	}

	data := TemplateData{
		Title:      "Админ панель",
		Year:       time.Now().Year(),
		Port:       "8081",
		UserCount:  len(users),
		AdminCount: adminCount,
		Users:      users,
	}

	renderTemplate(w, "admin", data)
}

func APIHandler(w http.ResponseWriter, r *http.Request) {
	data := TemplateData{
		Title: "Документация API",
		Year:  time.Now().Year(),
		Port:  "8081",
	}

	renderTemplate(w, "api", data)
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "ok",
		"time":    time.Now().Format(time.RFC3339),
		"version": "1.0.0",
	})
}

// ============ API ОБРАБОТЧИКИ ============

func AnimalsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(animals)
}

func AnimalHandler(w http.ResponseWriter, r *http.Request) {
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

func CatalogHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(catalog)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// Поиск пользователя
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

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// Проверка существования пользователя
	for _, user := range users {
		if user["username"] == data.Username {
			http.Error(w, "Пользователь уже существует", http.StatusConflict)
			return
		}
	}

	// Создание нового пользователя
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

func CreateAnimalHandler(w http.ResponseWriter, r *http.Request) {
	// Проверка прав администратора
	if r.Header.Get("X-Admin") != "true" {
		http.Error(w, "Требуются права администратора", http.StatusForbidden)
		return
	}

	var animal Animal
	if err := json.NewDecoder(r.Body).Decode(&animal); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// Валидация
	if animal.Name == "" || animal.Species == "" {
		http.Error(w, "Имя и вид обязательны", http.StatusBadRequest)
		return
	}

	// Создание нового животного
	animal.ID = len(animals) + 1
	animals = append(animals, animal)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(animal)
}

func CreateCatalogItemHandler(w http.ResponseWriter, r *http.Request) {
	// Проверка прав администратора
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

	// Создание нового товара
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
