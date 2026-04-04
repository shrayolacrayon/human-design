package csvreader

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestReadBirthDataCSV(t *testing.T) {
	// Create a temporary CSV file for testing
	tmpDir := t.TempDir()
	csvFile := filepath.Join(tmpDir, "test.csv")

	csvContent := `name,datetime,latitude,longitude,location,expected_type,expected_authority,expected_profile_conscious,expected_profile_unconscious,expected_strategy
Test Person,1990-06-15T14:30:00Z,40.7128,-74.0060,"New York, NY",Generator,Sacral,1,3,Wait to Respond`

	if err := os.WriteFile(csvFile, []byte(csvContent), 0644); err != nil {
		t.Fatalf("Failed to create test CSV: %v", err)
	}

	testCases, err := ReadBirthDataCSV(csvFile)
	if err != nil {
		t.Fatalf("ReadBirthDataCSV failed: %v", err)
	}

	if len(testCases) != 1 {
		t.Fatalf("Expected 1 test case, got %d", len(testCases))
	}

	tc := testCases[0]

	if tc.Name != "Test Person" {
		t.Errorf("Expected name 'Test Person', got '%s'", tc.Name)
	}

	expectedTime := time.Date(1990, 6, 15, 14, 30, 0, 0, time.UTC)
	if !tc.BirthData.DateTime.Equal(expectedTime) {
		t.Errorf("Expected datetime %v, got %v", expectedTime, tc.BirthData.DateTime)
	}

	if tc.BirthData.Latitude != 40.7128 {
		t.Errorf("Expected latitude 40.7128, got %f", tc.BirthData.Latitude)
	}

	if tc.BirthData.Longitude != -74.0060 {
		t.Errorf("Expected longitude -74.0060, got %f", tc.BirthData.Longitude)
	}

	if tc.BirthData.Location != "New York, NY" {
		t.Errorf("Expected location 'New York, NY', got '%s'", tc.BirthData.Location)
	}

	if tc.ExpectedType != "Generator" {
		t.Errorf("Expected type 'Generator', got '%s'", tc.ExpectedType)
	}

	if tc.ExpectedAuthority != "Sacral" {
		t.Errorf("Expected authority 'Sacral', got '%s'", tc.ExpectedAuthority)
	}

	if tc.ExpectedProfileConscious != 1 {
		t.Errorf("Expected conscious profile 1, got %d", tc.ExpectedProfileConscious)
	}

	if tc.ExpectedProfileUnconscious != 3 {
		t.Errorf("Expected unconscious profile 3, got %d", tc.ExpectedProfileUnconscious)
	}

	if tc.ExpectedStrategy != "Wait to Respond" {
		t.Errorf("Expected strategy 'Wait to Respond', got '%s'", tc.ExpectedStrategy)
	}
}

func TestReadBirthDataCSVWithGates(t *testing.T) {
	tmpDir := t.TempDir()
	csvFile := filepath.Join(tmpDir, "test_gates.csv")

	csvContent := `name,datetime,latitude,longitude,location,expected_gates,expected_channels,notes
Gate Test,1990-06-15T14:30:00Z,40.7128,-74.0060,"New York, NY","1,8,64,47","Inspiration,Abstraction",Testing gates and channels`

	if err := os.WriteFile(csvFile, []byte(csvContent), 0644); err != nil {
		t.Fatalf("Failed to create test CSV: %v", err)
	}

	testCases, err := ReadBirthDataCSV(csvFile)
	if err != nil {
		t.Fatalf("ReadBirthDataCSV failed: %v", err)
	}

	if len(testCases) != 1 {
		t.Fatalf("Expected 1 test case, got %d", len(testCases))
	}

	tc := testCases[0]

	expectedGates := []int{1, 8, 64, 47}
	if len(tc.ExpectedGates) != len(expectedGates) {
		t.Errorf("Expected %d gates, got %d", len(expectedGates), len(tc.ExpectedGates))
	}

	for i, gate := range expectedGates {
		if i >= len(tc.ExpectedGates) || tc.ExpectedGates[i] != gate {
			t.Errorf("Expected gate %d at index %d", gate, i)
		}
	}

	expectedChannels := []string{"Inspiration", "Abstraction"}
	if len(tc.ExpectedChannels) != len(expectedChannels) {
		t.Errorf("Expected %d channels, got %d", len(expectedChannels), len(tc.ExpectedChannels))
	}

	for i, channel := range expectedChannels {
		if i >= len(tc.ExpectedChannels) || tc.ExpectedChannels[i] != channel {
			t.Errorf("Expected channel %s at index %d", channel, i)
		}
	}

	if tc.Notes != "Testing gates and channels" {
		t.Errorf("Expected notes 'Testing gates and channels', got '%s'", tc.Notes)
	}
}

func TestWriteBirthDataCSV(t *testing.T) {
	tmpDir := t.TempDir()
	csvFile := filepath.Join(tmpDir, "output.csv")

	testCases := []TestCase{
		{
			Name: "Test Person",
			BirthData: struct {
				DateTime  time.Time `json:"datetime"`
				Latitude  float64   `json:"latitude"`
				Longitude float64   `json:"longitude"`
				Location  string    `json:"location"`
			}{
				DateTime:  time.Date(1990, 6, 15, 14, 30, 0, 0, time.UTC),
				Latitude:  40.7128,
				Longitude: -74.0060,
				Location:  "New York, NY",
			},
			ExpectedType:               "Generator",
			ExpectedAuthority:          "Sacral",
			ExpectedProfileConscious:   1,
			ExpectedProfileUnconscious: 3,
			ExpectedStrategy:           "Wait to Respond",
			ExpectedGates:              []int{1, 8},
			ExpectedChannels:           []string{"Inspiration"},
			Notes:                      "Test case",
		},
	}

	err := WriteBirthDataCSV(csvFile, testCases)
	if err != nil {
		t.Fatalf("WriteBirthDataCSV failed: %v", err)
	}

	// Read back and verify
	readCases, err := ReadBirthDataCSV(csvFile)
	if err != nil {
		t.Fatalf("Failed to read written CSV: %v", err)
	}

	if len(readCases) != 1 {
		t.Fatalf("Expected 1 test case, got %d", len(readCases))
	}

	if readCases[0].Name != "Test Person" {
		t.Errorf("Expected name 'Test Person', got '%s'", readCases[0].Name)
	}
}
