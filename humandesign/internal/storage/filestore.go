package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// FileStore implements Storage interface using filesystem
type FileStore struct {
	dataDir string
	mu      sync.RWMutex // For thread-safe operations
}

// NewFileStore creates a new filesystem-based storage
func NewFileStore(dataDir string) (*FileStore, error) {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	return &FileStore{
		dataDir: dataDir,
	}, nil
}

// Create creates a new person record
func (fs *FileStore) Create(person *Person) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	// Generate ID if not set
	if person.ID == "" {
		person.ID = uuid.New().String()
	}

	// Set timestamps
	now := time.Now()
	person.CreatedAt = now
	person.UpdatedAt = now
	person.Version = 1

	// Check if person with this ID already exists
	filename := fs.getFilename(person.ID)
	if _, err := os.Stat(filename); err == nil {
		return fmt.Errorf("person with ID %s already exists", person.ID)
	}

	// Write to file
	return fs.writePersonToFile(person)
}

// Get retrieves a person by ID
func (fs *FileStore) Get(id string) (*Person, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	filename := fs.getFilename(id)
	return fs.readPersonFromFile(filename)
}

// GetByName retrieves a person by name
func (fs *FileStore) GetByName(name string) (*Person, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	// List all files and find matching name
	files, err := os.ReadDir(fs.dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read data directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		filename := filepath.Join(fs.dataDir, file.Name())
		person, err := fs.readPersonFromFile(filename)
		if err != nil {
			continue // Skip invalid files
		}

		if person.Name == name {
			return person, nil
		}
	}

	return nil, fmt.Errorf("person with name %s not found", name)
}

// Update updates an existing person record
func (fs *FileStore) Update(person *Person) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	// Check if person exists
	filename := fs.getFilename(person.ID)
	existing, err := fs.readPersonFromFile(filename)
	if err != nil {
		return fmt.Errorf("person not found: %w", err)
	}

	// Optimistic locking check
	if person.Version != existing.Version {
		return fmt.Errorf("version conflict: expected %d, got %d", existing.Version, person.Version)
	}

	// Update timestamps and version
	person.UpdatedAt = time.Now()
	person.Version++

	// Preserve creation time
	person.CreatedAt = existing.CreatedAt

	// Write to file
	return fs.writePersonToFile(person)
}

// Delete deletes a person by ID
func (fs *FileStore) Delete(id string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	filename := fs.getFilename(id)
	if err := os.Remove(filename); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("person with ID %s not found", id)
		}
		return fmt.Errorf("failed to delete person: %w", err)
	}

	return nil
}

// List returns all people (optionally filtered by tags)
func (fs *FileStore) List(tags []string) ([]PersonSummary, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	files, err := os.ReadDir(fs.dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read data directory: %w", err)
	}

	var summaries []PersonSummary

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		filename := filepath.Join(fs.dataDir, file.Name())
		person, err := fs.readPersonFromFile(filename)
		if err != nil {
			continue // Skip invalid files
		}

		// Filter by tags if specified
		if len(tags) > 0 && !hasAnyTag(person.Tags, tags) {
			continue
		}

		summaries = append(summaries, person.ToSummary())
	}

	return summaries, nil
}

// Search searches for people by name (partial match)
func (fs *FileStore) Search(query string) ([]PersonSummary, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	files, err := os.ReadDir(fs.dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read data directory: %w", err)
	}

	query = strings.ToLower(query)
	var summaries []PersonSummary

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		filename := filepath.Join(fs.dataDir, file.Name())
		person, err := fs.readPersonFromFile(filename)
		if err != nil {
			continue // Skip invalid files
		}

		if strings.Contains(strings.ToLower(person.Name), query) {
			summaries = append(summaries, person.ToSummary())
		}
	}

	return summaries, nil
}

// Close closes the storage (no-op for filesystem)
func (fs *FileStore) Close() error {
	return nil
}

// Helper methods

func (fs *FileStore) getFilename(id string) string {
	return filepath.Join(fs.dataDir, fmt.Sprintf("%s.json", id))
}

func (fs *FileStore) writePersonToFile(person *Person) error {
	filename := fs.getFilename(person.ID)

	// Marshal to JSON with indentation for readability
	data, err := json.MarshalIndent(person, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal person: %w", err)
	}

	// Write to temporary file first (atomic write)
	tmpFile := filename + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	// Rename to final filename (atomic on most systems)
	if err := os.Rename(tmpFile, filename); err != nil {
		os.Remove(tmpFile) // Clean up
		return fmt.Errorf("failed to rename file: %w", err)
	}

	return nil
}

func (fs *FileStore) readPersonFromFile(filename string) (*Person, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("person not found")
		}
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var person Person
	if err := json.Unmarshal(data, &person); err != nil {
		return nil, fmt.Errorf("failed to unmarshal person: %w", err)
	}

	return &person, nil
}

func hasAnyTag(personTags, filterTags []string) bool {
	tagSet := make(map[string]bool)
	for _, tag := range personTags {
		tagSet[tag] = true
	}

	for _, tag := range filterTags {
		if tagSet[tag] {
			return true
		}
	}

	return false
}
