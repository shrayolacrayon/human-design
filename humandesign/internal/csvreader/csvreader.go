package csvreader

import (
	"encoding/csv"
	"fmt"
	"humandesign/internal/calculator"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

// TestCase represents a test case from CSV
type TestCase struct {
	Name                      string
	BirthData                 calculator.BirthData
	ExpectedType              string
	ExpectedAuthority         string
	ExpectedProfileConscious  int
	ExpectedProfileUnconscious int
	ExpectedStrategy          string
	ExpectedGates             []int
	ExpectedChannels          []string
	Notes                     string
}

// ReadBirthDataCSV reads birth data test cases from a CSV file
func ReadBirthDataCSV(filename string) ([]TestCase, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	// Create column index map
	colIndex := make(map[string]int)
	for i, col := range header {
		colIndex[col] = i
	}

	var testCases []TestCase

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read record: %w", err)
		}

		tc, err := parseTestCase(record, colIndex)
		if err != nil {
			return nil, fmt.Errorf("failed to parse test case: %w", err)
		}

		testCases = append(testCases, tc)
	}

	return testCases, nil
}

// parseTestCase parses a CSV record into a TestCase
func parseTestCase(record []string, colIndex map[string]int) (TestCase, error) {
	tc := TestCase{}

	// Parse name
	if idx, ok := colIndex["name"]; ok && idx < len(record) {
		tc.Name = record[idx]
	}

	// Parse datetime
	if idx, ok := colIndex["datetime"]; ok && idx < len(record) {
		datetime, err := time.Parse(time.RFC3339, record[idx])
		if err != nil {
			return tc, fmt.Errorf("failed to parse datetime: %w", err)
		}
		tc.BirthData.DateTime = datetime
	}

	// Parse latitude
	if idx, ok := colIndex["latitude"]; ok && idx < len(record) {
		lat, err := strconv.ParseFloat(record[idx], 64)
		if err != nil {
			return tc, fmt.Errorf("failed to parse latitude: %w", err)
		}
		tc.BirthData.Latitude = lat
	}

	// Parse longitude
	if idx, ok := colIndex["longitude"]; ok && idx < len(record) {
		lon, err := strconv.ParseFloat(record[idx], 64)
		if err != nil {
			return tc, fmt.Errorf("failed to parse longitude: %w", err)
		}
		tc.BirthData.Longitude = lon
	}

	// Parse location
	if idx, ok := colIndex["location"]; ok && idx < len(record) {
		tc.BirthData.Location = record[idx]
	}

	// Parse expected type
	if idx, ok := colIndex["expected_type"]; ok && idx < len(record) {
		tc.ExpectedType = record[idx]
	}

	// Parse expected authority
	if idx, ok := colIndex["expected_authority"]; ok && idx < len(record) {
		tc.ExpectedAuthority = record[idx]
	}

	// Parse expected profile conscious
	if idx, ok := colIndex["expected_profile_conscious"]; ok && idx < len(record) {
		if record[idx] != "" {
			conscious, err := strconv.Atoi(record[idx])
			if err != nil {
				return tc, fmt.Errorf("failed to parse profile conscious: %w", err)
			}
			tc.ExpectedProfileConscious = conscious
		}
	}

	// Parse expected profile unconscious
	if idx, ok := colIndex["expected_profile_unconscious"]; ok && idx < len(record) {
		if record[idx] != "" {
			unconscious, err := strconv.Atoi(record[idx])
			if err != nil {
				return tc, fmt.Errorf("failed to parse profile unconscious: %w", err)
			}
			tc.ExpectedProfileUnconscious = unconscious
		}
	}

	// Parse expected strategy
	if idx, ok := colIndex["expected_strategy"]; ok && idx < len(record) {
		tc.ExpectedStrategy = record[idx]
	}

	// Parse expected gates (comma-separated list)
	if idx, ok := colIndex["expected_gates"]; ok && idx < len(record) {
		if record[idx] != "" {
			gatesStr := strings.Split(record[idx], ",")
			for _, gateStr := range gatesStr {
				gate, err := strconv.Atoi(strings.TrimSpace(gateStr))
				if err != nil {
					return tc, fmt.Errorf("failed to parse gate: %w", err)
				}
				tc.ExpectedGates = append(tc.ExpectedGates, gate)
			}
		}
	}

	// Parse expected channels (comma-separated list)
	if idx, ok := colIndex["expected_channels"]; ok && idx < len(record) {
		if record[idx] != "" {
			tc.ExpectedChannels = strings.Split(record[idx], ",")
			for i := range tc.ExpectedChannels {
				tc.ExpectedChannels[i] = strings.TrimSpace(tc.ExpectedChannels[i])
			}
		}
	}

	// Parse notes
	if idx, ok := colIndex["notes"]; ok && idx < len(record) {
		tc.Notes = record[idx]
	}

	return tc, nil
}

// WriteBirthDataCSV writes birth data test cases to a CSV file
func WriteBirthDataCSV(filename string, testCases []TestCase) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"name", "datetime", "latitude", "longitude", "location",
		"expected_type", "expected_authority", "expected_profile_conscious",
		"expected_profile_unconscious", "expected_strategy", "expected_gates",
		"expected_channels", "notes",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write records
	for _, tc := range testCases {
		record := []string{
			tc.Name,
			tc.BirthData.DateTime.Format(time.RFC3339),
			fmt.Sprintf("%f", tc.BirthData.Latitude),
			fmt.Sprintf("%f", tc.BirthData.Longitude),
			tc.BirthData.Location,
			tc.ExpectedType,
			tc.ExpectedAuthority,
			fmt.Sprintf("%d", tc.ExpectedProfileConscious),
			fmt.Sprintf("%d", tc.ExpectedProfileUnconscious),
			tc.ExpectedStrategy,
			joinInts(tc.ExpectedGates),
			strings.Join(tc.ExpectedChannels, ","),
			tc.Notes,
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	return nil
}

// joinInts joins a slice of ints into a comma-separated string
func joinInts(ints []int) string {
	strs := make([]string, len(ints))
	for i, v := range ints {
		strs[i] = strconv.Itoa(v)
	}
	return strings.Join(strs, ",")
}
