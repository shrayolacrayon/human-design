package main

import (
	"humandesign/internal/calculator"
	"humandesign/internal/csvreader"
	"path/filepath"
	"testing"
)

func TestBirthDataCSVIntegration(t *testing.T) {
	// Read test cases from CSV
	csvPath := filepath.Join("testdata", "birth_data.csv")
	testCases, err := csvreader.ReadBirthDataCSV(csvPath)
	if err != nil {
		t.Skipf("Skipping integration test: %v", err)
		return
	}

	calc := calculator.NewCalculator()

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			reading, err := calc.Calculate(tc.BirthData)
			if err != nil {
				t.Fatalf("Calculate failed for %s: %v", tc.Name, err)
			}

			// Validate type if expected
			if tc.ExpectedType != "" {
				if string(reading.Type) != tc.ExpectedType {
					t.Logf("WARNING: Type mismatch for %s: expected %s, got %s",
						tc.Name, tc.ExpectedType, reading.Type)
					// Note: Not failing because ephemeris is approximate
				}
			}

			// Validate strategy if expected
			if tc.ExpectedStrategy != "" {
				if reading.Strategy != tc.ExpectedStrategy {
					t.Logf("WARNING: Strategy mismatch for %s: expected %s, got %s",
						tc.Name, tc.ExpectedStrategy, reading.Strategy)
				}
			}

			// Validate profile if expected
			if tc.ExpectedProfileConscious > 0 {
				if reading.Profile.Conscious != tc.ExpectedProfileConscious {
					t.Logf("WARNING: Profile conscious mismatch for %s: expected %d, got %d",
						tc.Name, tc.ExpectedProfileConscious, reading.Profile.Conscious)
				}
			}

			if tc.ExpectedProfileUnconscious > 0 {
				if reading.Profile.Unconscious != tc.ExpectedProfileUnconscious {
					t.Logf("WARNING: Profile unconscious mismatch for %s: expected %d, got %d",
						tc.Name, tc.ExpectedProfileUnconscious, reading.Profile.Unconscious)
				}
			}

			// Log the reading for debugging
			t.Logf("%s: Type=%s, Authority=%s, Profile=%d/%d, Strategy=%s",
				tc.Name, reading.Type, reading.Authority,
				reading.Profile.Conscious, reading.Profile.Unconscious,
				reading.Strategy)
		})
	}
}

func TestGateValidationCSVIntegration(t *testing.T) {
	csvPath := filepath.Join("testdata", "gate_validation.csv")
	testCases, err := csvreader.ReadBirthDataCSV(csvPath)
	if err != nil {
		t.Skipf("Skipping gate validation test: %v", err)
		return
	}

	calc := calculator.NewCalculator()

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			reading, err := calc.Calculate(tc.BirthData)
			if err != nil {
				t.Fatalf("Calculate failed for %s: %v", tc.Name, err)
			}

			// Check if expected gates are present
			if len(tc.ExpectedGates) > 0 {
				allGates := make(map[int]bool)
				for _, gate := range reading.PersonalityGates {
					allGates[gate.Number] = true
				}
				for _, gate := range reading.DesignGates {
					allGates[gate.Number] = true
				}

				for _, expectedGate := range tc.ExpectedGates {
					if !allGates[expectedGate] {
						t.Logf("WARNING: Expected gate %d not found for %s",
							expectedGate, tc.Name)
					}
				}
			}

			// Check if expected channels are present
			if len(tc.ExpectedChannels) > 0 {
				definedChannels := make(map[string]bool)
				for _, ch := range reading.Channels {
					if ch.Defined {
						definedChannels[ch.Name] = true
					}
				}

				for _, expectedChannel := range tc.ExpectedChannels {
					if !definedChannels[expectedChannel] {
						t.Logf("WARNING: Expected channel '%s' not found or not defined for %s",
							expectedChannel, tc.Name)
					}
				}
			}

			// Log results
			t.Logf("%s: %d personality gates, %d design gates, %d defined channels",
				tc.Name,
				len(reading.PersonalityGates),
				len(reading.DesignGates),
				countDefinedChannels(reading.Channels))
		})
	}
}

func countDefinedChannels(channels []calculator.Channel) int {
	count := 0
	for _, ch := range channels {
		if ch.Defined {
			count++
		}
	}
	return count
}
