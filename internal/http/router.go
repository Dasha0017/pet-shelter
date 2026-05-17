package httpserver

import (
	"net/http"

	"pet-shelter/internal/config"
	"pet-shelter/internal/http/handlers"
	"pet-shelter/internal/http/middleware"
)

func NewRouter(h *handlers.Handler, auth *middleware.AuthMiddleware, cfg config.Config) http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	mux.HandleFunc("/", h.HomePage)
	mux.HandleFunc("/animals", h.AnimalsPage)
	mux.HandleFunc("/catalog", h.CatalogPage)
	mux.HandleFunc("/api/docs", h.APIDocs)
	mux.HandleFunc("/health", h.Health)

	mux.HandleFunc("/login", h.LoginPage)
	mux.HandleFunc("/api/login", h.LoginAPI)
	mux.HandleFunc("/logout", h.Logout)
	if cfg.IsDev() {
		mux.HandleFunc("/dev-auto-login", h.DevAutoLogin)
	}

	mux.HandleFunc("/api/animals", h.AnimalsAPI)
	mux.HandleFunc("/api/catalog", h.CatalogAPI)

	adminOnly := auth.RequireRole("admin")
	mux.Handle("/admin", adminOnly(http.HandlerFunc(h.AdminPage)))
	mux.Handle("/admin/animals/create", adminOnly(http.HandlerFunc(h.AdminCreateAnimal)))
	mux.Handle("/admin/animals/delete", adminOnly(http.HandlerFunc(h.AdminDeleteAnimal)))
	mux.Handle("/admin/catalog/create", adminOnly(http.HandlerFunc(h.AdminCreateCatalogItem)))
	mux.Handle("/admin/catalog/delete", adminOnly(http.HandlerFunc(h.AdminDeleteCatalogItem)))

	return mux
}
