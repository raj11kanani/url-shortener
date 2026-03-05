package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateCode() string {
	code := make([]byte, 6)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}

// POST /shorten
func shortenHandler(w http.ResponseWriter, r *http.Request) {
	// Allow requests from Chrome extension
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.URL == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Generate a unique code
	code := generateCode()

	// Save to PostgreSQL instead of in-memory map
	if err := saveURL(code, body.URL); err != nil {
		http.Error(w, "Failed to save URL", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"short_url": fmt.Sprintf("%s/%s", getEnv("APP_URL", "http://localhost:8080"), code),
		"code":      code,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GET /:code
func redirectHandler(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/")

	if code == "" {
		http.Error(w, "No code provided", http.StatusBadRequest)
		return
	}

	// Look up in PostgreSQL instead of in-memory map
	originalURL, err := getURL(code)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	// Track the click
	incrementClicks(code)

	http.Redirect(w, r, originalURL, http.StatusFound)
}
