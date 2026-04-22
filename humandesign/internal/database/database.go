package database

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// Person represents a person in the database
type Person struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	BirthDate string    `json:"birth_date"` // YYYY-MM-DD
	BirthTime string    `json:"birth_time"` // HH:MM
	Location  string    `json:"location"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Notes     string    `json:"notes,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// Database provides file-backed storage for people
type Database struct {
	mu       sync.RWMutex
	filePath string
	people   []Person
}

// NewDatabase creates a new database backed by the given JSON file
func NewDatabase(filePath string) (*Database, error) {
	db := &Database{
		filePath: filePath,
		people:   []Person{},
	}

	// Load existing data if file exists
	if _, err := os.Stat(filePath); err == nil {
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("reading database file: %w", err)
		}
		if len(data) > 0 {
			if err := json.Unmarshal(data, &db.people); err != nil {
				return nil, fmt.Errorf("parsing database file: %w", err)
			}
		}
	}

	return db, nil
}

// List returns all people in the database
func (db *Database) List() []Person {
	db.mu.RLock()
	defer db.mu.RUnlock()
	result := make([]Person, len(db.people))
	copy(result, db.people)
	return result
}

// Get returns a person by ID
func (db *Database) Get(id string) (*Person, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	for _, p := range db.people {
		if p.ID == id {
			return &p, nil
		}
	}
	return nil, fmt.Errorf("person not found: %s", id)
}

// Add adds a new person to the database
func (db *Database) Add(p Person) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Generate ID from timestamp if not provided
	if p.ID == "" {
		p.ID = fmt.Sprintf("p%d", time.Now().UnixNano())
	}
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now()
	}

	db.people = append(db.people, p)
	return db.save()
}

// Update replaces a person's data by ID
func (db *Database) Update(p Person) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	for i, existing := range db.people {
		if existing.ID == p.ID {
			p.CreatedAt = existing.CreatedAt
			db.people[i] = p
			return db.save()
		}
	}
	return fmt.Errorf("person not found: %s", p.ID)
}

// Delete removes a person by ID
func (db *Database) Delete(id string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	for i, p := range db.people {
		if p.ID == id {
			db.people = append(db.people[:i], db.people[i+1:]...)
			return db.save()
		}
	}
	return fmt.Errorf("person not found: %s", id)
}

// save writes the database to disk
func (db *Database) save() error {
	data, err := json.MarshalIndent(db.people, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling database: %w", err)
	}
	return os.WriteFile(db.filePath, data, 0644)
}
