package main

import (
	"fmt"
	"net/http"
)

func main() {
	// Connect to PostgreSQL first
	initDB()

	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/", redirectHandler)

	fmt.Println("🚀 Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
