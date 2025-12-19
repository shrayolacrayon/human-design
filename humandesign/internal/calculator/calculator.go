package calculator

import (
	"fmt"
	"humandesign/internal/ephemeris"
)

// Calculator performs Human Design calculations
type Calculator struct {
	ephemeris *ephemeris.Ephemeris
}

// NewCalculator creates a new Human Design calculator
func NewCalculator() *Calculator {
	return &Calculator{
		ephemeris: ephemeris.NewEphemeris(),
	}
}

// Calculate generates a full Human Design reading from birth data
func (c *Calculator) Calculate(data BirthData) (*Reading, error) {
	reading := &Reading{
		BirthData: data,
		Centers:   make(map[string]*Center),
	}

	// Initialize all centers as undefined
	for _, name := range CenterNames {
		reading.Centers[name] = &Center{
			Name:    name,
			Defined: false,
			Gates:   []int{},
		}
	}

	// Calculate personality (conscious) positions - at birth time
	personalityPositions := c.ephemeris.CalculatePositions(data.DateTime)
	
	// Calculate design (unconscious) positions - 88 days before birth
	designPositions := c.ephemeris.CalculateDesignPositions(data.DateTime)

	// Convert positions to gates
	reading.PersonalityGates = c.positionsToGates(personalityPositions, false)
	reading.DesignGates = c.positionsToGates(designPositions, true)

	// Activate gates in centers
	allGates := append(reading.PersonalityGates, reading.DesignGates...)
	for _, gate := range allGates {
		centerName := GetGateCenter(gate.Number)
		if center, ok := reading.Centers[centerName]; ok {
			if !containsGate(center.Gates, gate.Number) {
				center.Gates = append(center.Gates, gate.Number)
			}
		}
	}

	// Find defined channels
	reading.Channels = c.findDefinedChannels(allGates)

	// Define centers based on channels
	for _, ch := range reading.Channels {
		if ch.Defined {
			if center, ok := reading.Centers[ch.Center1]; ok {
				center.Defined = true
			}
			if center, ok := reading.Centers[ch.Center2]; ok {
				center.Defined = true
			}
		}
	}

	// Determine Type
	reading.Type = c.determineType(reading.Centers)

	// Determine Strategy based on Type
	reading.Strategy = c.determineStrategy(reading.Type)

	// Determine Authority
	reading.Authority = c.determineAuthority(reading.Centers, reading.Type)

	// Determine Profile from Sun gates
	reading.Profile = c.determineProfile(personalityPositions, designPositions)

	// Determine Definition type
	reading.Definition = c.determineDefinition(reading.Centers, reading.Channels)

	// Set Not-Self Theme and Signature
	reading.NotSelfTheme, reading.Signature = c.getThemeAndSignature(reading.Type)

	// Calculate Incarnation Cross
	reading.IncarnationCross = c.calculateIncarnationCross(personalityPositions, designPositions)

	return reading, nil
}

func (c *Calculator) positionsToGates(positions []ephemeris.PlanetaryPosition, isDesign bool) []Gate {
	gates := []Gate{}
	for _, pos := range positions {
		gates = append(gates, Gate{
			Number:    pos.Gate,
			Name:      AllGates[pos.Gate].Name,
			Line:      pos.Line,
			Planet:    string(pos.Planet),
			Activated: true,
			Design:    isDesign,
		})
	}
	return gates
}

func (c *Calculator) findDefinedChannels(gates []Gate) []Channel {
	channels := []Channel{}
	activatedGates := make(map[int]bool)
	
	for _, gate := range gates {
		activatedGates[gate.Number] = true
	}

	for _, chDef := range AllChannels {
		isDefined := activatedGates[chDef.Gate1] && activatedGates[chDef.Gate2]
		channels = append(channels, Channel{
			Name:    chDef.Name,
			Gate1:   chDef.Gate1,
			Gate2:   chDef.Gate2,
			Center1: chDef.Center1,
			Center2: chDef.Center2,
			Defined: isDefined,
		})
	}

	return channels
}

func (c *Calculator) determineType(centers map[string]*Center) HumanDesignType {
	sacralDefined := centers["Sacral"].Defined
	throatDefined := centers["Throat"].Defined
	
	// Check for motor to throat connection
	motorToThroat := c.hasMotorToThroatConnection(centers)

	// Reflector: No centers defined (very rare)
	definedCount := 0
	for _, center := range centers {
		if center.Defined {
			definedCount++
		}
	}
	if definedCount == 0 {
		return TypeReflector
	}

	// Generator types have defined Sacral
	if sacralDefined {
		// Manifesting Generator has motor connected to throat
		if motorToThroat && throatDefined {
			return TypeManifestingGenerator
		}
		return TypeGenerator
	}

	// Manifestor: Motor to Throat but no Sacral
	if motorToThroat && !sacralDefined {
		return TypeManifestor
	}

	// Projector: No Sacral, no motor to throat
	return TypeProjector
}

func (c *Calculator) hasMotorToThroatConnection(centers map[string]*Center) bool {
	// Motors are: Sacral, Solar Plexus, Heart, Root
	motors := []string{"Sacral", "SolarPlexus", "Heart", "Root"}
	
	if !centers["Throat"].Defined {
		return false
	}

	// Simplified check - in reality need to trace through channels
	for _, motor := range motors {
		if centers[motor].Defined {
			// Check if there's a path from motor to throat
			// This is simplified - real implementation would trace channels
			return true
		}
	}
	return false
}

func (c *Calculator) determineStrategy(hdType HumanDesignType) string {
	switch hdType {
	case TypeGenerator, TypeManifestingGenerator:
		return "Wait to Respond"
	case TypeProjector:
		return "Wait for the Invitation"
	case TypeManifestor:
		return "Inform Before Acting"
	case TypeReflector:
		return "Wait a Lunar Cycle"
	default:
		return "Unknown"
	}
}

func (c *Calculator) determineAuthority(centers map[string]*Center, hdType HumanDesignType) Authority {
	// Authority hierarchy (in order of precedence)
	
	// 1. Emotional Authority (Solar Plexus defined)
	if centers["SolarPlexus"].Defined {
		return AuthorityEmotional
	}

	// 2. Sacral Authority (Sacral defined, no Solar Plexus)
	if centers["Sacral"].Defined {
		return AuthoritySacral
	}

	// 3. Splenic Authority (Spleen defined)
	if centers["Spleen"].Defined {
		return AuthoritySplenic
	}

	// 4. Ego/Heart Authority (Heart defined and connected to Throat)
	if centers["Heart"].Defined {
		return AuthorityEgo
	}

	// 5. Self-Projected (G Center connected to Throat)
	if centers["G"].Defined && centers["Throat"].Defined {
		return AuthoritySelf
	}

	// 6. Environmental/Mental (Projector with no inner authority)
	if hdType == TypeProjector {
		return AuthorityEnvironmental
	}

	// 7. Lunar (Reflector)
	if hdType == TypeReflector {
		return AuthorityLunar
	}

	return AuthorityEnvironmental
}

func (c *Calculator) determineProfile(personality, design []ephemeris.PlanetaryPosition) Profile {
	var conscioousLine, unconsciousLine int

	// Find Sun positions
	for _, pos := range personality {
		if pos.Planet == ephemeris.Sun {
			conscioousLine = pos.Line
			break
		}
	}

	for _, pos := range design {
		if pos.Planet == ephemeris.Sun {
			unconsciousLine = pos.Line
			break
		}
	}

	profileKey := fmt.Sprintf("%d/%d", conscioousLine, unconsciousLine)
	name := ProfileNames[profileKey]
	if name == "" {
		name = "Unknown Profile"
	}

	return Profile{
		Conscious:   conscioousLine,
		Unconscious: unconsciousLine,
		Name:        name,
	}
}

func (c *Calculator) determineDefinition(centers map[string]*Center, channels []Channel) string {
	definedCenters := []string{}
	for name, center := range centers {
		if center.Defined {
			definedCenters = append(definedCenters, name)
		}
	}

	if len(definedCenters) == 0 {
		return "No Definition"
	}

	// Count connected groups of centers
	// This is simplified - real implementation would use graph traversal
	definedCount := len(definedCenters)
	
	if definedCount <= 2 {
		return "Single Definition"
	} else if definedCount <= 4 {
		return "Split Definition"
	} else if definedCount <= 6 {
		return "Triple Split Definition"
	} else {
		return "Quadruple Split Definition"
	}
}

func (c *Calculator) getThemeAndSignature(hdType HumanDesignType) (string, string) {
	switch hdType {
	case TypeGenerator, TypeManifestingGenerator:
		return "Frustration", "Satisfaction"
	case TypeProjector:
		return "Bitterness", "Success"
	case TypeManifestor:
		return "Anger", "Peace"
	case TypeReflector:
		return "Disappointment", "Surprise"
	default:
		return "Unknown", "Unknown"
	}
}

func (c *Calculator) calculateIncarnationCross(personality, design []ephemeris.PlanetaryPosition) string {
	var sunP, earthP, sunD, earthD int

	for _, pos := range personality {
		switch pos.Planet {
		case ephemeris.Sun:
			sunP = pos.Gate
		case ephemeris.Earth:
			earthP = pos.Gate
		}
	}

	for _, pos := range design {
		switch pos.Planet {
		case ephemeris.Sun:
			sunD = pos.Gate
		case ephemeris.Earth:
			earthD = pos.Gate
		}
	}

	// Format: Right Angle Cross of... / Left Angle Cross of... / Juxtaposition Cross of...
	// Simplified - just return the gates
	return fmt.Sprintf("Cross of Gates %d/%d | %d/%d", sunP, earthP, sunD, earthD)
}

func containsGate(gates []int, gate int) bool {
	for _, g := range gates {
		if g == gate {
			return true
		}
	}
	return false
}
