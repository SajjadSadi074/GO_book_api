package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

var TokenAuth = jwtauth.New("HS256", []byte("your-secret-key"), nil)

func Login(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")

	if username != "admin" || password != "pass123" {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	_, token, _ := TokenAuth.Encode(map[string]interface{}{"user": username})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
