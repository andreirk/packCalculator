// service/pack_service_test.go
package service

import (
	"database/sql"
	"packCalculator/server/repository"
	"reflect"
	"sync"
	"testing"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}

	// Create pack_sizes table
	_, err = db.Exec(`
		CREATE TABLE pack_sizes (
			size INTEGER PRIMARY KEY
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create pack_sizes table: %v", err)
	}

	return db
}

func TestNewPackService(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert test data
	_, err := db.Exec("INSERT INTO pack_sizes (size) VALUES (250), (500), (1000)")
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	repo := repository.NewPackRepository(db)
	service, err := NewPackService(repo)
	if err != nil {
		t.Fatalf("Failed to create PackService: %v", err)
	}

	expected := []int{1000, 500, 250}
	actual := service.GetCurrentSizes()

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected pack sizes %v, got %v", expected, actual)
	}
}

func TestCalculatePacks(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert test data
	_, err := db.Exec("INSERT INTO pack_sizes (size) VALUES (250), (500), (1000), (2000), (5000)")
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	repo := repository.NewPackRepository(db)
	service, err := NewPackService(repo)
	if err != nil {
		t.Fatalf("Failed to create PackService: %v", err)
	}

	tests := []struct {
		name          string
		itemsOrdered  int
		expectedPacks map[int]int
	}{
		{
			name:          "1 item",
			itemsOrdered:  1,
			expectedPacks: map[int]int{250: 1},
		},
		{
			name:          "250 items",
			itemsOrdered:  250,
			expectedPacks: map[int]int{250: 1},
		},
		{
			name:          "251 items",
			itemsOrdered:  251,
			expectedPacks: map[int]int{500: 1},
		},
		{
			name:          "501 items",
			itemsOrdered:  501,
			expectedPacks: map[int]int{500: 1, 250: 1},
		},
		{
			name:          "12001 items",
			itemsOrdered:  12001,
			expectedPacks: map[int]int{5000: 2, 2000: 1, 250: 1},
		},
		{
			name:          "Zero items",
			itemsOrdered:  0,
			expectedPacks: map[int]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.CalculatePacks(tt.itemsOrdered)

			if !reflect.DeepEqual(result, tt.expectedPacks) {
				t.Errorf("CalculatePacks() = %v, want %v", result, tt.expectedPacks)
			}
		})
	}
}

func TestAddPackSize(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Initial data
	_, err := db.Exec("INSERT INTO pack_sizes (size) VALUES (250), (500)")
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	repo := repository.NewPackRepository(db)
	service, err := NewPackService(repo)
	if err != nil {
		t.Fatalf("Failed to create PackService: %v", err)
	}

	// Test adding a new size
	err = service.AddPackSize(1000)
	if err != nil {
		t.Fatalf("AddPackSize failed: %v", err)
	}

	expected := []int{1000, 500, 250}
	actual := service.GetCurrentSizes()

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("After AddPackSize, expected %v, got %v", expected, actual)
	}

	// Verify it was added to the database
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM pack_sizes WHERE size = 1000").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query database: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected size 1000 to be in database, but count was %d", count)
	}
}

func TestRemovePackSize(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Initial data
	_, err := db.Exec("INSERT INTO pack_sizes (size) VALUES (250), (500), (1000)")
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	repo := repository.NewPackRepository(db)
	service, err := NewPackService(repo)
	if err != nil {
		t.Fatalf("Failed to create PackService: %v", err)
	}

	// Test removing a size
	err = service.RemovePackSize(500)
	if err != nil {
		t.Fatalf("RemovePackSize failed: %v", err)
	}

	expected := []int{1000, 250}
	actual := service.GetCurrentSizes()

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("After RemovePackSize, expected %v, got %v", expected, actual)
	}

	// Verify it was removed from the database
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM pack_sizes WHERE size = 500").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query database: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected size 500 to be removed from database, but count was %d", count)
	}
}

func TestUpdatePackSizes(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Initial data
	_, err := db.Exec("INSERT INTO pack_sizes (size) VALUES (250), (500), (1000)")
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	repo := repository.NewPackRepository(db)
	service, err := NewPackService(repo)
	if err != nil {
		t.Fatalf("Failed to create PackService: %v", err)
	}

	// Test updating all sizes
	newSizes := []int{300, 600, 1200}
	err = service.UpdatePackSizes(newSizes)
	if err != nil {
		t.Fatalf("UpdatePackSizes failed: %v", err)
	}

	expected := []int{1200, 600, 300}
	actual := service.GetCurrentSizes()

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("After UpdatePackSizes, expected %v, got %v", expected, actual)
	}

	// Verify database was updated
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM pack_sizes").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query database: %v", err)
	}
	if count != 3 {
		t.Errorf("Expected 3 sizes in database, got %d", count)
	}
}

func TestConcurrentAccess(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Initial data
	_, err := db.Exec("INSERT INTO pack_sizes (size) VALUES (250), (500), (1000)")
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	repo := repository.NewPackRepository(db)
	service, err := NewPackService(repo)
	if err != nil {
		t.Fatalf("Failed to create PackService: %v", err)
	}

	var wg sync.WaitGroup
	iterations := 100

	// Test concurrent access to CalculatePacks
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			_ = service.CalculatePacks(i)
		}
	}()

	// Test concurrent access to GetCurrentSizes
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			_ = service.GetCurrentSizes()
		}
	}()

	// Test concurrent modifications
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			if i%2 == 0 {
				_ = service.AddPackSize(2000 + i)
			} else {
				_ = service.RemovePackSize(2000 + i - 1)
			}
		}
	}()

	wg.Wait()

	// Verify the service is still in a valid state
	sizes := service.GetCurrentSizes()
	if len(sizes) == 0 {
		t.Error("Expected some pack sizes after concurrent access, got none")
	}
}
