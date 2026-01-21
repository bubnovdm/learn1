package db

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	DB_HOST     = "localhost"
	DB_PORT     = 5432
	DB_USER     = "postgres"
	DB_PASSWORD = "postgres"
	DB_NAME     = "mydb"
)

//var DB *sqlx.DB

// InitDB initializes the database connection
func InitDB() (*sqlx.DB, error) {
	var db *sqlx.DB
	// Build connection string
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME,
	)

	var err error
	db, err = sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL database")
	return db, nil
}
