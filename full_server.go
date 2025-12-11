package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// Модели данных
type Animal struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Species     string    `json:"species"`
	Breed       string    `json:"breed,omitempty"`
	Age         int       `json:"age"`
	Adopted     bool      `json:"adopted"`
	ArrivalDate time.Time `json:"arrival_date"`
}

type CatalogItem struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	Role     string `json:"role"`
}

// База данных в памяти
var animals = []Animal{
	{ID: 1, Name: "Барсик", Species: "Кошка", Breed: "Персидская", Age: 3, Adopted: false, ArrivalDate: time.Now().AddDate(0, -2, 0)},
	{ID: 2, Name: "Шарик", Species: "Собака", Breed: "Лабрадор", Age: 2, Adopted: false, ArrivalDate: time.Now().AddDate(0, -1, -15)},
	{ID: 3, Name: "Мурка", Species: "Кошка", Breed: "Дворовая", Age: 1, Adopted: true, ArrivalDate: time.Now().AddDate(0, -3, 0)},
}

var catalog = []CatalogItem{
	{ID: 1, Name: "Сухой корм для собак", Type: "Корм", Price: 25.99, Quantity: 50},
	{ID: 2, Name: "Игрушка для кошек", Type: "Игрушка", Price: 5.99, Quantity: 100},
	{ID: 3, Name: "Наполнитель для туалета", Type: "Аксессуар", Price: 15.50, Quantity: 30},
}

var users = []User{
	{ID: 1, Username: "admin", Password: "admin123", Role: "admin"},
	{ID: 2, Username: "user", Password: "user123", Role: "user"},
}

var nextAnimalID = 4
var nextCatalogID = 4
var nextUserID = 3

func main() {
	r := mux.NewRouter()

	// Middleware для CORS
	r.Use(corsMiddleware)

	// Веб-интерфейс
	r.HandleFunc("/", homePage).Methods("GET")
	r.HandleFunc("/admin", adminPage).Methods("GET")

	// API: Аутентификация
	r.HandleFunc("/api/login", login).Methods("POST")
	r.HandleFunc("/api/register", register).Methods("POST")

	// API: Животные (публичные)
	r.HandleFunc("/api/animals", getAnimals).Methods("GET")
	r.HandleFunc("/api/animals/{id}", getAnimal).Methods("GET")

	// API: Животные (админ)
	r.HandleFunc("/api/admin/animals", createAnimal).Methods("POST")
	r.HandleFunc("/api/admin/animals/{id}", updateAnimal).Methods("PUT")
	r.HandleFunc("/api/admin/animals/{id}", deleteAnimal).Methods("DELETE")

	// API: Каталог (публичный)
	r.HandleFunc("/api/catalog", getCatalog).Methods("GET")
	r.HandleFunc("/api/catalog/{id}", getCatalogItem).Methods("GET")

	// API: Каталог (админ)
	r.HandleFunc("/api/admin/catalog", createCatalogItem).Methods("POST")
	r.HandleFunc("/api/admin/catalog/{id}", updateCatalogItem).Methods("PUT")
	r.HandleFunc("/api/admin/catalog/{id}", deleteCatalogItem).Methods("DELETE")

	// Запуск сервера
	port := "8080"
	log.Printf("🚀 Pet Shelter System запущен!")
	log.Printf("🌐 Веб-интерфейс: http://localhost:%s", port)
	log.Printf("👑 Админ панель: http://localhost:%s/admin", port)
	log.Printf("📡 API доступен по адресу: http://localhost:%s/api/*", port)
	
	log.Fatal(http.ListenAndServe(":"+port, r))
}

// Middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// Вспомогательные функции
func isAdmin(r *http.Request) bool {
	return r.Header.Get("Authorization") == "Bearer admin-token" || 
	       r.Header.Get("X-Admin") == "true" ||
	       strings.Contains(r.URL.Query().Get("token"), "admin")
}

func sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func sendError(w http.ResponseWriter, status int, message string) {
	sendJSON(w, status, map[string]string{"error": message})
}

// Веб-страницы
func homePage(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Pet Shelter - Главная</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: Arial, sans-serif; background: #f0f2f5; }
        .container { max-width: 1200px; margin: 0 auto; padding: 20px; }
        header { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 40px 0; text-align: center; border-radius: 0 0 20px 20px; }
        h1 { font-size: 2.5em; margin-bottom: 10px; }
        .stats { display: grid; grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); gap: 20px; margin: 40px 0; }
        .stat-card { background: white; padding: 20px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); text-align: center; }
        .stat-number { font-size: 2em; color: #667eea; font-weight: bold; }
        .actions { display: flex; gap: 10px; margin: 20px 0; flex-wrap: wrap; }
        .btn { padding: 10px 20px; border: none; border-radius: 5px; cursor: pointer; text-decoration: none; display: inline-block; }
        .btn-primary { background: #667eea; color: white; }
        .btn-secondary { background: #764ba2; color: white; }
        .btn-success { background: #48bb78; color: white; }
        .api-section { background: white; padding: 20px; border-radius: 10px; margin: 20px 0; }
        .endpoint { background: #f7fafc; padding: 10px; margin: 10px 0; border-left: 4px solid #667eea; }
        .method { display: inline-block; padding: 3px 8px; border-radius: 3px; font-weight: bold; margin-right: 10px; }
        .get { background: #48bb78; color: white; }
        .post { background: #4299e1; color: white; }
        .put { background: #ed8936; color: white; }
        .delete { background: #f56565; color: white; }
        .test-btn { background: #e2e8f0; border: none; padding: 5px 10px; border-radius: 3px; cursor: pointer; margin-left: 10px; }
    </style>
</head>
<body>
    <header>
        <div class="container">
            <h1>🐾 Pet Shelter Management System</h1>
            <p>Система управления приютом для животных</p>
        </div>
    </header>
    
    <div class="container">
        <div class="stats">
            <div class="stat-card">
                <div class="stat-number">` + fmt.Sprint(len(animals)) + `</div>
                <div>Животных в приюте</div>
            </div>
            <div class="stat-card">
                <div class="stat-number">` + fmt.Sprint(len(catalog)) + `</div>
                <div>Товаров в каталоге</div>
            </div>
            <div class="stat-card">
                <div class="stat-number">` + fmt.Sprint(len(users)) + `</div>
                <div>Пользователей</div>
            </div>
        </div>
        
        <div class="actions">
            <a href="/admin" class="btn btn-primary">👑 Панель администратора</a>
            <button class="btn btn-secondary" onclick="loadAnimals()">🦮 Показать животных</button>
            <button class="btn btn-success" onclick="loadCatalog()">🛒 Показать каталог</button>
        </div>
        
        <div id="results"></div>
        
        <div class="api-section">
            <h2>📡 REST API Endpoints</h2>
            
            <div class="endpoint">
                <span class="method get">GET</span> <code>/api/animals</code>
                <span class="description">Получить список всех животных</span>
                <button class="test-btn" onclick="testEndpoint('/api/animals', 'GET')">Тест</button>
            </div>
            
            <div class="endpoint">
                <span class="method post">POST</span> <code>/api/admin/animals</code>
                <span class="description">Добавить новое животное</span>
                <button class="test-btn" onclick="testCreateAnimal()">Тест</button>
            </div>
            
            <div class="endpoint">
                <span class="method get">GET</span> <code>/api/catalog</code>
                <span class="description">Получить каталог товаров</span>
                <button class="test-btn" onclick="testEndpoint('/api/catalog', 'GET')">Тест</button>
            </div>
            
            <div class="endpoint">
                <span class="method post">POST</span> <code>/api/login</code>
                <span class="description">Войти в систему</span>
                <button class="test-btn" onclick="testLogin()">Тест</button>
            </div>
        </div>
    </div>
    
    <script>
        function loadAnimals() {
            fetch('/api/animals')
                .then(r => r.json())
                .then(data => {
                    document.getElementById('results').innerHTML = 
                        '<h3>Животные:</h3><pre>' + JSON.stringify(data, null, 2) + '</pre>';
                });
        }
        
        function loadCatalog() {
            fetch('/api/catalog')
                .then(r => r.json())
                .then(data => {
                    document.getElementById('results').innerHTML = 
                        '<h3>Каталог:</h3><pre>' + JSON.stringify(data, null, 2) + '</pre>';
                });
        }
        
        function testEndpoint(url, method) {
            fetch(url, { method: method })
                .then(r => r.json())
                .then(data => alert(JSON.stringify(data, null, 2)))
                .catch(err => alert('Ошибка: ' + err));
        }
        
        function testCreateAnimal() {
            const animal = {
                name: "Новый питомец " + new Date().getSeconds(),
                species: "Кошка",
                breed: "Дворовая",
                age: 1
            };
            
            fetch('/api/admin/animals', {
                method: 'POST',
                headers: { 
                    'Content-Type': 'application/json',
                    'X-Admin': 'true'
                },
                body: JSON.stringify(animal)
            })
            .then(r => r.json())
            .then(data => alert('Животное добавлено: ' + JSON.stringify(data, null, 2)));
        }
        
        function testLogin() {
            fetch('/api/login', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ username: 'admin', password: 'admin123' })
            })
            .then(r => r.json())
            .then(data => alert('Токен: ' + data.token));
        }
        
        // Автозагрузка животных при открытии
        window.onload = loadAnimals;
    </script>
</body>
</html>`
	
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func adminPage(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>Панель администратора</title>
    <style>
        body { font-family: Arial; padding: 20px; }
        .tab { overflow: hidden; border: 1px solid #ccc; background: #f1f1f1; }
        .tab button { background: inherit; float: left; border: none; outline: none; cursor: pointer; padding: 14px 16px; transition: 0.3s; }
        .tab button:hover { background: #ddd; }
        .tab button.active { background: #667eea; color: white; }
        .tabcontent { display: none; padding: 20px; border: 1px solid #ccc; border-top: none; }
        .form-group { margin: 10px 0; }
        label { display: block; margin-bottom: 5px; }
        input, textarea { width: 100%; padding: 8px; }
        button { padding: 10px 20px; background: #667eea; color: white; border: none; cursor: pointer; }
    </style>
</head>
<body>
    <h1>👑 Панель администратора Pet Shelter</h1>
    
    <div class="tab">
        <button class="tablinks active" onclick="openTab(event, 'animals')">Животные</button>
        <button class="tablinks" onclick="openTab(event, 'catalog')">Каталог</button>
        <button class="tablinks" onclick="openTab(event, 'users')">Пользователи</button>
    </div>
    
    <div id="animals" class="tabcontent" style="display: block;">
        <h2>Управление животными</h2>
        <button onclick="loadAnimals()">Обновить список</button>
        <div id="animalsList"></div>
        
        <h3>Добавить новое животное</h3>
        <div class="form-group">
            <label>Имя:</label>
            <input type="text" id="animalName" placeholder="Барсик">
        </div>
        <div class="form-group">
            <label>Вид:</label>
            <input type="text" id="animalSpecies" placeholder="Кошка">
        </div>
        <div class="form-group">
            <label>Порода:</label>
            <input type="text" id="animalBreed" placeholder="Персидская">
        </div>
        <div class="form-group">
            <label>Возраст:</label>
            <input type="number" id="animalAge" placeholder="2">
        </div>
        <button onclick="createAnimal()">Добавить животное</button>
    </div>
    
    <div id="catalog" class="tabcontent">
        <h2>Управление каталогом</h2>
        <div id="catalogList"></div>
    </div>
    
    <div id="users" class="tabcontent">
        <h2>Управление пользователями</h2>
        <p>Всего пользователей: ` + fmt.Sprint(len(users)) + `</p>
    </div>
    
    <script>
        function openTab(evt, tabName) {
            var i, tabcontent, tablinks;
            tabcontent = document.getElementsByClassName("tabcontent");
            for (i = 0; i < tabcontent.length; i++) {
                tabcontent[i].style.display = "none";
            }
            tablinks = document.getElementsByClassName("tablinks");
            for (i = 0; i < tablinks.length; i++) {
                tablinks[i].className = tablinks[i].className.replace(" active", "");
            }
            document.getElementById(tabName).style.display = "block";
            evt.currentTarget.className += " active";
        }
        
        function loadAnimals() {
            fetch('/api/animals')
                .then(r => r.json())
                .then(data => {
                    let html = '<table border="1" cellpadding="10"><tr><th>ID</th><th>Имя</th><th>Вид</th><th>Порода</th><th>Возраст</th><th>Статус</th></tr>';
                    data.forEach(animal => {
                        html += `<tr>
                            <td>${animal.id}</td>
                            <td>${animal.name}</td>
                            <td>${animal.species}</td>
                            <td>${animal.breed || '-'}</td>
                            <td>${animal.age}</td>
                            <td>${animal.adopted ? 'Пристроен' : 'В приюте'}</td>
                        </tr>`;
                    });
                    html += '</table>';
                    document.getElementById('animalsList').innerHTML = html;
                });
        }
        
        function createAnimal() {
            const animal = {
                name: document.getElementById('animalName').value,
                species: document.getElementById('animalSpecies').value,
                breed: document.getElementById('animalBreed').value,
                age: parseInt(document.getElementById('animalAge').value)
            };
            
            fetch('/api/admin/animals', {
                method: 'POST',
                headers: { 
                    'Content-Type': 'application/json',
                    'X-Admin': 'true'
                },
                body: JSON.stringify(animal)
            })
            .then(r => r.json())
            .then(data => {
                alert('Животное добавлено! ID: ' + data.id);
                loadAnimals();
            });
        }
        
        // Загружаем животных при открытии
        window.onload = loadAnimals;
    </script>
</body>
</html>`
	
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

// API функции
func login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		sendError(w, http.StatusBadRequest, "Неверный формат данных")
		return
	}
	
	for _, user := range users {
		if user.Username == creds.Username && user.Password == creds.Password {
			token := fmt.Sprintf("%s-token-%d", user.Role, user.ID)
			sendJSON(w, http.StatusOK, map[string]interface{}{
				"token": token,
				"user": map[string]interface{}{
					"id":       user.ID,
					"username": user.Username,
					"role":     user.Role,
				},
			})
			return
		}
	}
	
	sendError(w, http.StatusUnauthorized, "Неверные учетные данные")
}

func register(w http.ResponseWriter, r *http.Request) {
	var newUser struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		sendError(w, http.StatusBadRequest, "Неверный формат данных")
		return
	}
	
	for _, user := range users {
		if user.Username == newUser.Username {
			sendError(w, http.StatusConflict, "Пользователь уже существует")
			return
		}
	}
	
	nextUserID++
	user := User{
		ID:       nextUserID,
		Username: newUser.Username,
		Password: newUser.Password,
		Role:     "user",
	}
	
	users = append(users, user)
	
	sendJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Пользователь создан",
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}

func getAnimals(w http.ResponseWriter, r *http.Request) {
	sendJSON(w, http.StatusOK, animals)
}

func getAnimal(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, "Неверный ID")
		return
	}
	
	for _, animal := range animals {
		if animal.ID == id {
			sendJSON(w, http.StatusOK, animal)
			return
		}
	}
	
	sendError(w, http.StatusNotFound, "Животное не найдено")
}

func createAnimal(w http.ResponseWriter, r *http.Request) {
	if !isAdmin(r) {
		sendError(w, http.StatusForbidden, "Требуются права администратора")
		return
	}
	
	var animal Animal
	if err := json.NewDecoder(r.Body).Decode(&animal); err != nil {
		sendError(w, http.StatusBadRequest, "Неверный формат данных")
		return
	}
	
	animal.ID = nextAnimalID
	nextAnimalID++
	animal.ArrivalDate = time.Now()
	animal.Adopted = false
	
	animals = append(animals, animal)
	
	sendJSON(w, http.StatusCreated, animal)
}

func updateAnimal(w http.ResponseWriter, r *http.Request) {
	if !isAdmin(r) {
		sendError(w, http.StatusForbidden, "Требуются права администратора")
		return
	}
	
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, "Неверный ID")
		return
	}
	
	var updates Animal
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		sendError(w, http.StatusBadRequest, "Неверный формат данных")
		return
	}
	
	for i, animal := range animals {
		if animal.ID == id {
			updates.ID = id
			updates.ArrivalDate = animal.ArrivalDate
			animals[i] = updates
			sendJSON(w, http.StatusOK, updates)
			return
		}
	}
	
	sendError(w, http.StatusNotFound, "Животное не найдено")
}

func deleteAnimal(w http.ResponseWriter, r *http.Request) {
	if !isAdmin(r) {
		sendError(w, http.StatusForbidden, "Требуются права администратора")
		return
	}
	
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, "Неверный ID")
		return
	}
	
	for i, animal := range animals {
		if animal.ID == id {
			animals = append(animals[:i], animals[i+1:]...)
			sendJSON(w, http.StatusOK, map[string]string{"message": "Животное удалено"})
			return
		}
	}
	
	sendError(w, http.StatusNotFound, "Животное не найдено")
}

func getCatalog(w http.ResponseWriter, r *http.Request) {
	sendJSON(w, http.StatusOK, catalog)
}

func getCatalogItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, "Неверный ID")
		return
	}
	
	for _, item := range catalog {
		if item.ID == id {
			sendJSON(w, http.StatusOK, item)
			return
		}
	}
	
	sendError(w, http.StatusNotFound, "Товар не найден")
}

func createCatalogItem(w http.ResponseWriter, r *http.Request) {
	if !isAdmin(r) {
		sendError(w, http.StatusForbidden, "Требуются права администратора")
		return
	}
	
	var item CatalogItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		sendError(w, http.StatusBadRequest, "Неверный формат данных")
		return
	}
	
	item.ID = nextCatalogID
	nextCatalogID++
	
	catalog = append(catalog, item)
	sendJSON(w, http.StatusCreated, item)
}

func updateCatalogItem(w http.ResponseWriter, r *http.Request) {
	if !isAdmin(r) {
		sendError(w, http.StatusForbidden, "Требуются права администратора")
		return
	}
	
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, "Неверный ID")
		return
	}
	
	var updates CatalogItem
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		sendError(w, http.StatusBadRequest, "Неверный формат данных")
		return
	}
	
	for i, item := range catalog {
		if item.ID == id {
			updates.ID = id
			catalog[i] = updates
			sendJSON(w, http.StatusOK, updates)
			return
		}
	}
	
	sendError(w, http.StatusNotFound, "Товар не найден")
}

func deleteCatalogItem(w http.ResponseWriter, r *http.Request) {
	if !isAdmin(r) {
		sendError(w, http.StatusForbidden, "Требуются права администратора")
		return
	}
	
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, "Неверный ID")
		return
	}
	
	for i, item := range catalog {
		if item.ID == id {
			catalog = append(catalog[:i], catalog[i+1:]...)
			sendJSON(w, http.StatusOK, map[string]string{"message": "Товар удален"})
			return
		}
	}
	
	sendError(w, http.StatusNotFound, "Товар не найден")
}