package calculator

import "testing"

func TestAllGatesExist(t *testing.T) {
	// Test that all 64 gates are defined
	if len(AllGates) != 64 {
		t.Errorf("Expected 64 gates, got %d", len(AllGates))
	}

	// Test each gate from 1 to 64
	for i := 1; i <= 64; i++ {
		gate, exists := AllGates[i]
		if !exists {
			t.Errorf("Gate %d is not defined", i)
			continue
		}

		if gate.Number != i {
			t.Errorf("Gate %d has incorrect number: %d", i, gate.Number)
		}

		if gate.Name == "" {
			t.Errorf("Gate %d has no name", i)
		}

		if gate.Center == "" {
			t.Errorf("Gate %d has no center assignment", i)
		}
	}
}

func TestGateCenterAssignments(t *testing.T) {
	validCenters := map[string]bool{
		"Head":        true,
		"Ajna":        true,
		"Throat":      true,
		"G":           true,
		"Heart":       true,
		"Sacral":      true,
		"SolarPlexus": true,
		"Spleen":      true,
		"Root":        true,
	}

	for gateNum, gate := range AllGates {
		if !validCenters[gate.Center] {
			t.Errorf("Gate %d has invalid center: %s", gateNum, gate.Center)
		}
	}
}

func TestGetGateCenter(t *testing.T) {
	tests := []struct {
		gate   int
		center string
	}{
		{1, "G"},
		{64, "Head"},
		{3, "Sacral"},
		{6, "SolarPlexus"},
		{18, "Spleen"},
		{21, "Heart"},
		{8, "Throat"},
		{4, "Ajna"},
		{19, "Root"},
	}

	for _, test := range tests {
		result := GetGateCenter(test.gate)
		if result != test.center {
			t.Errorf("GetGateCenter(%d) = %s, expected %s", test.gate, result, test.center)
		}
	}
}

func TestChannelDefinitions(t *testing.T) {
	// Test that channels have valid gates
	for _, channel := range AllChannels {
		if channel.Gate1 < 1 || channel.Gate1 > 64 {
			t.Errorf("Channel %s has invalid Gate1: %d", channel.Name, channel.Gate1)
		}
		if channel.Gate2 < 1 || channel.Gate2 > 64 {
			t.Errorf("Channel %s has invalid Gate2: %d", channel.Name, channel.Gate2)
		}

		// Test that gates exist
		if _, ok := AllGates[channel.Gate1]; !ok {
			t.Errorf("Channel %s references non-existent Gate1: %d", channel.Name, channel.Gate1)
		}
		if _, ok := AllGates[channel.Gate2]; !ok {
			t.Errorf("Channel %s references non-existent Gate2: %d", channel.Name, channel.Gate2)
		}

		// Test that channel has a name
		if channel.Name == "" {
			t.Errorf("Channel %d-%d has no name", channel.Gate1, channel.Gate2)
		}
	}
}

func TestGetChannelForGates(t *testing.T) {
	tests := []struct {
		gate1       int
		gate2       int
		shouldExist bool
		name        string
	}{
		{64, 47, true, "Abstraction"},
		{47, 64, true, "Abstraction"}, // Test reversed order
		{1, 8, true, "Inspiration"},
		{21, 45, true, "Money"},
		{1, 2, false, ""}, // These don't form a channel
	}

	for _, test := range tests {
		channel := GetChannelForGates(test.gate1, test.gate2)
		if test.shouldExist {
			if channel == nil {
				t.Errorf("Expected channel for gates %d-%d, got nil", test.gate1, test.gate2)
			} else if channel.Name != test.name {
				t.Errorf("Expected channel name %s for gates %d-%d, got %s",
					test.name, test.gate1, test.gate2, channel.Name)
			}
		} else {
			if channel != nil {
				t.Errorf("Expected no channel for gates %d-%d, got %s", test.gate1, test.gate2, channel.Name)
			}
		}
	}
}

func TestChannelCenterConnections(t *testing.T) {
	validCenters := map[string]bool{
		"Head":        true,
		"Ajna":        true,
		"Throat":      true,
		"G":           true,
		"Heart":       true,
		"Sacral":      true,
		"SolarPlexus": true,
		"Spleen":      true,
		"Root":        true,
	}

	for _, channel := range AllChannels {
		// Check that centers are valid
		if !validCenters[channel.Center1] {
			t.Errorf("Channel %s has invalid Center1: %s", channel.Name, channel.Center1)
		}
		if !validCenters[channel.Center2] {
			t.Errorf("Channel %s has invalid Center2: %s", channel.Name, channel.Center2)
		}

		// Check that gates belong to their respective centers
		gate1Center := GetGateCenter(channel.Gate1)
		gate2Center := GetGateCenter(channel.Gate2)

		// Gates should belong to one of the two centers in the channel
		gate1Valid := gate1Center == channel.Center1 || gate1Center == channel.Center2
		gate2Valid := gate2Center == channel.Center1 || gate2Center == channel.Center2

		if !gate1Valid {
			t.Errorf("Channel %s: Gate %d (center: %s) doesn't belong to centers %s or %s",
				channel.Name, channel.Gate1, gate1Center, channel.Center1, channel.Center2)
		}
		if !gate2Valid {
			t.Errorf("Channel %s: Gate %d (center: %s) doesn't belong to centers %s or %s",
				channel.Name, channel.Gate2, gate2Center, channel.Center1, channel.Center2)
		}
	}
}

func TestProfileNames(t *testing.T) {
	expectedProfiles := []string{
		"1/3", "1/4", "2/4", "2/5", "3/5", "3/6",
		"4/6", "4/1", "5/1", "5/2", "6/2", "6/3",
	}

	if len(ProfileNames) != len(expectedProfiles) {
		t.Errorf("Expected %d profiles, got %d", len(expectedProfiles), len(ProfileNames))
	}

	for _, profile := range expectedProfiles {
		if name, ok := ProfileNames[profile]; !ok {
			t.Errorf("Profile %s is not defined", profile)
		} else if name == "" {
			t.Errorf("Profile %s has empty name", profile)
		}
	}
}

func TestCenterNames(t *testing.T) {
	if len(CenterNames) != 9 {
		t.Errorf("Expected 9 centers, got %d", len(CenterNames))
	}

	expectedCenters := map[string]bool{
		"Head":        true,
		"Ajna":        true,
		"Throat":      true,
		"G":           true,
		"Heart":       true,
		"Sacral":      true,
		"SolarPlexus": true,
		"Spleen":      true,
		"Root":        true,
	}

	for _, center := range CenterNames {
		if !expectedCenters[center] {
			t.Errorf("Unexpected center: %s", center)
		}
	}
}
