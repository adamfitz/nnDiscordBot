package postgres

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"main/auth"
)

// Connect to the database

func Connect() (*sql.DB, error) {
	// Load the credentials
	config, _ := auth.LoadConfig()

	psqlConnect := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable",
		config.DbServer, config.DbPort, config.DbUser, config.DbPassword, config.DbName)

	// Open database connection
	dbConnection, err := sql.Open("postgres", psqlConnect)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Verify the connection is valid
	if err := dbConnection.Ping(); err != nil {
		dbConnection.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return dbConnection, nil
}

func DbVersion() string {
	// Get the database connection
	db, err := Connect()
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	defer db.Close() // Close when the program exits

	// Example: Query the database
	var version string
	err = db.QueryRow("SELECT version()").Scan(&version)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	return version
}

func InsertUpdate(mangaName string) (string, error) {

	// Get the database connection
	db, err := Connect()
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	defer db.Close() // Close when the program exits

	// Insert a new entry or update if it exists
	_, err = db.Exec("INSERT INTO manga (manga_name) VALUES ($1) ON CONFLICT (manga_name) DO UPDATE SET manga_name = EXCLUDED.manga_name", mangaName)
	if err != nil {
		log.Fatalf("Insert/Update failed: %v", err)
	}

	log.Printf("Successfully inserted/updated manga: %s", mangaName)

	return mangaName, err
}
