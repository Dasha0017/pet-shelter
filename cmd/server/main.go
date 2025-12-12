package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"pet-shelter/internal/config"
	"pet-shelter/internal/handlers"
	"pet-shelter/internal/repository"
)

var db *sql.DB

// alias, чтобы переиспользовать структуру из handlers
type Animal = handlers.Animal

func getLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "127.0.0.1"
	}
	defer conn.Close()
	if addr, ok := conn.LocalAddr().(*net.UDPAddr); ok {
		return addr.IP.String()
	}
	return "127.0.0.1"
}

func main() {
	// Загружаем конфиг и подключаемся к БД
	cfg := config.Load()

	var err error
	db, err = repository.NewDB(cfg)
	if err != nil {
		log.Fatal("DB connection error:", err)
	}
	defer db.Close()

	if err := handlers.InitTemplates(); err != nil {
		log.Fatal("InitTemplates error:", err)
	}
	log.Println("Templates initialized")

	// Статика (CSS и т.п.)
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// HTML-страницы
	http.HandleFunc("/", homePageHandler)
	http.HandleFunc("/animals", animalsPageHandler)
	http.HandleFunc("/catalog", catalogPageHandler)
	http.HandleFunc("/admin", adminPageHandler)
	http.HandleFunc("/api/docs", apiDocsHandler)
	http.HandleFunc("/health", healthHandler)

	// Публичный API (JSON)
	http.HandleFunc("/api/animals", animalsAPIHandler)
	http.HandleFunc("/api/catalog", catalogAPIHandler)

	// Админские формы
	http.HandleFunc("/admin/animals/create", adminCreateAnimalFormHandler)
	http.HandleFunc("/admin/animals/delete", adminDeleteAnimalFormHandler)
	http.HandleFunc("/admin/catalog/create", adminCreateCatalogFormHandler)
	http.HandleFunc("/admin/catalog/delete", adminDeleteCatalogFormHandler)

	port := cfg.ServerPort
	if port == "" {
		port = "8081"
	}

	host := getLocalIP()
	addr := ":" + port

	log.Println("======================================")
	log.Printf("Pet Shelter запущен на порту %s", port)
	log.Printf("Локально:   http://localhost:%s", port)
	log.Printf("В сети LAN: http://%s:%s", host, port)
	log.Println("Нажми Ctrl+C для остановки сервера")
	log.Println("======================================")

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}

// ---------------- HTML handlers ----------------

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	animals, err := repository.GetAllAnimals(db)
	if err != nil {
		log.Println("homePageHandler GetAllAnimals error:", err)
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	catalogItems, err := repository.GetAllCatalogItems(db)
	if err != nil {
		log.Println("homePageHandler GetAllCatalogItems error:", err)
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	data := handlers.HomeData{
		AnimalCount:  len(animals),
		CatalogCount: len(catalogItems),
		UserCount:    2, // пока захардкожено
		Year:         time.Now().Year(),
		Port:         portFromRequest(r),
	}

	handlers.Home(w, r, data)
}

func animalsPageHandler(w http.ResponseWriter, r *http.Request) {
	animals, err := repository.GetAllAnimals(db)
	if err != nil {
		log.Println("animalsPageHandler error:", err)
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	handlers.AnimalsPage(w, r, animals)
}

func catalogPageHandler(w http.ResponseWriter, r *http.Request) {
	items, err := repository.GetAllCatalogItems(db)
	if err != nil {
		log.Println("catalogPageHandler error:", err)
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	raw := make([]map[string]interface{}, 0, len(items))
	for _, it := range items {
		raw = append(raw, map[string]interface{}{
			"id":       it.ID,
			"name":     it.Name,
			"price":    it.Price,
			"quantity": it.Quantity,
		})
	}
	handlers.CatalogPage(w, r, raw)
}

func adminPageHandler(w http.ResponseWriter, r *http.Request) {
	rawUsers := []map[string]interface{}{
		{"id": 1, "username": "admin", "role": "admin"},
		{"id": 2, "username": "user", "role": "user"},
	}
	handlers.Admin(w, r, rawUsers)
}

func apiDocsHandler(w http.ResponseWriter, r *http.Request) {
	handlers.APIDocs(w, r, portFromRequest(r))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "ok",
		"time":    time.Now().Format(time.RFC3339),
		"version": "1.0.0",
	})
}

// ---------------- JSON API handlers ----------------

func animalsAPIHandler(w http.ResponseWriter, r *http.Request) {
	animals, err := repository.GetAllAnimals(db)
	if err != nil {
		log.Println("GetAllAnimals error:", err)
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(animals)
}

func catalogAPIHandler(w http.ResponseWriter, r *http.Request) {
	items, err := repository.GetAllCatalogItems(db)
	if err != nil {
		log.Println("GetAllCatalogItems error:", err)
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// ---------------- Admin forms (HTML) ----------------

func adminCreateAnimalFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Только POST", http.StatusMethodNotAllowed)
		return
	}

	name := r.FormValue("name")
	species := r.FormValue("species")
	ageStr := r.FormValue("age")

	if name == "" || species == "" {
		http.Error(w, "Имя и вид обязательны", http.StatusBadRequest)
		return
	}

	age, _ := strconv.Atoi(ageStr)

	a := handlers.Animal{
		Name:    name,
		Species: species,
		Age:     age,
	}

	if err := repository.CreateAnimal(db, &a); err != nil {
		log.Println("CreateAnimal error:", err)
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/animals", http.StatusSeeOther)
}

func adminDeleteAnimalFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Только POST", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	if err := repository.DeleteAnimal(db, id); err != nil {
		log.Println("DeleteAnimal error:", err)
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/animals", http.StatusSeeOther)
}

func adminCreateCatalogFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Только POST", http.StatusMethodNotAllowed)
		return
	}

	name := r.FormValue("name")
	priceStr := r.FormValue("price")
	qtyStr := r.FormValue("quantity")

	if name == "" {
		http.Error(w, "Название обязательно", http.StatusBadRequest)
		return
	}

	price, _ := strconv.ParseFloat(priceStr, 64)
	qty, _ := strconv.Atoi(qtyStr)

	it := repository.CatalogItem{
		Name:     name,
		Type:     "manual",
		Category: "",
		Price:    price,
		Quantity: qty,
	}

	if err := repository.CreateCatalogItem(db, &it); err != nil {
		log.Println("CreateCatalogItem error:", err)
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/catalog", http.StatusSeeOther)
}

func adminDeleteCatalogFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Только POST", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	if err := repository.DeleteCatalogItem(db, id); err != nil {
		log.Println("DeleteCatalogItem error:", err)
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/catalog", http.StatusSeeOther)
}

// вспомогательный хелпер, чтобы из r.Host вытащить порт при проксировании/других режимах
func portFromRequest(r *http.Request) string {
	host := r.Host
	// ожидаемый формат host:port
	for i := len(host) - 1; i >= 0; i-- {
		if host[i] == ':' {
			return host[i+1:]
		}
	}
	return ""
}
