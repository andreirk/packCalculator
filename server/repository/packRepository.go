// repository/pack_repository.go
package repository

import (
	"database/sql"
)

type PackRepository struct {
	db *sql.DB
}

func NewPackRepository(db *sql.DB) *PackRepository {
	return &PackRepository{db: db}
}

func (r *PackRepository) GetAll() ([]int, error) {
	rows, err := r.db.Query("SELECT size FROM pack_sizes ORDER BY size DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sizes []int
	for rows.Next() {
		var size int
		if err := rows.Scan(&size); err != nil {
			return nil, err
		}
		sizes = append(sizes, size)
	}

	return sizes, nil
}

func (r *PackRepository) Add(size int) error {
	_, err := r.db.Exec("INSERT OR IGNORE INTO pack_sizes (size) VALUES (?)", size)
	return err
}

func (r *PackRepository) Remove(size int) error {
	_, err := r.db.Exec("DELETE FROM pack_sizes WHERE size = ?", size)
	return err
}

func (r *PackRepository) ReplaceAll(sizes []int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Clear existing sizes
	_, err = tx.Exec("DELETE FROM pack_sizes")
	if err != nil {
		tx.Rollback()
		return err
	}

	// Insert new sizes
	stmt, err := tx.Prepare("INSERT INTO pack_sizes (size) VALUES (?)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, size := range sizes {
		if _, err := stmt.Exec(size); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
