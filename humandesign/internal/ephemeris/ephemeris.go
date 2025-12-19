package ephemeris

/*
#cgo CFLAGS: -I/usr/local/include -I/opt/homebrew/include
#cgo LDFLAGS: -L/usr/local/lib -L/opt/homebrew/lib -lswe -lm
#include <swephexp.h>
#include <stdlib.h>
*/
import "C"

import (
	"math"
	"time"
)

// Planet represents a celestial body
type Planet string

const (
	Sun       Planet = "Sun"
	Earth     Planet = "Earth"
	Moon      Planet = "Moon"
	Mercury   Planet = "Mercury"
	Venus     Planet = "Venus"
	Mars      Planet = "Mars"
	Jupiter   Planet = "Jupiter"
	Saturn    Planet = "Saturn"
	Uranus    Planet = "Uranus"
	Neptune   Planet = "Neptune"
	Pluto     Planet = "Pluto"
	NorthNode Planet = "North Node"
	SouthNode Planet = "South Node"
)

// PlanetaryPosition holds the position of a planet
type PlanetaryPosition struct {
	Planet    Planet
	Longitude float64 // 0-360 degrees
	Gate      int     // 1-64
	Line      int     // 1-6
}

// Ephemeris calculates planetary positions using Swiss Ephemeris
type Ephemeris struct{}

// NewEphemeris creates a new ephemeris calculator
func NewEphemeris() *Ephemeris {
	return &Ephemeris{}
}

// getPlanetNumber returns the Swiss Ephemeris planet number
func getPlanetNumber(planet Planet) C.int {
	switch planet {
	case Sun:
		return C.SE_SUN
	case Moon:
		return C.SE_MOON
	case Mercury:
		return C.SE_MERCURY
	case Venus:
		return C.SE_VENUS
	case Mars:
		return C.SE_MARS
	case Jupiter:
		return C.SE_JUPITER
	case Saturn:
		return C.SE_SATURN
	case Uranus:
		return C.SE_URANUS
	case Neptune:
		return C.SE_NEPTUNE
	case Pluto:
		return C.SE_PLUTO
	case NorthNode:
		return C.SE_TRUE_NODE
	default:
		return C.SE_SUN
	}
}

// calculateLongitude calculates the ecliptic longitude for a planet at a given Julian Day
func (e *Ephemeris) calculateLongitude(planet Planet, jd float64) float64 {
	// Handle special cases
	if planet == Earth {
		sunLong := e.calculateLongitude(Sun, jd)
		return normalizeAngle(sunLong + 180.0)
	}
	if planet == SouthNode {
		nodeLong := e.calculateLongitude(NorthNode, jd)
		return normalizeAngle(nodeLong + 180.0)
	}

	var xx [6]C.double
	var serr [256]C.char

	planetNum := getPlanetNumber(planet)
	// Use Moshier ephemeris (SEFLG_MOSEPH) - accurate and no data files needed
	iflag := C.int(C.SEFLG_MOSEPH)

	C.swe_calc(C.double(jd), planetNum, iflag, &xx[0], &serr[0])

	return float64(xx[0])
}

// CalculatePositions calculates planetary positions for a given time
func (e *Ephemeris) CalculatePositions(dt time.Time) []PlanetaryPosition {
	jd := julianDay(dt)
	positions := []PlanetaryPosition{}

	planets := []Planet{Sun, Earth, Moon, Mercury, Venus, Mars, Jupiter, Saturn, Uranus, Neptune, Pluto, NorthNode, SouthNode}

	for _, planet := range planets {
		longitude := e.calculateLongitude(planet, jd)
		gate, line := longitudeToGateLine(longitude)

		positions = append(positions, PlanetaryPosition{
			Planet:    planet,
			Longitude: longitude,
			Gate:      gate,
			Line:      line,
		})
	}

	return positions
}

// CalculateDesignPositions calculates the design (unconscious) positions
// Design is calculated for when the Sun was 88 DEGREES behind its birth position
func (e *Ephemeris) CalculateDesignPositions(dt time.Time) []PlanetaryPosition {
	jd := julianDay(dt)
	birthSunLong := e.calculateLongitude(Sun, jd)

	// Target longitude is 88° behind (earlier in the zodiac)
	targetLong := normalizeAngle(birthSunLong - 88.0)

	// Binary search for the exact Julian Day when Sun was at target longitude
	// Start approximately 89 days before birth
	designJD := jd - 89.0

	// Newton-Raphson iteration to find precise date
	for i := 0; i < 20; i++ {
		sunLong := e.calculateLongitude(Sun, designJD)
		diff := targetLong - sunLong

		// Handle wraparound
		if diff > 180 {
			diff -= 360
		} else if diff < -180 {
			diff += 360
		}

		// Adjust by approximate days needed (Sun moves ~0.9856° per day)
		designJD += diff / 0.9856

		if math.Abs(diff) < 0.0001 {
			break
		}
	}

	// Calculate all planetary positions at the design date
	return e.CalculatePositionsAtJD(designJD)
}

// CalculatePositionsAtJD calculates positions for a specific Julian Day
func (e *Ephemeris) CalculatePositionsAtJD(jd float64) []PlanetaryPosition {
	positions := []PlanetaryPosition{}

	planets := []Planet{Sun, Earth, Moon, Mercury, Venus, Mars, Jupiter, Saturn, Uranus, Neptune, Pluto, NorthNode, SouthNode}

	for _, planet := range planets {
		longitude := e.calculateLongitude(planet, jd)
		gate, line := longitudeToGateLine(longitude)

		positions = append(positions, PlanetaryPosition{
			Planet:    planet,
			Longitude: longitude,
			Gate:      gate,
			Line:      line,
		})
	}

	return positions
}

// julianDay calculates the Julian Day Number for a given time
func julianDay(t time.Time) float64 {
	year := float64(t.Year())
	month := float64(t.Month())
	day := float64(t.Day())
	hour := float64(t.Hour()) + float64(t.Minute())/60.0 + float64(t.Second())/3600.0

	if month <= 2 {
		year--
		month += 12
	}

	A := math.Floor(year / 100)
	B := 2 - A + math.Floor(A/4)

	jd := math.Floor(365.25*(year+4716)) + math.Floor(30.6001*(month+1)) + day + hour/24.0 + B - 1524.5

	return jd
}

func normalizeAngle(angle float64) float64 {
	angle = math.Mod(angle, 360)
	if angle < 0 {
		angle += 360
	}
	return angle
}

// Human Design Gate Wheel - starting from Gate 41 at 2° Aquarius (302° tropical)
// Each gate spans 5.625° (360/64)
var hdGateWheelFrom41 = []int{
	41, 19, 13, 49, 30, 55, 37, 63, 22, 36, // 302° - 358.25°
	25, 17, 21, 51, 42, 3, 27, 24, 2, 23,   // 358.25° - 56.25° (wraps through 0°)
	8, 20, 16, 35, 45, 12, 15, 52, 39, 53,  // 56.25° - 112.5°
	62, 56, 31, 33, 7, 4, 29, 59, 40, 64,   // 112.5° - 168.75°
	47, 6, 46, 18, 48, 57, 32, 50, 28, 44,  // 168.75° - 225°
	1, 43, 14, 34, 9, 5, 26, 11, 10, 58,    // 225° - 281.25°
	38, 54, 61, 60,                          // 281.25° - 302° (back to 41)
}

// HD wheel starts at 302° (2° Aquarius = Gate 41)
const hdWheelStartDegree = 302.0

// longitudeToGateLine converts ecliptic longitude to Human Design gate and line
func longitudeToGateLine(longitude float64) (int, int) {
	// Normalize longitude to 0-360
	longitude = normalizeAngle(longitude)

	// Convert to HD wheel position (offset from 302°)
	hdPosition := longitude - hdWheelStartDegree
	if hdPosition < 0 {
		hdPosition += 360
	}

	// Each gate spans 5.625 degrees (360/64)
	gateIndex := int(hdPosition / 5.625)
	if gateIndex >= 64 {
		gateIndex = 63
	}

	// Get the gate number from the wheel
	gate := hdGateWheelFrom41[gateIndex]

	// Each line spans 0.9375 degrees (5.625/6)
	positionInGate := math.Mod(hdPosition, 5.625)
	line := int(positionInGate/0.9375) + 1
	if line > 6 {
		line = 6
	}

	return gate, line
}
