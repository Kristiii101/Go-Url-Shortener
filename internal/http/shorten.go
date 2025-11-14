package http

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

// request body structure
type shortenRequest struct {
	URL string `json:"url"`
}

// response structure
type shortenResponse struct {
	ShortURL string `json:"short_url"`
}

// simple random code generator (placeholder)
func generateShortCode() string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())

	code := make([]byte, 6) // 6 characters
	for i := range code {
		code[i] = letters[rand.Intn(len(letters))]
	}
	return string(code)
}

// POST /api/shorten
func shortenURL(w http.ResponseWriter, r *http.Request) {
	var body shortenRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if body.URL == "" {
		http.Error(w, "URL field is required", http.StatusBadRequest)
		return
	}

	code := generateShortCode()
	shortURL := "http://localhost:8080/" + code

	response := shortenResponse{ShortURL: shortURL}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
