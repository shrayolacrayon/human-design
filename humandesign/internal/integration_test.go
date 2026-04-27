package integration_test

import (
	"encoding/json"
	"fmt"
	"humandesign/internal/calculator"
	"humandesign/internal/database"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestKnownChart verifies gate/line positions for a known birth chart.
// 09/16/1992 1:30 PM EDT (17:30 UTC), Nashua NH (42.7654, -71.4676)
func TestKnownChart(t *testing.T) {
	// 1:30 PM EDT = 17:30 UTC
	dt := time.Date(1992, 9, 16, 17, 30, 0, 0, time.UTC)
	data := calculator.BirthData{
		DateTime:  dt,
		Latitude:  42.7654,
		Longitude: -71.4676,
		Location:  "Nashua, US-NH",
	}

	calc := calculator.NewCalculator()
	reading, err := calc.Calculate(data)
	if err != nil {
		t.Fatalf("Calculate failed: %v", err)
	}

	// Print full positions as JSON for manual verification
	printPositions(t, reading)

	// Basic sanity checks
	if reading.Type == "" {
		t.Error("Type is empty")
	}
	if reading.Authority == "" {
		t.Error("Authority is empty")
	}
	if len(reading.PersonalityGates) != 13 {
		t.Errorf("Expected 13 personality gates (one per planet), got %d", len(reading.PersonalityGates))
	}
	if len(reading.DesignGates) != 13 {
		t.Errorf("Expected 13 design gates (one per planet), got %d", len(reading.DesignGates))
	}

	// Verify Sun personality gate is in valid range
	sunGate := findGateByPlanet(reading.PersonalityGates, "Sun")
	if sunGate == nil {
		t.Fatal("No Sun personality gate found")
	}
	if sunGate.Number < 1 || sunGate.Number > 64 {
		t.Errorf("Sun gate %d out of range 1-64", sunGate.Number)
	}
	if sunGate.Line < 1 || sunGate.Line > 6 {
		t.Errorf("Sun line %d out of range 1-6", sunGate.Line)
	}

	t.Logf("Type: %s", reading.Type)
	t.Logf("Authority: %s", reading.Authority)
	t.Logf("Profile: %d/%d - %s", reading.Profile.Conscious, reading.Profile.Unconscious, reading.Profile.Name)
	t.Logf("Incarnation Cross: %s", reading.IncarnationCross)
}

// TestGateLineRange checks all gates and lines are in valid ranges.
func TestGateLineRange(t *testing.T) {
	dates := []time.Time{
		time.Date(1990, 6, 15, 14, 30, 0, 0, time.UTC),
		time.Date(1985, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2000, 12, 31, 23, 59, 0, 0, time.UTC),
		time.Date(1970, 7, 4, 12, 0, 0, 0, time.UTC),
	}

	calc := calculator.NewCalculator()
	for _, dt := range dates {
		t.Run(dt.Format("2006-01-02"), func(t *testing.T) {
			data := calculator.BirthData{
				DateTime:  dt,
				Latitude:  40.7128,
				Longitude: -74.0060,
				Location:  "New York",
			}
			reading, err := calc.Calculate(data)
			if err != nil {
				t.Fatalf("Calculate failed: %v", err)
			}
			for _, g := range append(reading.PersonalityGates, reading.DesignGates...) {
				if g.Number < 1 || g.Number > 64 {
					t.Errorf("%s: gate %d out of range", g.Planet, g.Number)
				}
				if g.Line < 1 || g.Line > 6 {
					t.Errorf("%s: line %d out of range for gate %d", g.Planet, g.Line, g.Number)
				}
				if g.Longitude < 0 || g.Longitude >= 360 {
					t.Errorf("%s: longitude %.4f out of range", g.Planet, g.Longitude)
				}
			}
		})
	}
}

// TestAllCentersInitialized checks that all 9 centers are present in a reading.
func TestAllCentersInitialized(t *testing.T) {
	dt := time.Date(1992, 9, 16, 17, 30, 0, 0, time.UTC)
	calc := calculator.NewCalculator()
	reading, err := calc.Calculate(calculator.BirthData{
		DateTime: dt, Latitude: 42.7654, Longitude: -71.4676, Location: "Nashua",
	})
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{"Head", "Ajna", "Throat", "G", "Heart", "Sacral", "SolarPlexus", "Spleen", "Root"}
	for _, name := range expected {
		if _, ok := reading.Centers[name]; !ok {
			t.Errorf("Center %q missing from reading", name)
		}
	}
}

// TestDatabaseCRUD verifies basic database operations.
func TestDatabaseCRUD(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "test_people.json")

	db, err := database.NewDatabase(tmpFile)
	if err != nil {
		t.Fatalf("NewDatabase: %v", err)
	}

	// Add
	p := database.Person{
		Name:      "Test Person",
		BirthDate: "1992-09-16",
		BirthTime: "13:30",
		Location:  "Nashua, US-NH",
		Latitude:  42.7654,
		Longitude: -71.4676,
	}
	if err := db.Add(p); err != nil {
		t.Fatalf("Add: %v", err)
	}

	// List
	people := db.List()
	if len(people) != 1 {
		t.Fatalf("List: expected 1, got %d", len(people))
	}
	if people[0].Name != "Test Person" {
		t.Errorf("Name: expected 'Test Person', got %q", people[0].Name)
	}

	// Get
	got, err := db.Get(people[0].ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.BirthDate != "1992-09-16" {
		t.Errorf("BirthDate: expected '1992-09-16', got %q", got.BirthDate)
	}

	// Delete
	if err := db.Delete(people[0].ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if len(db.List()) != 0 {
		t.Error("Delete: list should be empty")
	}

	// Verify persisted to file
	data, _ := os.ReadFile(tmpFile)
	if string(data) != "[]" {
		t.Errorf("File should contain empty array, got %q", string(data))
	}
}

// TestDatabasePersistence verifies data survives re-opening.
func TestDatabasePersistence(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "test_people.json")

	db1, _ := database.NewDatabase(tmpFile)
	db1.Add(database.Person{Name: "Alice", BirthDate: "2000-01-01", BirthTime: "12:00"})
	db1.Add(database.Person{Name: "Bob", BirthDate: "1995-06-15", BirthTime: "08:00"})

	// Re-open
	db2, err := database.NewDatabase(tmpFile)
	if err != nil {
		t.Fatal(err)
	}
	if len(db2.List()) != 2 {
		t.Fatalf("Expected 2 people after re-open, got %d", len(db2.List()))
	}
}

func findGateByPlanet(gates []calculator.Gate, planet string) *calculator.Gate {
	for _, g := range gates {
		if g.Planet == planet {
			return &g
		}
	}
	return nil
}

func printPositions(t *testing.T, reading *calculator.Reading) {
	t.Helper()
	type entry struct {
		Planet string  `json:"planet"`
		Gate   int     `json:"gate"`
		Line   int     `json:"line"`
		Long   float64 `json:"longitude"`
	}
	toEntries := func(gates []calculator.Gate) []entry {
		out := make([]entry, len(gates))
		for i, g := range gates {
			out[i] = entry{g.Planet, g.Number, g.Line, g.Longitude}
		}
		return out
	}
	payload := map[string]interface{}{
		"type":              reading.Type,
		"authority":         reading.Authority,
		"profile":           fmt.Sprintf("%d/%d", reading.Profile.Conscious, reading.Profile.Unconscious),
		"personality_gates": toEntries(reading.PersonalityGates),
		"design_gates":      toEntries(reading.DesignGates),
	}
	b, _ := json.MarshalIndent(payload, "", "  ")
	t.Logf("Full chart positions:\n%s", string(b))
}
