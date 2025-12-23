package main

import (
	"fmt"
	"humandesign/internal/calculator"
	"humandesign/internal/ephemeris"
	"time"
)

func main() {
	// Shraya's birth data
	birthTime := time.Date(1992, 9, 16, 13, 30, 0, 0, time.UTC)

	eph := ephemeris.NewEphemeris()

	// Calculate personality (birth time)
	fmt.Println("=== PERSONALITY (Conscious) ===")
	personalityPos := eph.CalculatePositions(birthTime)

	for _, pos := range personalityPos {
		gateName := ""
		if info, ok := calculator.AllGates[pos.Gate]; ok {
			gateName = info.Name
		}
		fmt.Printf("%-12s: %7.3f° → Gate %2d.%d - %s\n",
			pos.Planet, pos.Longitude, pos.Gate, pos.Line, gateName)
	}

	fmt.Println("\n=== DESIGN (Unconscious) ===")
	designPos := eph.CalculateDesignPositions(birthTime)

	for _, pos := range designPos {
		gateName := ""
		if info, ok := calculator.AllGates[pos.Gate]; ok {
			gateName = info.Name
		}
		fmt.Printf("%-12s: %7.3f° → Gate %2d.%d - %s\n",
			pos.Planet, pos.Longitude, pos.Gate, pos.Line, gateName)
	}
}
