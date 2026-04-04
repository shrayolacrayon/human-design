package astrocartography

import (
	"humandesign/internal/ephemeris"
	"math"
	"time"
)

// LineType represents the type of astrocartography line
type LineType string

const (
	LineASC LineType = "ASC" // Ascendant line
	LineDSC LineType = "DSC" // Descendant line
	LineMC  LineType = "MC"  // Midheaven line
	LineIC  LineType = "IC"  // Imum Coeli line
)

// PlanetaryLine represents a single astrocartography line
type PlanetaryLine struct {
	Planet    string    `json:"planet"`
	LineType  LineType  `json:"line_type"`
	Points    []GeoPoint `json:"points"`
	Meaning   string    `json:"meaning"`
}

// GeoPoint represents a latitude/longitude coordinate
type GeoPoint struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// LocationInfluence represents planetary influence at a specific location
type LocationInfluence struct {
	Location    string  `json:"location"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Influences  []Influence `json:"influences"`
}

// Influence represents a single planetary influence at a location
type Influence struct {
	Planet   string  `json:"planet"`
	LineType LineType `json:"line_type"`
	Distance float64 `json:"distance_degrees"`
	Strength string  `json:"strength"` // "very strong", "strong", "moderate", "weak"
	Meaning  string  `json:"meaning"`
}

// AstrocartoChart represents the full astrocartography data
type AstrocartoChart struct {
	Lines     []PlanetaryLine   `json:"lines"`
	Location  *LocationInfluence `json:"location_influence,omitempty"`
}

// Calculator performs astrocartography calculations
type Calculator struct {
	ephemeris *ephemeris.Ephemeris
}

// NewCalculator creates a new astrocartography calculator
func NewCalculator() *Calculator {
	return &Calculator{
		ephemeris: ephemeris.NewEphemeris(),
	}
}

// Calculate generates astrocartography data for a birth time
func (c *Calculator) Calculate(dt time.Time, viewLat, viewLon float64) (*AstrocartoChart, error) {
	chart := &AstrocartoChart{}

	positions := c.ephemeris.CalculatePositions(dt)

	// Calculate planetary lines for each planet
	for _, pos := range positions {
		if pos.Planet == ephemeris.Earth || pos.Planet == ephemeris.SouthNode {
			continue
		}

		// MC line: where the planet is on the Midheaven
		mcLine := c.calculateMCLine(dt, pos)
		mcLine.Meaning = getMeaning(string(pos.Planet), LineMC)
		chart.Lines = append(chart.Lines, mcLine)

		// IC line: opposite the MC
		icLine := c.calculateICLine(dt, pos)
		icLine.Meaning = getMeaning(string(pos.Planet), LineIC)
		chart.Lines = append(chart.Lines, icLine)

		// ASC line: where the planet is rising
		ascLine := c.calculateASCLine(dt, pos)
		ascLine.Meaning = getMeaning(string(pos.Planet), LineASC)
		chart.Lines = append(chart.Lines, ascLine)

		// DSC line: where the planet is setting
		dscLine := c.calculateDSCLine(dt, pos)
		dscLine.Meaning = getMeaning(string(pos.Planet), LineDSC)
		chart.Lines = append(chart.Lines, dscLine)
	}

	// Calculate influences at the viewing location
	chart.Location = c.calculateLocationInfluences(dt, positions, viewLat, viewLon)

	return chart, nil
}

// calculateMCLine finds longitudes where a planet is on the MC at various latitudes
func (c *Calculator) calculateMCLine(dt time.Time, pos ephemeris.PlanetaryPosition) PlanetaryLine {
	line := PlanetaryLine{
		Planet:   string(pos.Planet),
		LineType: LineMC,
	}

	// MC line is a vertical line (same longitude for all latitudes)
	// The planet is on the MC when the local sidereal time aligns with the planet's RA
	mcLon := c.longitudeForMC(dt, pos.Longitude)

	for lat := -70.0; lat <= 70.0; lat += 2.0 {
		line.Points = append(line.Points, GeoPoint{
			Latitude:  lat,
			Longitude: mcLon,
		})
	}

	return line
}

// calculateICLine finds longitudes where a planet is on the IC
func (c *Calculator) calculateICLine(dt time.Time, pos ephemeris.PlanetaryPosition) PlanetaryLine {
	line := PlanetaryLine{
		Planet:   string(pos.Planet),
		LineType: LineIC,
	}

	mcLon := c.longitudeForMC(dt, pos.Longitude)
	icLon := normalizeAngle(mcLon + 180.0)
	if icLon > 180 {
		icLon -= 360
	}

	for lat := -70.0; lat <= 70.0; lat += 2.0 {
		line.Points = append(line.Points, GeoPoint{
			Latitude:  lat,
			Longitude: icLon,
		})
	}

	return line
}

// calculateASCLine finds where a planet is on the Ascendant at various latitudes
func (c *Calculator) calculateASCLine(dt time.Time, pos ephemeris.PlanetaryPosition) PlanetaryLine {
	line := PlanetaryLine{
		Planet:   string(pos.Planet),
		LineType: LineASC,
	}

	obliquity := 23.4393 * math.Pi / 180.0
	planetRad := pos.Longitude * math.Pi / 180.0

	for lat := -65.0; lat <= 65.0; lat += 2.0 {
		latRad := lat * math.Pi / 180.0

		// For ASC line, we need the longitude where the planet's ecliptic longitude
		// is on the eastern horizon (ascendant) for this latitude
		// Using the co-ascendant formula
		sinA := math.Sin(planetRad)
		cosA := math.Cos(planetRad)
		tanLat := math.Tan(latRad)
		sinObl := math.Sin(obliquity)
		cosObl := math.Cos(obliquity)

		// RAMC when planet is ascending
		ramc := math.Atan2(-cosA, sinA*cosObl+tanLat*sinObl)
		ramcDeg := ramc * 180.0 / math.Pi

		// Convert RAMC to geographic longitude
		geoLon := c.ramcToLongitude(dt, ramcDeg)

		if geoLon > 180 {
			geoLon -= 360
		}
		if geoLon < -180 {
			geoLon += 360
		}

		line.Points = append(line.Points, GeoPoint{
			Latitude:  lat,
			Longitude: geoLon,
		})
	}

	return line
}

// calculateDSCLine finds where a planet is on the Descendant
func (c *Calculator) calculateDSCLine(dt time.Time, pos ephemeris.PlanetaryPosition) PlanetaryLine {
	line := PlanetaryLine{
		Planet:   string(pos.Planet),
		LineType: LineDSC,
	}

	obliquity := 23.4393 * math.Pi / 180.0
	// Descendant is 180° from the ascendant
	dscLon := normalizeAngle(pos.Longitude + 180.0)
	planetRad := dscLon * math.Pi / 180.0

	for lat := -65.0; lat <= 65.0; lat += 2.0 {
		latRad := lat * math.Pi / 180.0

		sinA := math.Sin(planetRad)
		cosA := math.Cos(planetRad)
		tanLat := math.Tan(latRad)
		sinObl := math.Sin(obliquity)
		cosObl := math.Cos(obliquity)

		ramc := math.Atan2(-cosA, sinA*cosObl+tanLat*sinObl)
		ramcDeg := ramc * 180.0 / math.Pi

		geoLon := c.ramcToLongitude(dt, ramcDeg)

		if geoLon > 180 {
			geoLon -= 360
		}
		if geoLon < -180 {
			geoLon += 360
		}

		line.Points = append(line.Points, GeoPoint{
			Latitude:  lat,
			Longitude: geoLon,
		})
	}

	return line
}

// longitudeForMC calculates the geographic longitude where a planet's ecliptic longitude
// equals the local MC
func (c *Calculator) longitudeForMC(dt time.Time, eclipticLon float64) float64 {
	obliquity := 23.4393 * math.Pi / 180.0
	lonRad := eclipticLon * math.Pi / 180.0

	// Right Ascension of the planet
	ra := math.Atan2(math.Sin(lonRad)*math.Cos(obliquity), math.Cos(lonRad))
	raDeg := ra * 180.0 / math.Pi
	if raDeg < 0 {
		raDeg += 360
	}

	// GMST at birth
	gmst := greenwichSiderealTime(dt)

	// Geographic longitude where RAMC = planet's RA
	geoLon := raDeg - gmst
	if geoLon > 180 {
		geoLon -= 360
	}
	if geoLon < -180 {
		geoLon += 360
	}

	return geoLon
}

// ramcToLongitude converts Right Ascension of MC to geographic longitude
func (c *Calculator) ramcToLongitude(dt time.Time, ramcDeg float64) float64 {
	gmst := greenwichSiderealTime(dt)
	geoLon := ramcDeg - gmst
	return normalizeAngle(geoLon)
}

// calculateLocationInfluences determines which planetary lines are near a location
func (c *Calculator) calculateLocationInfluences(dt time.Time, positions []ephemeris.PlanetaryPosition, lat, lon float64) *LocationInfluence {
	loc := &LocationInfluence{
		Latitude:  lat,
		Longitude: lon,
	}

	for _, pos := range positions {
		if pos.Planet == ephemeris.Earth || pos.Planet == ephemeris.SouthNode {
			continue
		}

		// Check MC line proximity
		mcLon := c.longitudeForMC(dt, pos.Longitude)
		mcDist := math.Abs(lon - mcLon)
		if mcDist > 180 {
			mcDist = 360 - mcDist
		}
		if mcDist < 15 {
			loc.Influences = append(loc.Influences, Influence{
				Planet:   string(pos.Planet),
				LineType: LineMC,
				Distance: math.Round(mcDist*100) / 100,
				Strength: getStrength(mcDist),
				Meaning:  getMeaning(string(pos.Planet), LineMC),
			})
		}

		// Check IC line proximity
		icLon := normalizeAngle(mcLon + 180)
		if icLon > 180 {
			icLon -= 360
		}
		icDist := math.Abs(lon - icLon)
		if icDist > 180 {
			icDist = 360 - icDist
		}
		if icDist < 15 {
			loc.Influences = append(loc.Influences, Influence{
				Planet:   string(pos.Planet),
				LineType: LineIC,
				Distance: math.Round(icDist*100) / 100,
				Strength: getStrength(icDist),
				Meaning:  getMeaning(string(pos.Planet), LineIC),
			})
		}
	}

	return loc
}

// greenwichSiderealTime calculates GMST in degrees
func greenwichSiderealTime(dt time.Time) float64 {
	year := float64(dt.Year())
	month := float64(dt.Month())
	day := float64(dt.Day())
	hour := float64(dt.Hour()) + float64(dt.Minute())/60.0 + float64(dt.Second())/3600.0

	if month <= 2 {
		year--
		month += 12
	}

	A := math.Floor(year / 100)
	B := 2 - A + math.Floor(A/4)
	jd := math.Floor(365.25*(year+4716)) + math.Floor(30.6001*(month+1)) + day + hour/24.0 + B - 1524.5

	T := (jd - 2451545.0) / 36525.0

	gmst := 280.46061837 + 360.98564736629*(jd-2451545.0) + 0.000387933*T*T - T*T*T/38710000.0
	gmst = math.Mod(gmst, 360.0)
	if gmst < 0 {
		gmst += 360.0
	}

	return gmst
}

func normalizeAngle(angle float64) float64 {
	angle = math.Mod(angle, 360)
	if angle < 0 {
		angle += 360
	}
	return angle
}

func getStrength(distance float64) string {
	if distance < 2 {
		return "very strong"
	} else if distance < 5 {
		return "strong"
	} else if distance < 10 {
		return "moderate"
	}
	return "weak"
}

// getMeaning returns the interpretation for a planet on a given angle
func getMeaning(planet string, lineType LineType) string {
	meanings := map[string]map[LineType]string{
		"Sun": {
			LineMC:  "Career prominence, public recognition, leadership opportunities",
			LineIC:  "Deep sense of belonging, family connections, feeling at home",
			LineASC: "Strong vitality, self-expression, personal identity shines",
			LineDSC: "Important partnerships, meeting significant others",
		},
		"Moon": {
			LineMC:  "Public emotional connection, nurturing career, popularity",
			LineIC:  "Emotional roots, comfort, domestic happiness",
			LineASC: "Emotional sensitivity heightened, intuitive connections",
			LineDSC: "Emotional partnerships, family-oriented relationships",
		},
		"Mercury": {
			LineMC:  "Communication in career, intellectual recognition, writing/speaking",
			LineIC:  "Mental stimulation at home, learning environment",
			LineASC: "Quick thinking, social connections, intellectual curiosity",
			LineDSC: "Intellectual partnerships, communication in relationships",
		},
		"Venus": {
			LineMC:  "Career in arts/beauty, social popularity, financial success",
			LineIC:  "Beautiful home, harmonious family life, comfort",
			LineASC: "Attractiveness, charm, pleasure, romantic encounters",
			LineDSC: "Love relationships, artistic partnerships, harmony",
		},
		"Mars": {
			LineMC:  "Ambitious drive, competitive career, leadership energy",
			LineIC:  "Active home life, renovation projects, family assertiveness",
			LineASC: "Physical energy, courage, pioneering spirit, independence",
			LineDSC: "Passionate relationships, dynamic partnerships",
		},
		"Jupiter": {
			LineMC:  "Career expansion, success, abundance, opportunities",
			LineIC:  "Large/comfortable home, family growth, generosity",
			LineASC: "Optimism, growth, luck, adventurous spirit",
			LineDSC: "Beneficial partnerships, generous relationships",
		},
		"Saturn": {
			LineMC:  "Career discipline, authority, long-term achievements",
			LineIC:  "Structured home, responsibility, ancestral connections",
			LineASC: "Self-discipline, maturity, serious demeanor",
			LineDSC: "Committed relationships, lessons through partnerships",
		},
		"Uranus": {
			LineMC:  "Unconventional career, innovation, sudden changes",
			LineIC:  "Unusual living situations, freedom in home life",
			LineASC: "Individuality, eccentricity, sudden awakenings",
			LineDSC: "Exciting but unstable relationships, freedom in partnerships",
		},
		"Neptune": {
			LineMC:  "Spiritual/creative career, inspiration, idealism",
			LineIC:  "Dreamy home life, spiritual family connections",
			LineASC: "Heightened intuition, spiritual experiences, creativity",
			LineDSC: "Soulmate connections, spiritual partnerships, idealized love",
		},
		"Pluto": {
			LineMC:  "Transformative career, power, deep impact on public",
			LineIC:  "Deep psychological roots, transformative home experiences",
			LineASC: "Personal transformation, intensity, magnetic presence",
			LineDSC: "Intense relationships, power dynamics in partnerships",
		},
		"North Node": {
			LineMC:  "Karmic career path, destined public role",
			LineIC:  "Soul's home, karmic family connections",
			LineASC: "Life purpose activated, destined path opens",
			LineDSC: "Fated relationships, karmic partnerships",
		},
	}

	if planetMeanings, ok := meanings[planet]; ok {
		if meaning, ok := planetMeanings[lineType]; ok {
			return meaning
		}
	}
	return "Planetary influence at this angle"
}
