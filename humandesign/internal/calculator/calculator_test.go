package calculator

import (
	"testing"
	"time"
)

func TestNewCalculator(t *testing.T) {
	calc := NewCalculator()
	if calc == nil {
		t.Fatal("NewCalculator returned nil")
	}
	if calc.ephemeris == nil {
		t.Fatal("Calculator ephemeris is nil")
	}
}

func TestDetermineType(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name           string
		sacralDefined  bool
		throatDefined  bool
		motorDefined   bool
		expectedType   HumanDesignType
	}{
		{
			name:          "Generator - Sacral defined, no motor to throat",
			sacralDefined: true,
			throatDefined: false,
			motorDefined:  false,
			expectedType:  TypeGenerator,
		},
		{
			name:          "Manifesting Generator - Sacral and throat defined with motor",
			sacralDefined: true,
			throatDefined: true,
			motorDefined:  true,
			expectedType:  TypeManifestingGenerator,
		},
		{
			name:          "Projector - No sacral, no motor to throat, but other centers defined",
			sacralDefined: false,
			throatDefined: false,
			motorDefined:  false,
			expectedType:  TypeProjector,
		},
		{
			name:          "Manifestor - Motor to throat, no sacral",
			sacralDefined: false,
			throatDefined: true,
			motorDefined:  true,
			expectedType:  TypeManifestor,
		},
		{
			name:          "Reflector - No centers defined",
			sacralDefined: false,
			throatDefined: false,
			motorDefined:  false,
			expectedType:  TypeReflector,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			centers := make(map[string]*Center)
			for _, name := range CenterNames {
				centers[name] = &Center{
					Name:    name,
					Defined: false,
					Gates:   []int{},
				}
			}

			centers["Sacral"].Defined = test.sacralDefined
			centers["Throat"].Defined = test.throatDefined
			if test.motorDefined {
				centers["Heart"].Defined = true
			}

			// For Projector test, define some non-motor centers (Ajna, G)
			// to distinguish from Reflector (which has NO centers defined)
			if test.name == "Projector - No sacral, no motor to throat, but other centers defined" {
				centers["Ajna"].Defined = true
				centers["G"].Defined = true
			}

			result := calc.determineType(centers)
			if result != test.expectedType {
				t.Errorf("Expected type %s, got %s", test.expectedType, result)
			}
		})
	}
}

func TestDetermineStrategy(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		hdType           HumanDesignType
		expectedStrategy string
	}{
		{TypeGenerator, "Wait to Respond"},
		{TypeManifestingGenerator, "Wait to Respond"},
		{TypeProjector, "Wait for the Invitation"},
		{TypeManifestor, "Inform Before Acting"},
		{TypeReflector, "Wait a Lunar Cycle"},
	}

	for _, test := range tests {
		result := calc.determineStrategy(test.hdType)
		if result != test.expectedStrategy {
			t.Errorf("For type %s, expected strategy %s, got %s",
				test.hdType, test.expectedStrategy, result)
		}
	}
}

func TestDetermineAuthority(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name              string
		solarPlexusDefined bool
		sacralDefined     bool
		spleenDefined     bool
		heartDefined      bool
		gDefined          bool
		throatDefined     bool
		hdType            HumanDesignType
		expectedAuthority Authority
	}{
		{
			name:               "Emotional Authority - Solar Plexus defined",
			solarPlexusDefined: true,
			expectedAuthority:  AuthorityEmotional,
		},
		{
			name:              "Sacral Authority - Sacral defined, no Solar Plexus",
			solarPlexusDefined: false,
			sacralDefined:     true,
			expectedAuthority: AuthoritySacral,
		},
		{
			name:              "Splenic Authority - Spleen defined",
			solarPlexusDefined: false,
			sacralDefined:     false,
			spleenDefined:     true,
			expectedAuthority: AuthoritySplenic,
		},
		{
			name:              "Ego Authority - Heart defined",
			solarPlexusDefined: false,
			sacralDefined:     false,
			spleenDefined:     false,
			heartDefined:      true,
			expectedAuthority: AuthorityEgo,
		},
		{
			name:              "Self-Projected Authority - G and Throat defined",
			solarPlexusDefined: false,
			sacralDefined:     false,
			spleenDefined:     false,
			heartDefined:      false,
			gDefined:          true,
			throatDefined:     true,
			expectedAuthority: AuthoritySelf,
		},
		{
			name:              "Lunar Authority - Reflector",
			solarPlexusDefined: false,
			sacralDefined:     false,
			spleenDefined:     false,
			heartDefined:      false,
			gDefined:          false,
			throatDefined:     false,
			hdType:            TypeReflector,
			expectedAuthority: AuthorityLunar,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			centers := make(map[string]*Center)
			for _, name := range CenterNames {
				centers[name] = &Center{
					Name:    name,
					Defined: false,
					Gates:   []int{},
				}
			}

			centers["SolarPlexus"].Defined = test.solarPlexusDefined
			centers["Sacral"].Defined = test.sacralDefined
			centers["Spleen"].Defined = test.spleenDefined
			centers["Heart"].Defined = test.heartDefined
			centers["G"].Defined = test.gDefined
			centers["Throat"].Defined = test.throatDefined

			result := calc.determineAuthority(centers, test.hdType)
			if result != test.expectedAuthority {
				t.Errorf("Expected authority %s, got %s", test.expectedAuthority, result)
			}
		})
	}
}

func TestGetThemeAndSignature(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		hdType            HumanDesignType
		expectedTheme     string
		expectedSignature string
	}{
		{TypeGenerator, "Frustration", "Satisfaction"},
		{TypeManifestingGenerator, "Frustration", "Satisfaction"},
		{TypeProjector, "Bitterness", "Success"},
		{TypeManifestor, "Anger", "Peace"},
		{TypeReflector, "Disappointment", "Surprise"},
	}

	for _, test := range tests {
		theme, signature := calc.getThemeAndSignature(test.hdType)
		if theme != test.expectedTheme {
			t.Errorf("For type %s, expected theme %s, got %s",
				test.hdType, test.expectedTheme, theme)
		}
		if signature != test.expectedSignature {
			t.Errorf("For type %s, expected signature %s, got %s",
				test.hdType, test.expectedSignature, signature)
		}
	}
}

func TestContainsGate(t *testing.T) {
	gates := []int{1, 8, 13, 25, 64}

	tests := []struct {
		gate     int
		expected bool
	}{
		{1, true},
		{8, true},
		{13, true},
		{25, true},
		{64, true},
		{2, false},
		{10, false},
		{100, false},
	}

	for _, test := range tests {
		result := containsGate(gates, test.gate)
		if result != test.expected {
			t.Errorf("containsGate(%v, %d) = %v, expected %v",
				gates, test.gate, result, test.expected)
		}
	}
}

func TestFindDefinedChannels(t *testing.T) {
	calc := NewCalculator()

	// Create gates that form the 64-47 channel (Head-Ajna Abstraction)
	gates := []Gate{
		{Number: 64, Activated: true},
		{Number: 47, Activated: true},
	}

	channels := calc.findDefinedChannels(gates)

	// Find the Abstraction channel
	var abstractionChannel *Channel
	for i := range channels {
		if channels[i].Name == "Abstraction" {
			abstractionChannel = &channels[i]
			break
		}
	}

	if abstractionChannel == nil {
		t.Fatal("Abstraction channel not found in results")
	}

	if !abstractionChannel.Defined {
		t.Error("Abstraction channel should be defined when gates 64 and 47 are activated")
	}
}

func TestCalculate(t *testing.T) {
	calc := NewCalculator()

	birthData := BirthData{
		DateTime:  time.Date(1990, 6, 15, 14, 30, 0, 0, time.UTC),
		Latitude:  40.7128,
		Longitude: -74.0060,
		Location:  "New York, NY",
	}

	reading, err := calc.Calculate(birthData)
	if err != nil {
		t.Fatalf("Calculate failed: %v", err)
	}

	// Test that basic fields are populated
	if reading == nil {
		t.Fatal("Reading is nil")
	}

	if reading.Type == "" {
		t.Error("Type is empty")
	}

	if reading.Authority == "" {
		t.Error("Authority is empty")
	}

	if reading.Strategy == "" {
		t.Error("Strategy is empty")
	}

	if reading.Profile.Conscious < 1 || reading.Profile.Conscious > 6 {
		t.Errorf("Invalid conscious profile: %d", reading.Profile.Conscious)
	}

	if reading.Profile.Unconscious < 1 || reading.Profile.Unconscious > 6 {
		t.Errorf("Invalid unconscious profile: %d", reading.Profile.Unconscious)
	}

	if len(reading.Centers) != 9 {
		t.Errorf("Expected 9 centers, got %d", len(reading.Centers))
	}

	if len(reading.PersonalityGates) == 0 {
		t.Error("No personality gates calculated")
	}

	if len(reading.DesignGates) == 0 {
		t.Error("No design gates calculated")
	}
}

func TestDetermineDefinition(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name               string
		definedCenterCount int
		expectedDefinition string
	}{
		{"No Definition", 0, "No Definition"},
		{"Single Definition", 2, "Single Definition"},
		{"Split Definition", 4, "Split Definition"},
		{"Triple Split", 6, "Triple Split Definition"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			centers := make(map[string]*Center)
			for i, name := range CenterNames {
				centers[name] = &Center{
					Name:    name,
					Defined: i < test.definedCenterCount,
					Gates:   []int{},
				}
			}

			result := calc.determineDefinition(centers, []Channel{})
			if result != test.expectedDefinition {
				t.Errorf("Expected %s, got %s", test.expectedDefinition, result)
			}
		})
	}
}
