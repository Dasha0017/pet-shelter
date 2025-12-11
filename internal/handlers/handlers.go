package handlers

import (
	"encoding/json"
	"net/http"
	"pet-shelter/internal/models"
	"pet-shelter/internal/service"
	"strconv"

	"github.com/gorilla/mux"
)

type AnimalHandler struct {
	animalService *service.AnimalService
}

func NewAnimalHandler(animalService *service.AnimalService) *AnimalHandler {
	return &AnimalHandler{animalService: animalService}
}

func (h *AnimalHandler) CreateAnimal(w http.ResponseWriter, r *http.Request) {
	var animal models.Animal
	if err := json.NewDecoder(r.Body).Decode(&animal); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.animalService.Create(&animal); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(animal)
}

func (h *AnimalHandler) GetAnimals(w http.ResponseWriter, r *http.Request) {
	animals, err := h.animalService.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(animals)
}

func (h *AnimalHandler) GetAnimal(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid animal ID", http.StatusBadRequest)
		return
	}

	animal, err := h.animalService.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if animal == nil {
		http.Error(w, "Animal not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(animal)
}

func (h *AnimalHandler) UpdateAnimal(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid animal ID", http.StatusBadRequest)
		return
	}

	var animal models.Animal
	if err := json.NewDecoder(r.Body).Decode(&animal); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	animal.ID = id

	if err := h.animalService.Update(&animal); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(animal)
}

func (h *AnimalHandler) DeleteAnimal(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid animal ID", http.StatusBadRequest)
		return
	}

	if err := h.animalService.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
