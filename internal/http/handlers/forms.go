package handlers

import (
	"net/http"
	"strconv"

	"pet-shelter/internal/models"
)

func (h *Handler) AdminCreateAnimal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}

	age, _ := strconv.Atoi(r.FormValue("age"))
	animal := models.Animal{
		Name:         r.FormValue("name"),
		Species:      r.FormValue("species"),
		Breed:        r.FormValue("breed"),
		Age:          age,
		Gender:       r.FormValue("gender"),
		HealthStatus: r.FormValue("health_status"),
		Description:  r.FormValue("description"),
	}

	if err := h.animals.Create(r.Context(), &animal); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (h *Handler) AdminDeleteAnimal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}

	if err := h.animals.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (h *Handler) AdminCreateCatalogItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}

	price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
	quantity, _ := strconv.Atoi(r.FormValue("quantity"))

	item := models.CatalogItem{
		Name:        r.FormValue("name"),
		Type:        r.FormValue("type"),
		Category:    r.FormValue("category"),
		Price:       price,
		Quantity:    quantity,
		Description: r.FormValue("description"),
		Supplier:    r.FormValue("supplier"),
	}

	if err := h.catalog.Create(r.Context(), &item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (h *Handler) AdminDeleteCatalogItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}

	if err := h.catalog.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
