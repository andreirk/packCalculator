// db/database.go
package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	DB *sql.DB
}

func NewDatabase(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return &Database{DB: db}, nil
}

func (d *Database) Initialize() error {
	// Create pack_sizes table if it doesn't exist
	_, err := d.DB.Exec(`
		CREATE TABLE IF NOT EXISTS pack_sizes (
			size INTEGER PRIMARY KEY
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create pack_sizes table: %v", err)
	}

	// Check if we need to insert default sizes
	var count int
	err = d.DB.QueryRow("SELECT COUNT(*) FROM pack_sizes").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to count pack sizes: %v", err)
	}

	if count == 0 {
		defaultSizes := []int{250, 500, 1000, 2000, 5000}
		for _, size := range defaultSizes {
			_, err := d.DB.Exec("INSERT OR IGNORE INTO pack_sizes (size) VALUES (?)", size)
			if err != nil {
				return fmt.Errorf("failed to insert default pack size %d: %v", size, err)
			}
		}
		log.Println("Initialized database with default pack sizes")
	}

	return nil
}

func (d *Database) Close() error {
	return d.DB.Close()
}
