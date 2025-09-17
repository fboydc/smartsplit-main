package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {

	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=password dbname=cust_data sslmode=disable")
	if err != nil {
		return nil, err
	}
	DB = db
	return db, nil
}

func CloseDB(db *sql.DB) {
	db.Close()
}

func executeQuery(query string, db *sql.DB) (*sql.Rows, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
