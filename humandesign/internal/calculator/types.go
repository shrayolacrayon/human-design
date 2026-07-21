package calculator

import "time"

// BirthData represents the input for a Human Design reading
type BirthData struct {
	DateTime  time.Time `json:"datetime"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Location  string    `json:"location"`
}

// HumanDesignType represents the 4 main types in Human Design
type HumanDesignType string

const (
	TypeGenerator          HumanDesignType = "Generator"
	TypeManifestingGenerator HumanDesignType = "Manifesting Generator"
	TypeProjector          HumanDesignType = "Projector"
	TypeManifestor         HumanDesignType = "Manifestor"
	TypeReflector          HumanDesignType = "Reflector"
)

// Authority represents the decision-making authority
type Authority string

const (
	AuthorityEmotional     Authority = "Emotional (Solar Plexus)"
	AuthoritySacral        Authority = "Sacral"
	AuthoritySplenic       Authority = "Splenic"
	AuthorityEgo           Authority = "Ego/Heart"
	AuthoritySelf          Authority = "Self-Projected"
	AuthorityEnvironmental Authority = "Environmental (Mental Projector)"
	AuthorityLunar         Authority = "Lunar (Reflector)"
)

// Center represents one of the 9 centers in the body graph
type Center struct {
	Name    string `json:"name"`
	Defined bool   `json:"defined"`
	Gates   []int  `json:"gates"`
}

// Channel represents a connection between two centers
type Channel struct {
	Name     string `json:"name"`
	Gate1    int    `json:"gate1"`
	Gate2    int    `json:"gate2"`
	Center1  string `json:"center1"`
	Center2  string `json:"center2"`
	Defined  bool   `json:"defined"`
}

// Gate represents one of the 64 gates
type Gate struct {
	Number    int     `json:"number"`
	Name      string  `json:"name"`
	Line      int     `json:"line"`
	Planet    string  `json:"planet"`
	Longitude float64 `json:"longitude"` // ecliptic longitude in degrees
	Activated bool    `json:"activated"`
	Design    bool    `json:"design"` // true = design (red), false = personality (black)
}

// Profile represents the profile lines (e.g., 1/3, 4/6)
type Profile struct {
	Conscious   int    `json:"conscious"`
	Unconscious int    `json:"unconscious"`
	Name        string `json:"name"`
}

// Reading represents a complete Human Design reading
type Reading struct {
	BirthData       BirthData       `json:"birth_data"`
	Type            HumanDesignType `json:"type"`
	Authority       Authority       `json:"authority"`
	Profile         Profile         `json:"profile"`
	Definition      string          `json:"definition"`
	Strategy        string          `json:"strategy"`
	NotSelfTheme    string          `json:"not_self_theme"`
	Signature       string          `json:"signature"`
	Centers         map[string]*Center `json:"centers"`
	Channels        []Channel       `json:"channels"`
	PersonalityGates []Gate         `json:"personality_gates"` // Black - conscious
	DesignGates     []Gate          `json:"design_gates"`      // Red - unconscious
	IncarnationCross string         `json:"incarnation_cross"`
}

// ProfileNames maps profile numbers to their names
var ProfileNames = map[string]string{
	"1/3": "Investigator/Martyr",
	"1/4": "Investigator/Opportunist",
	"2/4": "Hermit/Opportunist",
	"2/5": "Hermit/Heretic",
	"3/5": "Martyr/Heretic",
	"3/6": "Martyr/Role Model",
	"4/6": "Opportunist/Role Model",
	"4/1": "Opportunist/Investigator",
	"5/1": "Heretic/Investigator",
	"5/2": "Heretic/Hermit",
	"6/2": "Role Model/Hermit",
	"6/3": "Role Model/Martyr",
}

// CenterNames lists all 9 centers
var CenterNames = []string{
	"Head",
	"Ajna",
	"Throat",
	"G",
	"Heart",
	"Sacral",
	"SolarPlexus",
	"Spleen",
	"Root",
}
