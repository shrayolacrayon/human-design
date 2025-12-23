# Storage System

The Human Design project includes a flexible storage system for managing person records with their Human Design readings.

## Overview

The storage system uses a filesystem-based JSON database by default, with a schema designed for easy migration to a relational or NoSQL database in the future.

## Architecture

### Storage Interface

The `Storage` interface (`internal/storage/types.go`) defines the contract for all storage backends:

```go
type Storage interface {
    Create(person *Person) error
    Get(id string) (*Person, error)
    GetByName(name string) (*Person, error)
    Update(person *Person) error
    Delete(id string) error
    List(tags []string) ([]PersonSummary, error)
    Search(query string) ([]PersonSummary, error)
    Close() error
}
```

### Person Schema

The `Person` type is designed to be database-compatible:

```go
type Person struct {
    ID        string                 // UUID primary key
    Name      string                 // Person's name
    Email     string                 // Optional email
    Notes     string                 // Optional notes
    BirthData calculator.BirthData   // Birth information
    Reading   *calculator.Reading    // Calculated Human Design reading
    CreatedAt time.Time              // Creation timestamp
    UpdatedAt time.Time              // Last update timestamp
    Version   int                    // For optimistic locking
    Tags      []string               // Categorization tags
}
```

## Filesystem Storage

The default implementation stores each person as a separate JSON file.

### Directory Structure

```
data/
└── people/
    ├── 550e8400-e29b-41d4-a716-446655440000.json
    ├── 6ba7b810-9dad-11d1-80b4-00c04fd430c8.json
    └── ...
```

### File Format

Each JSON file contains a complete `Person` record:

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "John Doe",
  "email": "john@example.com",
  "notes": "Family member",
  "birth_data": {
    "datetime": "1990-06-15T14:30:00Z",
    "latitude": 40.7128,
    "longitude": -74.0060,
    "location": "New York, NY"
  },
  "reading": {
    "type": "Generator",
    "authority": "Sacral",
    "profile": {
      "conscious": 1,
      "unconscious": 3,
      "name": "Investigator/Martyr"
    },
    ...
  },
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z",
  "version": 1,
  "tags": ["family", "friend"]
}
```

## CLI Usage

### Save a Person

```bash
humandesign-cli save \
  -name "John Doe" \
  -date "1990-06-15T14:30:00Z" \
  -lat 40.7128 \
  -lon -74.0060 \
  -location "New York, NY" \
  -email "john@example.com" \
  -notes "Family member" \
  -tags "family,friends"
```

### Load a Person

By name:
```bash
humandesign-cli load -name "John Doe"
```

By ID:
```bash
humandesign-cli load -id "550e8400-e29b-41d4-a716-446655440000"
```

JSON output:
```bash
humandesign-cli load -name "John Doe" -output json
```

### List All People

```bash
# List all
humandesign-cli list

# Filter by tags
humandesign-cli list -tags "family,friends"

# JSON output
humandesign-cli list -output json
```

### Search for People

```bash
humandesign-cli search -query "John"
humandesign-cli search -query "Doe" -output json
```

### Delete a Person

```bash
# By name
humandesign-cli delete -name "John Doe"

# By ID
humandesign-cli delete -id "550e8400-e29b-41d4-a716-446655440000"
```

## Programmatic Usage

### Creating a Person

```go
package main

import (
    "humandesign/internal/calculator"
    "humandesign/internal/storage"
    "time"
)

func main() {
    // Initialize storage
    store, err := storage.NewFileStore("./data/people")
    if err != nil {
        panic(err)
    }
    defer store.Close()

    // Create birth data
    birthData := calculator.BirthData{
        DateTime:  time.Date(1990, 6, 15, 14, 30, 0, 0, time.UTC),
        Latitude:  40.7128,
        Longitude: -74.0060,
        Location:  "New York, NY",
    }

    // Calculate reading
    calc := calculator.NewCalculator()
    reading, err := calc.Calculate(birthData)
    if err != nil {
        panic(err)
    }

    // Create person
    person := &storage.Person{
        Name:      "John Doe",
        Email:     "john@example.com",
        BirthData: birthData,
        Reading:   reading,
        Tags:      []string{"family", "friends"},
    }

    // Save to storage
    if err := store.Create(person); err != nil {
        panic(err)
    }

    println("Saved person with ID:", person.ID)
}
```

### Retrieving a Person

```go
// By ID
person, err := store.Get("550e8400-e29b-41d4-a716-446655440000")
if err != nil {
    panic(err)
}

// By name
person, err := store.GetByName("John Doe")
if err != nil {
    panic(err)
}
```

### Updating a Person

```go
person, err := store.GetByName("John Doe")
if err != nil {
    panic(err)
}

person.Notes = "Updated notes"
person.Tags = append(person.Tags, "new-tag")

if err := store.Update(person); err != nil {
    panic(err)
}
```

### Listing and Searching

```go
// List all people
summaries, err := store.List(nil)

// Filter by tags
summaries, err := store.List([]string{"family", "friends"})

// Search by name
summaries, err := store.Search("John")
```

## Features

### Optimistic Locking

The storage system uses version numbers for optimistic locking:

```go
person, _ := store.Get(id)
// ... someone else updates the person ...

person.Notes = "My update"
err := store.Update(person)
// Returns error if version mismatch
```

### Atomic Writes

File writes are atomic using a two-phase commit:
1. Write to temporary file (`.json.tmp`)
2. Atomic rename to final filename

This prevents data corruption from partial writes.

### Thread Safety

The filesystem storage uses mutexes for thread-safe concurrent access.

## Future: Database Migration

The schema is designed for easy migration to SQL databases:

### PostgreSQL Example

```sql
CREATE TABLE people (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    notes TEXT,
    birth_data JSONB NOT NULL,
    reading JSONB,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    version INTEGER NOT NULL DEFAULT 1,
    tags TEXT[]
);

CREATE INDEX idx_people_name ON people(name);
CREATE INDEX idx_people_tags ON people USING GIN(tags);
CREATE INDEX idx_people_created_at ON people(created_at);
```

### Implementation

To add a PostgreSQL backend:

1. Implement the `Storage` interface:
```go
type PostgresStore struct {
    db *sql.DB
}

func (ps *PostgresStore) Create(person *Person) error {
    // INSERT INTO people ...
}
// ... implement other methods
```

2. Use the storage abstraction:
```go
var store storage.Storage

if config.UsePostgres {
    store = NewPostgresStore(dbURL)
} else {
    store = NewFileStore(dataDir)
}
```

## Custom Data Directory

Specify a custom data directory:

```bash
# CLI
humandesign-cli save -name "John" -datadir /custom/path ...

# Programmatically
store, err := storage.NewFileStore("/custom/path")
```

## Backup and Export

### Manual Backup

Simply copy the data directory:
```bash
cp -r ./data/people ./backups/people-$(date +%Y%m%d)
```

### Export to CSV

```go
summaries, _ := store.List(nil)
// Convert summaries to CSV format
```

## Performance Considerations

### Filesystem Storage

- **Small scale** (<1000 people): Excellent performance
- **Medium scale** (1000-10000 people): Good performance, consider indexing for search
- **Large scale** (>10000 people): Consider migrating to a database

### Optimization Tips

1. **Use tags** for categorization instead of searching
2. **Cache frequently accessed** people in memory
3. **Batch operations** when possible
4. **Consider database** for production with many users

## Error Handling

All storage operations return errors. Always check them:

```go
person, err := store.Get(id)
if err != nil {
    if strings.Contains(err.Error(), "not found") {
        // Handle not found
    } else {
        // Handle other errors
    }
}
```

## Testing

Run storage tests:

```bash
go test ./internal/storage -v
```

The tests use temporary directories and clean up automatically.

## Best Practices

1. **Always close storage** when done:
   ```go
   defer store.Close()
   ```

2. **Use transactions** (when migrating to DB):
   Store readings atomically with person data

3. **Validate data** before storing:
   Ensure birth data is complete and valid

4. **Use tags** for organization:
   Better than relying on search alone

5. **Regular backups**:
   Especially important for filesystem storage

## Troubleshooting

### Permission Errors

Ensure the data directory is writable:
```bash
chmod 755 ./data/people
```

### Corruption

If a JSON file becomes corrupted:
1. Restore from backup, or
2. Manually edit the JSON file, or
3. Delete the file (data will be lost)

### Version Conflicts

If you get version conflicts:
1. Reload the person
2. Apply your changes
3. Save again

## Migration Guide

See [MIGRATIONS.md](MIGRATIONS.md) for guides on:
- Migrating from filesystem to PostgreSQL
- Migrating from filesystem to MongoDB
- Importing from CSV
- Exporting to different formats
