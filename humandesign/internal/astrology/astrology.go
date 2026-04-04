package astrology

import (
	"fmt"
	"humandesign/internal/ephemeris"
	"math"
	"time"
)

// ZodiacSign represents one of the 12 zodiac signs
type ZodiacSign struct {
	Name     string  `json:"name"`
	Symbol   string  `json:"symbol"`
	Element  string  `json:"element"`
	Modality string  `json:"modality"`
	Ruler    string  `json:"ruler"`
	StartDeg float64 `json:"-"`
}

// PlanetPlacement represents a planet's position in a sign/house
type PlanetPlacement struct {
	Planet    string  `json:"planet"`
	Sign      string  `json:"sign"`
	SignSymbol string  `json:"sign_symbol"`
	Degree    float64 `json:"degree"`
	DegreeInSign float64 `json:"degree_in_sign"`
	House     int     `json:"house"`
	Retrograde bool   `json:"retrograde"`
}

// Aspect represents an angular relationship between two planets
type Aspect struct {
	Planet1   string  `json:"planet1"`
	Planet2   string  `json:"planet2"`
	Type      string  `json:"type"`
	Angle     float64 `json:"angle"`
	Orb       float64 `json:"orb"`
	Harmony   string  `json:"harmony"` // "harmonious", "challenging", "neutral"
}

// HouseData represents a house cusp
type HouseData struct {
	Number   int     `json:"number"`
	Sign     string  `json:"sign"`
	Degree   float64 `json:"degree"`
}

// NatalChart represents a full natal astrology chart
type NatalChart struct {
	Planets     []PlanetPlacement `json:"planets"`
	Houses      []HouseData       `json:"houses"`
	Aspects     []Aspect          `json:"aspects"`
	Ascendant   float64           `json:"ascendant"`
	Midheaven   float64           `json:"midheaven"`
	AscSign     string            `json:"asc_sign"`
	MCSign      string            `json:"mc_sign"`
	SunSign     string            `json:"sun_sign"`
	MoonSign    string            `json:"moon_sign"`
	RisingSign  string            `json:"rising_sign"`
	Elements    map[string]int    `json:"elements"`
	Modalities  map[string]int    `json:"modalities"`
}

// Calculator performs astrology calculations
type Calculator struct {
	ephemeris *ephemeris.Ephemeris
}

// NewCalculator creates a new astrology calculator
func NewCalculator() *Calculator {
	return &Calculator{
		ephemeris: ephemeris.NewEphemeris(),
	}
}

// Signs holds all 12 zodiac signs in order
var Signs = []ZodiacSign{
	{Name: "Aries", Symbol: "\u2648", Element: "Fire", Modality: "Cardinal", Ruler: "Mars", StartDeg: 0},
	{Name: "Taurus", Symbol: "\u2649", Element: "Earth", Modality: "Fixed", Ruler: "Venus", StartDeg: 30},
	{Name: "Gemini", Symbol: "\u264A", Element: "Air", Modality: "Mutable", Ruler: "Mercury", StartDeg: 60},
	{Name: "Cancer", Symbol: "\u264B", Element: "Water", Modality: "Cardinal", Ruler: "Moon", StartDeg: 90},
	{Name: "Leo", Symbol: "\u264C", Element: "Fire", Modality: "Fixed", Ruler: "Sun", StartDeg: 120},
	{Name: "Virgo", Symbol: "\u264D", Element: "Earth", Modality: "Mutable", Ruler: "Mercury", StartDeg: 150},
	{Name: "Libra", Symbol: "\u264E", Element: "Air", Modality: "Cardinal", Ruler: "Venus", StartDeg: 180},
	{Name: "Scorpio", Symbol: "\u264F", Element: "Water", Modality: "Fixed", Ruler: "Pluto", StartDeg: 210},
	{Name: "Sagittarius", Symbol: "\u2650", Element: "Fire", Modality: "Mutable", Ruler: "Jupiter", StartDeg: 240},
	{Name: "Capricorn", Symbol: "\u2651", Element: "Earth", Modality: "Cardinal", Ruler: "Saturn", StartDeg: 270},
	{Name: "Aquarius", Symbol: "\u2652", Element: "Air", Modality: "Fixed", Ruler: "Uranus", StartDeg: 300},
	{Name: "Pisces", Symbol: "\u2653", Element: "Water", Modality: "Mutable", Ruler: "Neptune", StartDeg: 330},
}

// AspectDefinition defines aspect types
type AspectDefinition struct {
	Name    string
	Angle   float64
	Orb     float64
	Harmony string
}

var AspectDefs = []AspectDefinition{
	{"Conjunction", 0, 8, "neutral"},
	{"Opposition", 180, 8, "challenging"},
	{"Trine", 120, 8, "harmonious"},
	{"Square", 90, 7, "challenging"},
	{"Sextile", 60, 6, "harmonious"},
	{"Quincunx", 150, 3, "challenging"},
	{"Semi-Sextile", 30, 2, "neutral"},
}

// GetSign returns the zodiac sign for a given ecliptic longitude
func GetSign(longitude float64) ZodiacSign {
	longitude = normalizeAngle(longitude)
	idx := int(longitude / 30.0)
	if idx >= 12 {
		idx = 11
	}
	return Signs[idx]
}

// CalculateChart generates a full natal chart
func (c *Calculator) CalculateChart(dt time.Time, latitude, longitude float64) (*NatalChart, error) {
	chart := &NatalChart{
		Elements:   make(map[string]int),
		Modalities: make(map[string]int),
	}

	// Get planetary positions from ephemeris
	positions := c.ephemeris.CalculatePositions(dt)

	// Calculate Ascendant and Midheaven
	chart.Ascendant = calculateAscendant(dt, latitude, longitude)
	chart.Midheaven = calculateMidheaven(dt, longitude)

	ascSign := GetSign(chart.Ascendant)
	mcSign := GetSign(chart.Midheaven)
	chart.AscSign = ascSign.Name
	chart.MCSign = mcSign.Name
	chart.RisingSign = ascSign.Name

	// Calculate house cusps using Equal House system (from Ascendant)
	chart.Houses = calculateHouses(chart.Ascendant)

	// Convert positions to placements
	for _, pos := range positions {
		if pos.Planet == ephemeris.Earth {
			continue // Skip Earth in traditional astrology
		}

		sign := GetSign(pos.Longitude)
		degInSign := pos.Longitude - sign.StartDeg
		if degInSign < 0 {
			degInSign += 360
		}
		degInSign = math.Mod(degInSign, 30.0)

		house := getHouseForDegree(pos.Longitude, chart.Houses)

		placement := PlanetPlacement{
			Planet:       string(pos.Planet),
			Sign:         sign.Name,
			SignSymbol:   sign.Symbol,
			Degree:       pos.Longitude,
			DegreeInSign: degInSign,
			House:        house,
		}

		chart.Planets = append(chart.Planets, placement)

		// Track elements and modalities for personal planets
		chart.Elements[sign.Element]++
		chart.Modalities[sign.Modality]++

		if pos.Planet == ephemeris.Sun {
			chart.SunSign = sign.Name
		}
		if pos.Planet == ephemeris.Moon {
			chart.MoonSign = sign.Name
		}
	}

	// Calculate aspects between planets
	chart.Aspects = calculateAspects(chart.Planets)

	return chart, nil
}

// calculateAscendant computes the Ascendant (rising sign degree)
func calculateAscendant(dt time.Time, lat, lon float64) float64 {
	// Calculate Local Sidereal Time
	lst := localSiderealTime(dt, lon)

	// Convert to radians
	lstRad := lst * math.Pi / 12.0 // LST is in hours, convert to radians
	latRad := lat * math.Pi / 180.0

	// Obliquity of the ecliptic (approximately 23.4393 degrees)
	obliquity := 23.4393 * math.Pi / 180.0

	// Ascendant formula
	y := -math.Cos(lstRad)
	x := math.Sin(lstRad)*math.Cos(obliquity) + math.Tan(latRad)*math.Sin(obliquity)

	asc := math.Atan2(y, x) * 180.0 / math.Pi
	asc = normalizeAngle(asc)

	return asc
}

// calculateMidheaven computes the MC (Medium Coeli)
func calculateMidheaven(dt time.Time, lon float64) float64 {
	lst := localSiderealTime(dt, lon)
	lstRad := lst * math.Pi / 12.0
	obliquity := 23.4393 * math.Pi / 180.0

	mc := math.Atan2(math.Sin(lstRad), math.Cos(lstRad)*math.Cos(obliquity)) * 180.0 / math.Pi
	mc = normalizeAngle(mc)

	return mc
}

// localSiderealTime calculates the Local Sidereal Time in hours
func localSiderealTime(dt time.Time, lon float64) float64 {
	// Julian Date
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

	// Julian centuries from J2000.0
	T := (jd - 2451545.0) / 36525.0

	// Greenwich Mean Sidereal Time in degrees
	gmst := 280.46061837 + 360.98564736629*(jd-2451545.0) + 0.000387933*T*T - T*T*T/38710000.0
	gmst = math.Mod(gmst, 360.0)
	if gmst < 0 {
		gmst += 360.0
	}

	// Convert to Local Sidereal Time
	lst := gmst + lon
	lst = math.Mod(lst, 360.0)
	if lst < 0 {
		lst += 360.0
	}

	// Convert from degrees to hours
	return lst / 15.0
}

// calculateHouses using Equal House system (each house = 30 degrees from ASC)
func calculateHouses(ascendant float64) []HouseData {
	houses := make([]HouseData, 12)
	for i := 0; i < 12; i++ {
		cusp := normalizeAngle(ascendant + float64(i)*30.0)
		sign := GetSign(cusp)
		houses[i] = HouseData{
			Number: i + 1,
			Sign:   sign.Name,
			Degree: cusp,
		}
	}
	return houses
}

// getHouseForDegree determines which house a degree falls in
func getHouseForDegree(degree float64, houses []HouseData) int {
	degree = normalizeAngle(degree)
	for i := 0; i < 12; i++ {
		nextIdx := (i + 1) % 12
		start := houses[i].Degree
		end := houses[nextIdx].Degree

		if end < start { // wraps around 360
			if degree >= start || degree < end {
				return i + 1
			}
		} else {
			if degree >= start && degree < end {
				return i + 1
			}
		}
	}
	return 1
}

// calculateAspects finds all aspects between planets
func calculateAspects(planets []PlanetPlacement) []Aspect {
	var aspects []Aspect

	for i := 0; i < len(planets); i++ {
		for j := i + 1; j < len(planets); j++ {
			p1 := planets[i]
			p2 := planets[j]

			// Skip South Node aspects (redundant with North Node)
			if p1.Planet == "South Node" || p2.Planet == "South Node" {
				continue
			}

			diff := math.Abs(p1.Degree - p2.Degree)
			if diff > 180 {
				diff = 360 - diff
			}

			for _, def := range AspectDefs {
				orb := math.Abs(diff - def.Angle)
				if orb <= def.Orb {
					aspects = append(aspects, Aspect{
						Planet1: p1.Planet,
						Planet2: p2.Planet,
						Type:    def.Name,
						Angle:   def.Angle,
						Orb:     math.Round(orb*100) / 100,
						Harmony: def.Harmony,
					})
					break
				}
			}
		}
	}

	return aspects
}

// FormatDegree formats a longitude as "DD° Sign MM'"
func FormatDegree(longitude float64) string {
	sign := GetSign(longitude)
	degInSign := longitude - sign.StartDeg
	if degInSign < 0 {
		degInSign += 360
	}
	degInSign = math.Mod(degInSign, 30.0)

	deg := int(degInSign)
	min := int((degInSign - float64(deg)) * 60)

	return fmt.Sprintf("%d° %s %d'", deg, sign.Name, min)
}

func normalizeAngle(angle float64) float64 {
	angle = math.Mod(angle, 360)
	if angle < 0 {
		angle += 360
	}
	return angle
}
