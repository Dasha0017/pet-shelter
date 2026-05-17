package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"pet-shelter/internal/dto"
)

func (h *Handler) LoginForm(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}

	role := r.FormValue("role")
	username := r.FormValue("username")
	password := r.FormValue("password")

	token, _, err := h.auth.Login(r.Context(), username, password, role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	h.setAccessCookie(w, token)
	if role == "admin" {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) LoginAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	token, user, err := h.auth.Login(r.Context(), req.Username, req.Password, req.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	h.setAccessCookie(w, token)
	writeJSON(w, http.StatusOK, dto.LoginResponse{Status: "ok", Token: token, User: user})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   h.cfg.CookieSecure(),
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) DevAutoLogin(w http.ResponseWriter, r *http.Request) {
	if !h.cfg.IsDev() {
		http.NotFound(w, r)
		return
	}

	token, err := h.auth.DevAdminToken()
	if err != nil {
		http.Error(w, "token error", http.StatusInternalServerError)
		return
	}

	h.setAccessCookie(w, token)
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (h *Handler) setAccessCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   h.cfg.CookieSecure(),
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})
}
