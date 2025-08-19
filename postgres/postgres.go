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

func InsertUpdate(mangaName string, tableName string) (string, error) {

	// Get the database connection
	db, err := Connect()
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	defer db.Close() // Close when the program exits

	// Insert a new entry or update if it exists
	_, err = db.Exec("INSERT INTO manga (name) VALUES ($1) ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name", mangaName)
	if err != nil {
		log.Fatalf("Insert/Update failed: %v", err)
	}

	log.Printf("Successfully inserted/updated manga: %s into table: %s", mangaName, tableName)

	return mangaName, err
}

// Insert new entry into table if the name does not already exist and returns the name if successful
func InsertIfNotExists(mangaName string, tableName string) (string, error) {
	db, err := Connect()
	if err != nil {
		return "", fmt.Errorf("database connection error: %w", err)
	}
	defer db.Close()

	// Check if name exists (case-insensitive)
	checkQuery := fmt.Sprintf("SELECT name FROM %s WHERE LOWER(name) = LOWER($1)", tableName)
	var existingName string
	err = db.QueryRow(checkQuery, mangaName).Scan(&existingName)
	if err == nil {
		// name already exists, return it
		return existingName, nil
	}
	if err != sql.ErrNoRows {
		// some other error happened
		return "", err
	}

	// Insert since not existing
	insertQuery := fmt.Sprintf("INSERT INTO %s (name) VALUES ($1) RETURNING name", tableName)
	var newName string
	if err := db.QueryRow(insertQuery, mangaName).Scan(&newName); err != nil {
		return "", fmt.Errorf("failed to insert: %w", err)
	}

	return newName, nil
}

// Inserts managa name if the entry dos not already exist and returns full row as a map
func InsertIfNotExistsFull(mangaName string, tableName string) (map[string]any, error) {
	db, err := Connect()
	if err != nil {
		return nil, fmt.Errorf("database connection error: %w", err)
	}
	defer db.Close()

	// 1. Check if name exists (case-insensitive)
	checkQuery := fmt.Sprintf("SELECT * FROM %s WHERE LOWER(name) = LOWER($1)", tableName)
	rows, err := db.Query(checkQuery, mangaName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		// Existing row found
		cols, _ := rows.Columns()
		vals := make([]any, len(cols))
		valPtrs := make([]any, len(cols))
		for i := range vals {
			valPtrs[i] = &vals[i]
		}

		if err := rows.Scan(valPtrs...); err != nil {
			return nil, err
		}

		rowMap := make(map[string]any)
		for i, col := range cols {
			rowMap[col] = vals[i]
		}
		return rowMap, nil
	}

	// 2. Insert since not existing
	insertQuery := fmt.Sprintf("INSERT INTO %s (name) VALUES ($1) RETURNING *", tableName)
	insertRows, err := db.Query(insertQuery, mangaName)
	if err != nil {
		return nil, fmt.Errorf("failed to insert: %w", err)
	}
	defer insertRows.Close()

	if insertRows.Next() {
		cols, _ := insertRows.Columns()
		vals := make([]any, len(cols))
		valPtrs := make([]any, len(cols))
		for i := range vals {
			valPtrs[i] = &vals[i]
		}

		if err := insertRows.Scan(valPtrs...); err != nil {
			return nil, err
		}

		rowMap := make(map[string]any)
		for i, col := range cols {
			rowMap[col] = vals[i]
		}
		return rowMap, nil
	}

	return nil, fmt.Errorf("insert failed, no rows returned")
}
