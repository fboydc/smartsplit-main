package services

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// DatabaseService handles database operations
type DatabaseService struct {
	db *sql.DB
}

// NewDatabaseService creates a new database service instance
func NewDatabaseService() *DatabaseService {
	return &DatabaseService{}
}

// InitDB initializes the database connection
func (ds *DatabaseService) InitDB() (*sql.DB, error) {
	// Database connection string - should come from environment variables in production
	connStr := "host=localhost port=5432 user=postgres password=password dbname=cust_data sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	ds.db = db
	return db, nil
}

// CloseDB closes the database connection
func (ds *DatabaseService) CloseDB(db *sql.DB) error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// ExecuteQuery executes a query and returns the rows
func (ds *DatabaseService) ExecuteQuery(query string, db *sql.DB) (*sql.Rows, error) {
	return db.Query(query)
}

// Global functions for backward compatibility (these should eventually be refactored)

// InitDB initializes database connection (global function for compatibility)
func InitDB() (*sql.DB, error) {
	service := NewDatabaseService()
	return service.InitDB()
}

// CloseDB closes database connection (global function for compatibility)
func CloseDB(db *sql.DB) error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// ExecuteQuery executes a query (global function for compatibility)
func ExecuteQuery(query string, db *sql.DB) (*sql.Rows, error) {
	return db.Query(query)
}
