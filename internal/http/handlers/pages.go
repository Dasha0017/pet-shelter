package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"pet-shelter/internal/models"
)

type HomeData struct {
	AnimalCount  int
	CatalogCount int
	UserCount    int
	Year         int
	Port         string
}

type AdminData struct {
	Users   []models.User
	Animals []models.Animal
	Catalog []models.CatalogItem
}

type LoginPageData struct {
	Role string
}

type APIDocsData struct {
	Port string
}

func (h *Handler) HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	ctx := r.Context()
	animalCount, err := h.animals.Count(ctx)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	catalogCount, err := h.catalog.Count(ctx)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	userCount, err := h.auth.UserCount(ctx)
	if err != nil {
		userCount = 0
	}

	data := HomeData{
		AnimalCount:  animalCount,
		CatalogCount: catalogCount,
		UserCount:    userCount,
		Year:         time.Now().Year(),
		Port:         portFromRequest(r),
	}

	if err := h.renderer.Render(w, "home", data); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

func (h *Handler) AnimalsPage(w http.ResponseWriter, r *http.Request) {
	animals, err := h.animals.List(r.Context())
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	if err := h.renderer.Render(w, "animals", animals); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

func (h *Handler) CatalogPage(w http.ResponseWriter, r *http.Request) {
	items, err := h.catalog.List(r.Context())
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	if err := h.renderer.Render(w, "catalog", items); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

func (h *Handler) AdminPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := h.auth.ListUsers(ctx)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	animals, err := h.animals.List(ctx)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	catalog, err := h.catalog.List(ctx)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	data := AdminData{Users: users, Animals: animals, Catalog: catalog}
	if err := h.renderer.Render(w, "admin", data); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.LoginForm(w, r)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	role := r.URL.Query().Get("role")
	if role == "" {
		role = "user"
	}

	if err := h.renderer.Render(w, "login", LoginPageData{Role: role}); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

func (h *Handler) APIDocs(w http.ResponseWriter, r *http.Request) {
	if err := h.renderer.Render(w, "api_docs", APIDocsData{Port: portFromRequest(r)}); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status":  "ok",
		"time":    time.Now().Format(time.RFC3339),
		"version": "refactored-1.0.0",
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func portFromRequest(r *http.Request) string {
	host := r.Host
	for i := len(host) - 1; i >= 0; i-- {
		if host[i] == ':' {
			return host[i+1:]
		}
	}
	return ""
}
