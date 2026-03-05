package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "admin")
	password := getEnv("DB_PASSWORD", "password")
	dbname := getEnv("DB_NAME", "urlshortener")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		host, port, user, password, dbname,
	)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("❌ Failed to open DB connection:", err)
	}

	// Retry until PostgreSQL is ready (max 10 attempts)
	for i := 1; i <= 10; i++ {
		err = db.Ping()
		if err == nil {
			break
		}
		fmt.Printf("⏳ Waiting for DB... attempt %d/10\n", i)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("❌ Could not connect to DB after 10 attempts:", err)
	}

	createTable()
	fmt.Println("✅ Connected to PostgreSQL successfully!")
}

func createTable() {
	query := `
	CREATE TABLE IF NOT EXISTS urls (
		id        SERIAL PRIMARY KEY,
		code      VARCHAR(10) UNIQUE NOT NULL,
		original  TEXT NOT NULL,
		clicks    INTEGER DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("❌ Failed to create table:", err)
	}
}

func saveURL(code, originalURL string) error {
	_, err := db.Exec(
		"INSERT INTO urls (code, original) VALUES ($1, $2)",
		code, originalURL,
	)
	return err
}

func getURL(code string) (string, error) {
	var original string
	err := db.QueryRow(
		"SELECT original FROM urls WHERE code = $1", code,
	).Scan(&original)
	return original, err
}

func incrementClicks(code string) {
	db.Exec("UPDATE urls SET clicks = clicks + 1 WHERE code = $1", code)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
