package calculator

import (
	"testing"
)

// TestContainsGate tests the containsGate helper function
func TestContainsGate(t *testing.T) {
	tests := []struct {
		name     string
		gates    []int
		gate     int
		expected bool
	}{
		{
			name:     "Gate present in middle",
			gates:    []int{1, 5, 10, 15, 20},
			gate:     10,
			expected: true,
		},
		{
			name:     "Gate present at start",
			gates:    []int{1, 5, 10, 15, 20},
			gate:     1,
			expected: true,
		},
		{
			name:     "Gate present at end",
			gates:    []int{1, 5, 10, 15, 20},
			gate:     20,
			expected: true,
		},
		{
			name:     "Gate not present",
			gates:    []int{1, 5, 10, 15, 20},
			gate:     7,
			expected: false,
		},
		{
			name:     "Empty slice",
			gates:    []int{},
			gate:     5,
			expected: false,
		},
		{
			name:     "Single element - match",
			gates:    []int{42},
			gate:     42,
			expected: true,
		},
		{
			name:     "Single element - no match",
			gates:    []int{42},
			gate:     10,
			expected: false,
		},
		{
			name:     "Duplicate gates",
			gates:    []int{5, 5, 5},
			gate:     5,
			expected: true,
		},
		{
			name:     "Zero gate in slice",
			gates:    []int{0, 1, 2},
			gate:     0,
			expected: true,
		},
		{
			name:     "Negative gate",
			gates:    []int{1, 2, 3},
			gate:     -1,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := containsGate(tt.gates, tt.gate)
			if got != tt.expected {
				t.Errorf("containsGate(%v, %d) = %v, want %v", tt.gates, tt.gate, got, tt.expected)
			}
		})
	}
}

// TestDetermineStrategy tests the strategy determination logic
func TestDetermineStrategy(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name         string
		hdType       HumanDesignType
		wantStrategy string
	}{
		{
			name:         "Generator strategy",
			hdType:       TypeGenerator,
			wantStrategy: "Wait to Respond",
		},
		{
			name:         "Manifesting Generator strategy",
			hdType:       TypeManifestingGenerator,
			wantStrategy: "Wait to Respond",
		},
		{
			name:         "Projector strategy",
			hdType:       TypeProjector,
			wantStrategy: "Wait for the Invitation",
		},
		{
			name:         "Manifestor strategy",
			hdType:       TypeManifestor,
			wantStrategy: "Inform Before Acting",
		},
		{
			name:         "Reflector strategy",
			hdType:       TypeReflector,
			wantStrategy: "Wait a Lunar Cycle",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calc.determineStrategy(tt.hdType)
			if got != tt.wantStrategy {
				t.Errorf("determineStrategy(%v) = %v, want %v", tt.hdType, got, tt.wantStrategy)
			}
		})
	}
}

// TestGetThemeAndSignature tests the theme and signature determination
func TestGetThemeAndSignature(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name          string
		hdType        HumanDesignType
		wantNotSelf   string
		wantSignature string
	}{
		{
			name:          "Generator",
			hdType:        TypeGenerator,
			wantNotSelf:   "Frustration",
			wantSignature: "Satisfaction",
		},
		{
			name:          "Manifesting Generator",
			hdType:        TypeManifestingGenerator,
			wantNotSelf:   "Frustration",
			wantSignature: "Satisfaction",
		},
		{
			name:          "Projector",
			hdType:        TypeProjector,
			wantNotSelf:   "Bitterness",
			wantSignature: "Success",
		},
		{
			name:          "Manifestor",
			hdType:        TypeManifestor,
			wantNotSelf:   "Anger",
			wantSignature: "Peace",
		},
		{
			name:          "Reflector",
			hdType:        TypeReflector,
			wantNotSelf:   "Disappointment",
			wantSignature: "Surprise",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNotSelf, gotSignature := calc.getThemeAndSignature(tt.hdType)
			if gotNotSelf != tt.wantNotSelf {
				t.Errorf("getThemeAndSignature(%v) NotSelf = %v, want %v", tt.hdType, gotNotSelf, tt.wantNotSelf)
			}
			if gotSignature != tt.wantSignature {
				t.Errorf("getThemeAndSignature(%v) Signature = %v, want %v", tt.hdType, gotSignature, tt.wantSignature)
			}
		})
	}
}

// TestDetermineType tests the type determination logic
func TestDetermineType(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name             string
		definedCenters   map[string]bool
		wantType         HumanDesignType
	}{
		{
			name: "Reflector - no centers defined",
			definedCenters: map[string]bool{
				"Head":        false,
				"Ajna":        false,
				"Throat":      false,
				"G":           false,
				"Heart":       false,
				"Sacral":      false,
				"SolarPlexus": false,
				"Spleen":      false,
				"Root":        false,
			},
			wantType: TypeReflector,
		},
		{
			name: "Generator - Sacral defined only",
			definedCenters: map[string]bool{
				"Head":        false,
				"Ajna":        false,
				"Throat":      false,
				"G":           false,
				"Heart":       false,
				"Sacral":      true,
				"SolarPlexus": false,
				"Spleen":      false,
				"Root":        false,
			},
			wantType: TypeGenerator,
		},
		{
			name: "Manifesting Generator - Sacral and Throat defined",
			definedCenters: map[string]bool{
				"Head":        false,
				"Ajna":        false,
				"Throat":      true,
				"G":           false,
				"Heart":       true, // Motor
				"Sacral":      true,
				"SolarPlexus": false,
				"Spleen":      false,
				"Root":        false,
			},
			wantType: TypeManifestingGenerator,
		},
		{
			name: "Projector - No Sacral, no motor to throat",
			definedCenters: map[string]bool{
				"Head":        true,
				"Ajna":        true,
				"Throat":      true,
				"G":           false,
				"Heart":       false,
				"Sacral":      false,
				"SolarPlexus": false,
				"Spleen":      true,
				"Root":        false,
			},
			wantType: TypeProjector,
		},
		{
			name: "Manifestor - Motor to Throat, no Sacral",
			definedCenters: map[string]bool{
				"Head":        false,
				"Ajna":        false,
				"Throat":      true,
				"G":           false,
				"Heart":       true, // Motor
				"Sacral":      false,
				"SolarPlexus": false,
				"Spleen":      false,
				"Root":        false,
			},
			wantType: TypeManifestor,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			centers := make(map[string]*Center)
			for name, defined := range tt.definedCenters {
				centers[name] = &Center{
					Name:    name,
					Defined: defined,
					Gates:   []int{},
				}
			}

			got := calc.determineType(centers)
			if got != tt.wantType {
				t.Errorf("determineType() = %v, want %v", got, tt.wantType)
			}
		})
	}
}

// TestDetermineAuthority tests the authority determination logic
func TestDetermineAuthority(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name           string
		definedCenters map[string]bool
		hdType         HumanDesignType
		wantAuthority  Authority
	}{
		{
			name: "Emotional Authority - Solar Plexus defined",
			definedCenters: map[string]bool{
				"SolarPlexus": true,
				"Sacral":      true,
				"Spleen":      false,
				"Heart":       false,
				"G":           false,
				"Throat":      false,
			},
			hdType:        TypeGenerator,
			wantAuthority: AuthorityEmotional,
		},
		{
			name: "Sacral Authority - Sacral defined, no Solar Plexus",
			definedCenters: map[string]bool{
				"SolarPlexus": false,
				"Sacral":      true,
				"Spleen":      false,
				"Heart":       false,
				"G":           false,
				"Throat":      false,
			},
			hdType:        TypeGenerator,
			wantAuthority: AuthoritySacral,
		},
		{
			name: "Splenic Authority - Spleen defined, no Solar Plexus or Sacral",
			definedCenters: map[string]bool{
				"SolarPlexus": false,
				"Sacral":      false,
				"Spleen":      true,
				"Heart":       false,
				"G":           false,
				"Throat":      false,
			},
			hdType:        TypeProjector,
			wantAuthority: AuthoritySplenic,
		},
		{
			name: "Ego Authority - Heart defined, no Solar Plexus/Sacral/Spleen",
			definedCenters: map[string]bool{
				"SolarPlexus": false,
				"Sacral":      false,
				"Spleen":      false,
				"Heart":       true,
				"G":           false,
				"Throat":      false,
			},
			hdType:        TypeManifestor,
			wantAuthority: AuthorityEgo,
		},
		{
			name: "Self-Projected Authority - G and Throat defined",
			definedCenters: map[string]bool{
				"SolarPlexus": false,
				"Sacral":      false,
				"Spleen":      false,
				"Heart":       false,
				"G":           true,
				"Throat":      true,
			},
			hdType:        TypeProjector,
			wantAuthority: AuthoritySelf,
		},
		{
			name: "Environmental Authority - Projector with no inner authority",
			definedCenters: map[string]bool{
				"SolarPlexus": false,
				"Sacral":      false,
				"Spleen":      false,
				"Heart":       false,
				"G":           false,
				"Throat":      true,
			},
			hdType:        TypeProjector,
			wantAuthority: AuthorityEnvironmental,
		},
		{
			name: "Lunar Authority - Reflector",
			definedCenters: map[string]bool{
				"SolarPlexus": false,
				"Sacral":      false,
				"Spleen":      false,
				"Heart":       false,
				"G":           false,
				"Throat":      false,
			},
			hdType:        TypeReflector,
			wantAuthority: AuthorityLunar,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			centers := make(map[string]*Center)
			for name, defined := range tt.definedCenters {
				centers[name] = &Center{
					Name:    name,
					Defined: defined,
					Gates:   []int{},
				}
			}

			got := calc.determineAuthority(centers, tt.hdType)
			if got != tt.wantAuthority {
				t.Errorf("determineAuthority() = %v, want %v", got, tt.wantAuthority)
			}
		})
	}
}

// TestDetermineDefinition tests the definition determination logic
func TestDetermineDefinition(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name             string
		definedCenters   []string
		wantDefinition   string
	}{
		{
			name:           "No Definition",
			definedCenters: []string{},
			wantDefinition: "No Definition",
		},
		{
			name:           "Single Definition - 1 center",
			definedCenters: []string{"Sacral"},
			wantDefinition: "Single Definition",
		},
		{
			name:           "Single Definition - 2 centers",
			definedCenters: []string{"Sacral", "Throat"},
			wantDefinition: "Single Definition",
		},
		{
			name:           "Split Definition - 3 centers",
			definedCenters: []string{"Sacral", "Throat", "Ajna"},
			wantDefinition: "Split Definition",
		},
		{
			name:           "Split Definition - 4 centers",
			definedCenters: []string{"Sacral", "Throat", "Ajna", "G"},
			wantDefinition: "Split Definition",
		},
		{
			name:           "Triple Split - 5 centers",
			definedCenters: []string{"Sacral", "Throat", "Ajna", "G", "Root"},
			wantDefinition: "Triple Split Definition",
		},
		{
			name:           "Triple Split - 6 centers",
			definedCenters: []string{"Sacral", "Throat", "Ajna", "G", "Root", "Spleen"},
			wantDefinition: "Triple Split Definition",
		},
		{
			name:           "Quadruple Split - 7 centers",
			definedCenters: []string{"Sacral", "Throat", "Ajna", "G", "Root", "Spleen", "SolarPlexus"},
			wantDefinition: "Quadruple Split Definition",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			centers := make(map[string]*Center)
			// Initialize all centers as undefined
			for _, name := range CenterNames {
				centers[name] = &Center{
					Name:    name,
					Defined: false,
					Gates:   []int{},
				}
			}
			// Define the specified centers
			for _, name := range tt.definedCenters {
				if center, ok := centers[name]; ok {
					center.Defined = true
				}
			}

			got := calc.determineDefinition(centers, []Channel{})
			if got != tt.wantDefinition {
				t.Errorf("determineDefinition() = %v, want %v", got, tt.wantDefinition)
			}
		})
	}
}

// TestProfileNames ensures all profile combinations are defined
func TestProfileNames(t *testing.T) {
	// All valid profile combinations in Human Design
	validProfiles := []string{
		"1/3", "1/4",
		"2/4", "2/5",
		"3/5", "3/6",
		"4/6", "4/1",
		"5/1", "5/2",
		"6/2", "6/3",
	}

	for _, profile := range validProfiles {
		name, exists := ProfileNames[profile]
		if !exists {
			t.Errorf("Profile %s is not defined in ProfileNames", profile)
		}
		if name == "" {
			t.Errorf("Profile %s has empty name", profile)
		}
	}

	// Ensure we have exactly 12 profiles
	if len(ProfileNames) != 12 {
		t.Errorf("Expected 12 profiles, but found %d", len(ProfileNames))
	}
}

// TestCenterNames ensures all 9 centers are present
func TestCenterNames(t *testing.T) {
	expectedCenters := []string{
		"Head", "Ajna", "Throat", "G", "Heart",
		"Sacral", "SolarPlexus", "Spleen", "Root",
	}

	if len(CenterNames) != 9 {
		t.Errorf("Expected 9 centers, but found %d", len(CenterNames))
	}

	centerMap := make(map[string]bool)
	for _, center := range CenterNames {
		centerMap[center] = true
	}

	for _, expected := range expectedCenters {
		if !centerMap[expected] {
			t.Errorf("Center %s is missing from CenterNames", expected)
		}
	}
}

// TestFindDefinedChannels tests the channel finding logic
func TestFindDefinedChannels(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name            string
		gates           []Gate
		wantDefinedCh   []string // Names of channels that should be defined
		wantUndefinedCh []string // Names of channels that should not be defined
	}{
		{
			name: "Channel 64-47 defined",
			gates: []Gate{
				{Number: 64, Name: "Confusion", Activated: true},
				{Number: 47, Name: "Realization", Activated: true},
			},
			wantDefinedCh:   []string{"Abstraction"},
			wantUndefinedCh: []string{"Logic", "Awareness"},
		},
		{
			name: "No channels defined",
			gates: []Gate{
				{Number: 64, Name: "Confusion", Activated: true},
				{Number: 24, Name: "Rationalization", Activated: true},
			},
			wantDefinedCh:   []string{},
			wantUndefinedCh: []string{"Abstraction", "Logic", "Awareness"},
		},
		{
			name: "Multiple channels defined",
			gates: []Gate{
				{Number: 64, Name: "Confusion", Activated: true},
				{Number: 47, Name: "Realization", Activated: true},
				{Number: 63, Name: "Doubt", Activated: true},
				{Number: 4, Name: "Formulization", Activated: true},
			},
			wantDefinedCh:   []string{"Abstraction", "Logic"},
			wantUndefinedCh: []string{"Awareness"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			channels := calc.findDefinedChannels(tt.gates)

			// Check defined channels
			for _, wantCh := range tt.wantDefinedCh {
				found := false
				for _, ch := range channels {
					if ch.Name == wantCh && ch.Defined {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected channel %s to be defined, but it wasn't", wantCh)
				}
			}

			// Check undefined channels
			for _, wantCh := range tt.wantUndefinedCh {
				for _, ch := range channels {
					if ch.Name == wantCh && ch.Defined {
						t.Errorf("Expected channel %s to be undefined, but it was defined", wantCh)
					}
				}
			}
		})
	}
}

// TestNewCalculator ensures calculator is initialized properly
func TestNewCalculator(t *testing.T) {
	calc := NewCalculator()
	if calc == nil {
		t.Error("NewCalculator() returned nil")
	}
	if calc.ephemeris == nil {
		t.Error("NewCalculator() did not initialize ephemeris")
	}
}
