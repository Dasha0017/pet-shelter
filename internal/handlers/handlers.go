package handlers

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var templates *template.Template

// InitTemplates загружает все шаблоны из web/templates
func InitTemplates() error {
	templatesDir := filepath.Join("web", "templates")

	// проверим, что директория существует (не обязательно, но полезно)
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		log.Printf("директория шаблонов не найдена: %s", templatesDir)
		return err
	}

	pattern := filepath.Join(templatesDir, "*.html")
	t, err := template.ParseGlob(pattern)
	if err != nil {
		log.Printf("ошибка загрузки шаблонов: %v", err)
		return err
	}

	templates = t
	log.Printf("загружено шаблонов: %d", len(t.Templates()))
	return nil
}

// ---- Home ----

type HomeData struct {
	AnimalCount  int
	CatalogCount int
	UserCount    int
	Year         int
	Port         string
}

func Home(w http.ResponseWriter, r *http.Request, data HomeData) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := templates.ExecuteTemplate(w, "home.html", data); err != nil {
		log.Println("Home template error:", err)
		http.Error(w, "server error", http.StatusInternalServerError)
	}
}

// ---- Admin ----

type UserView struct {
	ID       int
	Username string
	Role     string
}

type AdminData struct {
	UserCount  int
	AdminCount int
	Users      []UserView
}

func Admin(w http.ResponseWriter, r *http.Request, rawUsers []map[string]interface{}) {
	views := make([]UserView, 0, len(rawUsers))
	adminCount := 0

	for _, u := range rawUsers {
		id, _ := u["id"].(int)
		username, _ := u["username"].(string)
		role, _ := u["role"].(string)
		views = append(views, UserView{
			ID:       id,
			Username: username,
			Role:     role,
		})
		if role == "admin" {
			adminCount++
		}
	}

	data := AdminData{
		UserCount:  len(views),
		AdminCount: adminCount,
		Users:      views,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := templates.ExecuteTemplate(w, "admin.html", data); err != nil {
		log.Println("Admin template error:", err)
		http.Error(w, "server error", http.StatusInternalServerError)
	}
}

// ---- Animals page ----

type Animal struct {
	ID      int
	Name    string
	Species string
	Age     int
}

type AnimalsPageData struct {
	Animals []Animal
}

func AnimalsPage(w http.ResponseWriter, r *http.Request, animals []Animal) {
	data := AnimalsPageData{
		Animals: animals,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := templates.ExecuteTemplate(w, "animals.html", data); err != nil {
		log.Println("Animals template error:", err)
		http.Error(w, "server error", http.StatusInternalServerError)
	}
}

// ---- Catalog page ----

type CatalogItemView struct {
	ID       int
	Name     string
	Price    float64
	Quantity int
}

type CatalogPageData struct {
	Items []CatalogItemView
}

func CatalogPage(w http.ResponseWriter, r *http.Request, catalogRaw []map[string]interface{}) {
	data := CatalogPageData{}

	for _, it := range catalogRaw {
		id, _ := it["id"].(int)
		name, _ := it["name"].(string)
		price, _ := it["price"].(float64)
		qty, _ := it["quantity"].(int)

		data.Items = append(data.Items, CatalogItemView{
			ID:       id,
			Name:     name,
			Price:    price,
			Quantity: qty,
		})
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := templates.ExecuteTemplate(w, "catalog.html", data); err != nil {
		log.Println("Catalog template error:", err)
		http.Error(w, "server error", http.StatusInternalServerError)
	}
}

// ---- API docs ----

type APIDocsData struct {
	Port string
}

func APIDocs(w http.ResponseWriter, r *http.Request, port string) {
	data := APIDocsData{Port: port}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := templates.ExecuteTemplate(w, "api.html", data); err != nil {
		log.Println("API docs template error:", err)
		http.Error(w, "server error", http.StatusInternalServerError)
	}
}
