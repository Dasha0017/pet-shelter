package handlers

import "net/http"

func (h *Handler) AnimalsAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	animals, err := h.animals.List(r.Context())
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, animals)
}

func (h *Handler) CatalogAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	items, err := h.catalog.List(r.Context())
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, items)
}
