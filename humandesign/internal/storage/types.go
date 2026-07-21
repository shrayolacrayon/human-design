package storage

import (
	"humandesign/internal/calculator"
	"time"
)

// Person represents a person with their birth data and Human Design reading
// This schema is designed to be database-compatible
type Person struct {
	// Primary key (will be UUID in database)
	ID string `json:"id"`

	// Person information
	Name  string `json:"name"`
	Email string `json:"email,omitempty"` // Optional
	Notes string `json:"notes,omitempty"` // Optional notes about the person

	// Birth data
	BirthData calculator.BirthData `json:"birth_data"`

	// Calculated reading (stored to avoid recalculation)
	Reading *calculator.Reading `json:"reading,omitempty"`

	// Metadata (for database tracking)
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int       `json:"version"` // For optimistic locking in database

	// Tags for categorization
	Tags []string `json:"tags,omitempty"`
}

// PersonSummary is a lightweight version for listing
type PersonSummary struct {
	ID        string                  `json:"id"`
	Name      string                  `json:"name"`
	Type      calculator.HumanDesignType `json:"type,omitempty"`
	Authority calculator.Authority    `json:"authority,omitempty"`
	Profile   string                  `json:"profile,omitempty"`
	CreatedAt time.Time               `json:"created_at"`
	Tags      []string                `json:"tags,omitempty"`
}

// ToSummary converts a Person to a PersonSummary
func (p *Person) ToSummary() PersonSummary {
	summary := PersonSummary{
		ID:        p.ID,
		Name:      p.Name,
		CreatedAt: p.CreatedAt,
		Tags:      p.Tags,
	}

	if p.Reading != nil {
		summary.Type = p.Reading.Type
		summary.Authority = p.Reading.Authority
		if p.Reading.Profile.Conscious > 0 && p.Reading.Profile.Unconscious > 0 {
			summary.Profile = formatProfile(p.Reading.Profile)
		}
	}

	return summary
}

func formatProfile(profile calculator.Profile) string {
	return profile.Name
}

// Storage defines the interface for person storage
// This can be implemented by different backends (filesystem, SQL, NoSQL)
type Storage interface {
	// Create creates a new person record
	Create(person *Person) error

	// Get retrieves a person by ID
	Get(id string) (*Person, error)

	// GetByName retrieves a person by name (for convenience)
	GetByName(name string) (*Person, error)

	// Update updates an existing person record
	Update(person *Person) error

	// Delete deletes a person by ID
	Delete(id string) error

	// List returns all people (optionally filtered by tags)
	List(tags []string) ([]PersonSummary, error)

	// Search searches for people by name (partial match)
	Search(query string) ([]PersonSummary, error)

	// Close closes the storage connection
	Close() error
}

// StorageConfig holds configuration for storage backends
type StorageConfig struct {
	// For filesystem storage
	DataDir string

	// For future database storage
	DatabaseURL string
	DatabaseType string // "postgres", "mysql", "sqlite", etc.
}
