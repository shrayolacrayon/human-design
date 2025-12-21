package calculator

import (
	"testing"
)

// TestGetGateCenter tests the GetGateCenter function
func TestGetGateCenter(t *testing.T) {
	tests := []struct {
		name       string
		gateNum    int
		wantCenter string
	}{
		// Test Head center gates
		{name: "Gate 61 - Head", gateNum: 61, wantCenter: "Head"},
		{name: "Gate 63 - Head", gateNum: 63, wantCenter: "Head"},
		{name: "Gate 64 - Head", gateNum: 64, wantCenter: "Head"},

		// Test Ajna center gates
		{name: "Gate 4 - Ajna", gateNum: 4, wantCenter: "Ajna"},
		{name: "Gate 11 - Ajna", gateNum: 11, wantCenter: "Ajna"},
		{name: "Gate 17 - Ajna", gateNum: 17, wantCenter: "Ajna"},
		{name: "Gate 24 - Ajna", gateNum: 24, wantCenter: "Ajna"},
		{name: "Gate 43 - Ajna", gateNum: 43, wantCenter: "Ajna"},
		{name: "Gate 47 - Ajna", gateNum: 47, wantCenter: "Ajna"},

		// Test Throat center gates
		{name: "Gate 8 - Throat", gateNum: 8, wantCenter: "Throat"},
		{name: "Gate 12 - Throat", gateNum: 12, wantCenter: "Throat"},
		{name: "Gate 16 - Throat", gateNum: 16, wantCenter: "Throat"},
		{name: "Gate 20 - Throat", gateNum: 20, wantCenter: "Throat"},
		{name: "Gate 23 - Throat", gateNum: 23, wantCenter: "Throat"},
		{name: "Gate 31 - Throat", gateNum: 31, wantCenter: "Throat"},
		{name: "Gate 33 - Throat", gateNum: 33, wantCenter: "Throat"},
		{name: "Gate 35 - Throat", gateNum: 35, wantCenter: "Throat"},
		{name: "Gate 45 - Throat", gateNum: 45, wantCenter: "Throat"},
		{name: "Gate 56 - Throat", gateNum: 56, wantCenter: "Throat"},
		{name: "Gate 62 - Throat", gateNum: 62, wantCenter: "Throat"},

		// Test G center gates
		{name: "Gate 1 - G", gateNum: 1, wantCenter: "G"},
		{name: "Gate 2 - G", gateNum: 2, wantCenter: "G"},
		{name: "Gate 7 - G", gateNum: 7, wantCenter: "G"},
		{name: "Gate 10 - G", gateNum: 10, wantCenter: "G"},
		{name: "Gate 13 - G", gateNum: 13, wantCenter: "G"},
		{name: "Gate 15 - G", gateNum: 15, wantCenter: "G"},
		{name: "Gate 25 - G", gateNum: 25, wantCenter: "G"},
		{name: "Gate 46 - G", gateNum: 46, wantCenter: "G"},

		// Test Heart/Ego center gates
		{name: "Gate 21 - Heart", gateNum: 21, wantCenter: "Heart"},
		{name: "Gate 26 - Heart", gateNum: 26, wantCenter: "Heart"},
		{name: "Gate 40 - Heart", gateNum: 40, wantCenter: "Heart"},
		{name: "Gate 51 - Heart", gateNum: 51, wantCenter: "Heart"},

		// Test Sacral center gates
		{name: "Gate 3 - Sacral", gateNum: 3, wantCenter: "Sacral"},
		{name: "Gate 5 - Sacral", gateNum: 5, wantCenter: "Sacral"},
		{name: "Gate 9 - Sacral", gateNum: 9, wantCenter: "Sacral"},
		{name: "Gate 14 - Sacral", gateNum: 14, wantCenter: "Sacral"},
		{name: "Gate 27 - Sacral", gateNum: 27, wantCenter: "Sacral"},
		{name: "Gate 29 - Sacral", gateNum: 29, wantCenter: "Sacral"},
		{name: "Gate 34 - Sacral", gateNum: 34, wantCenter: "Sacral"},
		{name: "Gate 42 - Sacral", gateNum: 42, wantCenter: "Sacral"},
		{name: "Gate 59 - Sacral", gateNum: 59, wantCenter: "Sacral"},

		// Test Solar Plexus center gates
		{name: "Gate 6 - SolarPlexus", gateNum: 6, wantCenter: "SolarPlexus"},
		{name: "Gate 22 - SolarPlexus", gateNum: 22, wantCenter: "SolarPlexus"},
		{name: "Gate 30 - SolarPlexus", gateNum: 30, wantCenter: "SolarPlexus"},
		{name: "Gate 36 - SolarPlexus", gateNum: 36, wantCenter: "SolarPlexus"},
		{name: "Gate 37 - SolarPlexus", gateNum: 37, wantCenter: "SolarPlexus"},
		{name: "Gate 49 - SolarPlexus", gateNum: 49, wantCenter: "SolarPlexus"},
		{name: "Gate 55 - SolarPlexus", gateNum: 55, wantCenter: "SolarPlexus"},

		// Test Spleen center gates
		{name: "Gate 18 - Spleen", gateNum: 18, wantCenter: "Spleen"},
		{name: "Gate 28 - Spleen", gateNum: 28, wantCenter: "Spleen"},
		{name: "Gate 32 - Spleen", gateNum: 32, wantCenter: "Spleen"},
		{name: "Gate 44 - Spleen", gateNum: 44, wantCenter: "Spleen"},
		{name: "Gate 48 - Spleen", gateNum: 48, wantCenter: "Spleen"},
		{name: "Gate 50 - Spleen", gateNum: 50, wantCenter: "Spleen"},
		{name: "Gate 57 - Spleen", gateNum: 57, wantCenter: "Spleen"},

		// Test Root center gates
		{name: "Gate 19 - Root", gateNum: 19, wantCenter: "Root"},
		{name: "Gate 38 - Root", gateNum: 38, wantCenter: "Root"},
		{name: "Gate 39 - Root", gateNum: 39, wantCenter: "Root"},
		{name: "Gate 41 - Root", gateNum: 41, wantCenter: "Root"},
		{name: "Gate 52 - Root", gateNum: 52, wantCenter: "Root"},
		{name: "Gate 53 - Root", gateNum: 53, wantCenter: "Root"},
		{name: "Gate 54 - Root", gateNum: 54, wantCenter: "Root"},
		{name: "Gate 58 - Root", gateNum: 58, wantCenter: "Root"},
		{name: "Gate 60 - Root", gateNum: 60, wantCenter: "Root"},

		// Test edge cases
		{name: "Invalid gate 0", gateNum: 0, wantCenter: ""},
		{name: "Invalid gate -1", gateNum: -1, wantCenter: ""},
		{name: "Invalid gate 65", gateNum: 65, wantCenter: ""},
		{name: "Invalid gate 100", gateNum: 100, wantCenter: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetGateCenter(tt.gateNum)
			if got != tt.wantCenter {
				t.Errorf("GetGateCenter(%d) = %v, want %v", tt.gateNum, got, tt.wantCenter)
			}
		})
	}
}

// TestGetGateCenterAllGates ensures all 64 gates are mapped correctly
func TestGetGateCenterAllGates(t *testing.T) {
	for i := 1; i <= 64; i++ {
		center := GetGateCenter(i)
		if center == "" {
			t.Errorf("Gate %d has no center mapping", i)
		}
	}
}

// TestGetChannelForGates tests the GetChannelForGates function
func TestGetChannelForGates(t *testing.T) {
	tests := []struct {
		name         string
		gate1        int
		gate2        int
		wantChannel  bool
		wantName     string
		wantCenter1  string
		wantCenter2  string
	}{
		// Test valid channels - Head to Ajna
		{
			name:        "Channel 64-47 - Abstraction",
			gate1:       64,
			gate2:       47,
			wantChannel: true,
			wantName:    "Abstraction",
			wantCenter1: "Head",
			wantCenter2: "Ajna",
		},
		{
			name:        "Channel 63-4 - Logic",
			gate1:       63,
			gate2:       4,
			wantChannel: true,
			wantName:    "Logic",
			wantCenter1: "Head",
			wantCenter2: "Ajna",
		},
		{
			name:        "Channel 61-24 - Awareness",
			gate1:       61,
			gate2:       24,
			wantChannel: true,
			wantName:    "Awareness",
			wantCenter1: "Head",
			wantCenter2: "Ajna",
		},

		// Test Ajna to Throat channels
		{
			name:        "Channel 43-23 - Structuring",
			gate1:       43,
			gate2:       23,
			wantChannel: true,
			wantName:    "Structuring",
			wantCenter1: "Ajna",
			wantCenter2: "Throat",
		},
		{
			name:        "Channel 11-56 - Curiosity",
			gate1:       11,
			gate2:       56,
			wantChannel: true,
			wantName:    "Curiosity",
			wantCenter1: "Ajna",
			wantCenter2: "Throat",
		},

		// Test G to Throat channels
		{
			name:        "Channel 7-31 - The Alpha",
			gate1:       7,
			gate2:       31,
			wantChannel: true,
			wantName:    "The Alpha",
			wantCenter1: "G",
			wantCenter2: "Throat",
		},
		{
			name:        "Channel 1-8 - Inspiration",
			gate1:       1,
			gate2:       8,
			wantChannel: true,
			wantName:    "Inspiration",
			wantCenter1: "G",
			wantCenter2: "Throat",
		},

		// Test Sacral channels
		{
			name:        "Channel 3-60 - Mutation",
			gate1:       3,
			gate2:       60,
			wantChannel: true,
			wantName:    "Mutation",
			wantCenter1: "Sacral",
			wantCenter2: "Root",
		},
		{
			name:        "Channel 5-15 - Rhythm",
			gate1:       5,
			gate2:       15,
			wantChannel: true,
			wantName:    "Rhythm",
			wantCenter1: "G",
			wantCenter2: "Sacral",
		},

		// Test order independence (gate2, gate1 should work same as gate1, gate2)
		{
			name:        "Channel 64-47 - Reversed order",
			gate1:       47,
			gate2:       64,
			wantChannel: true,
			wantName:    "Abstraction",
			wantCenter1: "Head",
			wantCenter2: "Ajna",
		},
		{
			name:        "Channel 7-31 - Reversed order",
			gate1:       31,
			gate2:       7,
			wantChannel: true,
			wantName:    "The Alpha",
			wantCenter1: "G",
			wantCenter2: "Throat",
		},

		// Test invalid channel combinations
		{
			name:        "Invalid - Gates 1 and 2",
			gate1:       1,
			gate2:       2,
			wantChannel: false,
		},
		{
			name:        "Invalid - Gates 5 and 6",
			gate1:       5,
			gate2:       6,
			wantChannel: false,
		},
		{
			name:        "Invalid - Same gate",
			gate1:       10,
			gate2:       10,
			wantChannel: false,
		},

		// Test edge cases
		{
			name:        "Invalid - Gate 0 and 1",
			gate1:       0,
			gate2:       1,
			wantChannel: false,
		},
		{
			name:        "Invalid - Gate 64 and 65",
			gate1:       64,
			gate2:       65,
			wantChannel: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetChannelForGates(tt.gate1, tt.gate2)
			if tt.wantChannel {
				if got == nil {
					t.Errorf("GetChannelForGates(%d, %d) returned nil, expected channel", tt.gate1, tt.gate2)
					return
				}
				if got.Name != tt.wantName {
					t.Errorf("GetChannelForGates(%d, %d) Name = %v, want %v", tt.gate1, tt.gate2, got.Name, tt.wantName)
				}
				// Check that centers match (either order)
				centersMatch := (got.Center1 == tt.wantCenter1 && got.Center2 == tt.wantCenter2) ||
					(got.Center1 == tt.wantCenter2 && got.Center2 == tt.wantCenter1)
				if !centersMatch {
					t.Errorf("GetChannelForGates(%d, %d) Centers = (%v, %v), want (%v, %v)",
						tt.gate1, tt.gate2, got.Center1, got.Center2, tt.wantCenter1, tt.wantCenter2)
				}
			} else {
				if got != nil {
					t.Errorf("GetChannelForGates(%d, %d) = %v, want nil", tt.gate1, tt.gate2, got.Name)
				}
			}
		})
	}
}

// TestAllChannelsValid ensures all defined channels have valid gates
func TestAllChannelsValid(t *testing.T) {
	for _, ch := range AllChannels {
		// Test that both gates exist
		center1 := GetGateCenter(ch.Gate1)
		center2 := GetGateCenter(ch.Gate2)

		if center1 == "" {
			t.Errorf("Channel %s has invalid Gate1: %d", ch.Name, ch.Gate1)
		}
		if center2 == "" {
			t.Errorf("Channel %s has invalid Gate2: %d", ch.Name, ch.Gate2)
		}

		// Test that gates are different
		if ch.Gate1 == ch.Gate2 {
			t.Errorf("Channel %s has same gate for both ends: %d", ch.Name, ch.Gate1)
		}

		// Test that channel can be found
		found := GetChannelForGates(ch.Gate1, ch.Gate2)
		if found == nil {
			t.Errorf("Channel %s (%d-%d) cannot be found by GetChannelForGates", ch.Name, ch.Gate1, ch.Gate2)
		}
	}
}

// TestChannelCount ensures we have all 36 channels
func TestChannelCount(t *testing.T) {
	// Human Design has 36 channels
	expectedChannels := 36
	actualChannels := len(AllChannels)

	if actualChannels != expectedChannels {
		t.Errorf("Expected %d channels, but found %d", expectedChannels, actualChannels)
	}
}

// TestGateCount ensures we have all 64 gates
func TestGateCount(t *testing.T) {
	// Human Design has 64 gates
	expectedGates := 64
	actualGates := len(AllGates)

	if actualGates != expectedGates {
		t.Errorf("Expected %d gates, but found %d", expectedGates, actualGates)
	}
}

// TestGateNames ensures all gates have names
func TestGateNames(t *testing.T) {
	for gateNum, gateInfo := range AllGates {
		if gateInfo.Name == "" {
			t.Errorf("Gate %d has empty name", gateNum)
		}
		if gateInfo.Number != gateNum {
			t.Errorf("Gate %d has mismatched Number field: %d", gateNum, gateInfo.Number)
		}
		if gateInfo.Center == "" {
			t.Errorf("Gate %d has empty center", gateNum)
		}
	}
}

// TestChannelBidirectionality ensures channels can be found in both directions
func TestChannelBidirectionality(t *testing.T) {
	for _, ch := range AllChannels {
		ch1 := GetChannelForGates(ch.Gate1, ch.Gate2)
		ch2 := GetChannelForGates(ch.Gate2, ch.Gate1)

		if ch1 == nil || ch2 == nil {
			t.Errorf("Channel %s cannot be found bidirectionally", ch.Name)
			continue
		}

		if ch1.Name != ch2.Name {
			t.Errorf("Channel names don't match when reversed: %s vs %s", ch1.Name, ch2.Name)
		}
	}
}

// TestCenterNamesValid ensures all centers in gates are valid
func TestCenterNamesValid(t *testing.T) {
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

	for gateNum, gateInfo := range AllGates {
		if !validCenters[gateInfo.Center] {
			t.Errorf("Gate %d has invalid center: %s", gateNum, gateInfo.Center)
		}
	}
}
