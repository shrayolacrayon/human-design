package storage

import (
	"human

design/internal/calculator"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFileStoreCreateAndGet(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	// Create a person
	person := &Person{
		Name:  "Test Person",
		Email: "test@example.com",
		BirthData: calculator.BirthData{
			DateTime:  time.Date(1990, 6, 15, 14, 30, 0, 0, time.UTC),
			Latitude:  40.7128,
			Longitude: -74.0060,
			Location:  "New York, NY",
		},
		Tags: []string{"test", "example"},
	}

	// Save person
	err = store.Create(person)
	if err != nil {
		t.Fatalf("Failed to create person: %v", err)
	}

	// Verify ID was generated
	if person.ID == "" {
		t.Error("ID was not generated")
	}

	// Verify timestamps were set
	if person.CreatedAt.IsZero() {
		t.Error("CreatedAt was not set")
	}
	if person.UpdatedAt.IsZero() {
		t.Error("UpdatedAt was not set")
	}

	// Retrieve person by ID
	retrieved, err := store.Get(person.ID)
	if err != nil {
		t.Fatalf("Failed to get person: %v", err)
	}

	// Verify data
	if retrieved.Name != person.Name {
		t.Errorf("Expected name %s, got %s", person.Name, retrieved.Name)
	}
	if retrieved.Email != person.Email {
		t.Errorf("Expected email %s, got %s", person.Email, retrieved.Email)
	}
}

func TestFileStoreGetByName(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	person := &Person{
		Name: "John Doe",
		BirthData: calculator.BirthData{
			DateTime: time.Now(),
			Location: "Test Location",
		},
	}

	err = store.Create(person)
	if err != nil {
		t.Fatalf("Failed to create person: %v", err)
	}

	// Retrieve by name
	retrieved, err := store.GetByName("John Doe")
	if err != nil {
		t.Fatalf("Failed to get person by name: %v", err)
	}

	if retrieved.ID != person.ID {
		t.Errorf("Retrieved wrong person")
	}

	// Try non-existent name
	_, err = store.GetByName("Non Existent")
	if err == nil {
		t.Error("Expected error for non-existent person")
	}
}

func TestFileStoreUpdate(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	person := &Person{
		Name: "Test Person",
		BirthData: calculator.BirthData{
			DateTime: time.Now(),
			Location: "Original Location",
		},
	}

	err = store.Create(person)
	if err != nil {
		t.Fatalf("Failed to create person: %v", err)
	}

	originalVersion := person.Version
	originalCreatedAt := person.CreatedAt

	// Update person
	person.Notes = "Updated notes"
	person.Tags = []string{"updated"}

	err = store.Update(person)
	if err != nil {
		t.Fatalf("Failed to update person: %v", err)
	}

	// Verify version incremented
	if person.Version != originalVersion+1 {
		t.Errorf("Expected version %d, got %d", originalVersion+1, person.Version)
	}

	// Verify CreatedAt preserved
	if !person.CreatedAt.Equal(originalCreatedAt) {
		t.Error("CreatedAt was modified")
	}

	// Verify UpdatedAt changed
	if person.UpdatedAt.Equal(person.CreatedAt) {
		t.Error("UpdatedAt was not updated")
	}

	// Retrieve and verify
	retrieved, err := store.Get(person.ID)
	if err != nil {
		t.Fatalf("Failed to get person: %v", err)
	}

	if retrieved.Notes != "Updated notes" {
		t.Errorf("Notes not updated")
	}
}

func TestFileStoreDelete(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	person := &Person{
		Name: "Test Person",
		BirthData: calculator.BirthData{
			DateTime: time.Now(),
			Location: "Test Location",
		},
	}

	err = store.Create(person)
	if err != nil {
		t.Fatalf("Failed to create person: %v", err)
	}

	// Delete person
	err = store.Delete(person.ID)
	if err != nil {
		t.Fatalf("Failed to delete person: %v", err)
	}

	// Verify person is gone
	_, err = store.Get(person.ID)
	if err == nil {
		t.Error("Person still exists after deletion")
	}

	// Try deleting non-existent person
	err = store.Delete("non-existent-id")
	if err == nil {
		t.Error("Expected error when deleting non-existent person")
	}
}

func TestFileStoreList(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	// Create multiple people
	people := []*Person{
		{
			Name: "Person 1",
			BirthData: calculator.BirthData{
				DateTime: time.Now(),
				Location: "Location 1",
			},
			Tags: []string{"tag1", "tag2"},
		},
		{
			Name: "Person 2",
			BirthData: calculator.BirthData{
				DateTime: time.Now(),
				Location: "Location 2",
			},
			Tags: []string{"tag2", "tag3"},
		},
		{
			Name: "Person 3",
			BirthData: calculator.BirthData{
				DateTime: time.Now(),
				Location: "Location 3",
			},
			Tags: []string{"tag3"},
		},
	}

	for _, p := range people {
		if err := store.Create(p); err != nil {
			t.Fatalf("Failed to create person: %v", err)
		}
	}

	// List all
	summaries, err := store.List(nil)
	if err != nil {
		t.Fatalf("Failed to list people: %v", err)
	}

	if len(summaries) != 3 {
		t.Errorf("Expected 3 people, got %d", len(summaries))
	}

	// List by tag
	summaries, err = store.List([]string{"tag1"})
	if err != nil {
		t.Fatalf("Failed to list people by tag: %v", err)
	}

	if len(summaries) != 1 {
		t.Errorf("Expected 1 person with tag1, got %d", len(summaries))
	}

	summaries, err = store.List([]string{"tag2"})
	if err != nil {
		t.Fatalf("Failed to list people by tag: %v", err)
	}

	if len(summaries) != 2 {
		t.Errorf("Expected 2 people with tag2, got %d", len(summaries))
	}
}

func TestFileStoreSearch(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	// Create people
	people := []*Person{
		{
			Name: "John Doe",
			BirthData: calculator.BirthData{
				DateTime: time.Now(),
				Location: "Test",
			},
		},
		{
			Name: "Jane Doe",
			BirthData: calculator.BirthData{
				DateTime: time.Now(),
				Location: "Test",
			},
		},
		{
			Name: "Bob Smith",
			BirthData: calculator.BirthData{
				DateTime: time.Now(),
				Location: "Test",
			},
		},
	}

	for _, p := range people {
		if err := store.Create(p); err != nil {
			t.Fatalf("Failed to create person: %v", err)
		}
	}

	// Search for "Doe"
	summaries, err := store.Search("Doe")
	if err != nil {
		t.Fatalf("Failed to search: %v", err)
	}

	if len(summaries) != 2 {
		t.Errorf("Expected 2 results for 'Doe', got %d", len(summaries))
	}

	// Search for "john" (case insensitive)
	summaries, err = store.Search("john")
	if err != nil {
		t.Fatalf("Failed to search: %v", err)
	}

	if len(summaries) != 1 {
		t.Errorf("Expected 1 result for 'john', got %d", len(summaries))
	}
}

func TestFileStoreVersionConflict(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	person := &Person{
		Name: "Test Person",
		BirthData: calculator.BirthData{
			DateTime: time.Now(),
			Location: "Test",
		},
	}

	err = store.Create(person)
	if err != nil {
		t.Fatalf("Failed to create person: %v", err)
	}

	// Simulate concurrent update
	person.Notes = "Update 1"
	err = store.Update(person)
	if err != nil {
		t.Fatalf("Failed first update: %v", err)
	}

	// Try to update with old version
	oldPerson := *person
	oldPerson.Version = 1 // Old version
	oldPerson.Notes = "Update 2 (should fail)"

	err = store.Update(&oldPerson)
	if err == nil {
		t.Error("Expected version conflict error")
	}
}

func TestFileStoreAtomicWrite(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewFileStore(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	person := &Person{
		Name: "Test Person",
		BirthData: calculator.BirthData{
			DateTime: time.Now(),
			Location: "Test",
		},
	}

	err = store.Create(person)
	if err != nil {
		t.Fatalf("Failed to create person: %v", err)
	}

	// Verify no .tmp file left behind
	files, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to read dir: %v", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".tmp" {
			t.Error("Temporary file not cleaned up")
		}
	}
}
