package db

import (
	"database/sql"
	"fmt"
)

type Migration struct {
	Version int
	Up      func(*sql.DB) error
}

var migrations = []Migration{
	{
		Version: 1,
		Up: func(db *sql.DB) error {
			_, err := db.Exec(`
				CREATE TABLE IF NOT EXISTS pack_sizes (
					size INTEGER PRIMARY KEY
				)
			`)
			return err
		},
	},
	// Add more migrations as needed
}

func (d *Database) Migrate() error {
	// Create migrations table if it doesn't exist
	if _, err := d.DB.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			version INTEGER PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`); err != nil {
		return fmt.Errorf("failed to create migrations table: %v", err)
	}

	// Get current version
	var currentVersion int
	err := d.DB.QueryRow("SELECT MAX(version) FROM migrations").Scan(&currentVersion)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to get current migration version: %v", err)
	}

	// Apply missing migrations
	for _, migration := range migrations {
		if migration.Version > currentVersion {
			if err := migration.Up(d.DB); err != nil {
				return fmt.Errorf("failed to apply migration %d: %v", migration.Version, err)
			}

			if _, err := d.DB.Exec("INSERT INTO migrations (version) VALUES (?)", migration.Version); err != nil {
				return fmt.Errorf("failed to record migration %d: %v", migration.Version, err)
			}
		}
	}

	return nil
}
