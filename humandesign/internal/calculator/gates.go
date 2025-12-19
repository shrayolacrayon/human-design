package calculator

// GateInfo holds metadata about each gate
type GateInfo struct {
	Number int
	Name   string
	Center string
}

// AllGates maps gate numbers to their info
var AllGates = map[int]GateInfo{
	1:  {1, "Self-Expression", "G"},
	2:  {2, "The Direction of Self", "G"},
	3:  {3, "Ordering", "Sacral"},
	4:  {4, "Formulization", "Ajna"},
	5:  {5, "Fixed Patterns", "Sacral"},
	6:  {6, "Friction", "SolarPlexus"},
	7:  {7, "The Role of Self", "G"},
	8:  {8, "Contribution", "Throat"},
	9:  {9, "Focus", "Sacral"},
	10: {10, "Behavior of Self", "G"},
	11: {11, "Ideas", "Ajna"},
	12: {12, "Caution", "Throat"},
	13: {13, "The Listener", "G"},
	14: {14, "Power Skills", "Sacral"},
	15: {15, "Extremes", "G"},
	16: {16, "Skills", "Throat"},
	17: {17, "Opinions", "Ajna"},
	18: {18, "Correction", "Spleen"},
	19: {19, "Wanting", "Root"},
	20: {20, "The Now", "Throat"},
	21: {21, "Hunter/Huntress", "Heart"},
	22: {22, "Openness", "SolarPlexus"},
	23: {23, "Assimilation", "Throat"},
	24: {24, "Rationalization", "Ajna"},
	25: {25, "Spirit of Self", "G"},
	26: {26, "The Egoist", "Heart"},
	27: {27, "Caring", "Sacral"},
	28: {28, "The Game Player", "Spleen"},
	29: {29, "Perseverance", "Sacral"},
	30: {30, "Recognition of Feelings", "SolarPlexus"},
	31: {31, "Influence", "Throat"},
	32: {32, "Continuity", "Spleen"},
	33: {33, "Privacy", "Throat"},
	34: {34, "Power", "Sacral"},
	35: {35, "Change", "Throat"},
	36: {36, "Crisis", "SolarPlexus"},
	37: {37, "Friendship", "SolarPlexus"},
	38: {38, "The Fighter", "Root"},
	39: {39, "Provocation", "Root"},
	40: {40, "Aloneness", "Heart"},
	41: {41, "Contraction", "Root"},
	42: {42, "Growth", "Sacral"},
	43: {43, "Insight", "Ajna"},
	44: {44, "Alertness", "Spleen"},
	45: {45, "The Gatherer", "Throat"},
	46: {46, "Love of Body", "G"},
	47: {47, "Realization", "Ajna"},
	48: {48, "Depth", "Spleen"},
	49: {49, "Principles", "SolarPlexus"},
	50: {50, "Values", "Spleen"},
	51: {51, "Shock", "Heart"},
	52: {52, "Stillness", "Root"},
	53: {53, "Beginnings", "Root"},
	54: {54, "Ambition", "Root"},
	55: {55, "Spirit", "SolarPlexus"},
	56: {56, "Stimulation", "Throat"},
	57: {57, "Intuition", "Spleen"},
	58: {58, "Vitality", "Root"},
	59: {59, "Sexuality", "Sacral"},
	60: {60, "Acceptance", "Root"},
	61: {61, "Mystery", "Head"},
	62: {62, "Details", "Throat"},
	63: {63, "Doubt", "Head"},
	64: {64, "Confusion", "Head"},
}

// ChannelDefinition defines which gates form a channel and connects which centers
type ChannelDefinition struct {
	Name    string
	Gate1   int
	Gate2   int
	Center1 string
	Center2 string
}

// AllChannels defines all 36 channels in Human Design
var AllChannels = []ChannelDefinition{
	// Head to Ajna (3 channels)
	{"Abstraction", 64, 47, "Head", "Ajna"},
	{"Logic", 63, 4, "Head", "Ajna"},
	{"Awareness", 61, 24, "Head", "Ajna"},

	// Ajna to Throat (3 channels)
	{"Structuring", 43, 23, "Ajna", "Throat"},
	{"Curiosity", 11, 56, "Ajna", "Throat"},
	{"Acceptance", 17, 62, "Ajna", "Throat"},

	// Throat to G (4 channels)
	{"The Alpha", 7, 31, "G", "Throat"},
	{"The Prodigal", 13, 33, "G", "Throat"},
	{"Inspiration", 1, 8, "G", "Throat"},
	{"Awakening", 10, 20, "G", "Throat"},

	// Throat to Heart (2 channels)
	{"Money", 21, 45, "Heart", "Throat"},

	// Throat to Spleen (2 channels)
	{"The Brain Wave", 20, 57, "Spleen", "Throat"},
	{"Wavelength", 16, 48, "Spleen", "Throat"},

	// Throat to Solar Plexus (2 channels)
	{"Openness", 12, 22, "Throat", "SolarPlexus"},
	{"Transitoriness", 35, 36, "Throat", "SolarPlexus"},

	// G to Sacral (4 channels)
	{"Exploration", 10, 34, "G", "Sacral"},
	{"Discovery", 29, 46, "G", "Sacral"},
	{"The Beat", 2, 14, "G", "Sacral"},
	{"Rhythm", 5, 15, "G", "Sacral"},

	// G to Spleen (1 channel) - THIS WAS MISSING!
	{"Perfected Form", 10, 57, "G", "Spleen"},

	// G to Heart (1 channel)
	{"Initiation", 25, 51, "G", "Heart"},

	// Heart to Spleen (1 channel)
	{"Surrender", 26, 44, "Heart", "Spleen"},

	// Heart to Solar Plexus (1 channel)
	{"Community", 37, 40, "Heart", "SolarPlexus"},

	// Sacral to Throat (2 channels)
	{"Charisma", 20, 34, "Sacral", "Throat"},

	// Sacral to Spleen (2 channels)
	{"Power", 34, 57, "Sacral", "Spleen"},
	{"Preservation", 27, 50, "Sacral", "Spleen"},

	// Sacral to Solar Plexus (1 channel)
	{"Intimacy", 6, 59, "Sacral", "SolarPlexus"},

	// Sacral to Root (4 channels)
	{"Mutation", 3, 60, "Sacral", "Root"},
	{"Concentration", 9, 52, "Sacral", "Root"},
	{"Maturation", 42, 53, "Sacral", "Root"},

	// Spleen to Root (4 channels)
	{"Struggle", 28, 38, "Spleen", "Root"},
	{"Judgment", 18, 58, "Spleen", "Root"},
	{"Transformation", 32, 54, "Spleen", "Root"},

	// Solar Plexus to Root (3 channels)
	{"Recognition", 30, 41, "Root", "SolarPlexus"},
	{"Emoting", 39, 55, "Root", "SolarPlexus"},
	{"Sensitivity", 19, 49, "Root", "SolarPlexus"},
}

// GetGateCenter returns which center a gate belongs to
func GetGateCenter(gateNum int) string {
	if info, ok := AllGates[gateNum]; ok {
		return info.Center
	}
	return ""
}

// GetChannelForGates finds if two gates form a channel
func GetChannelForGates(gate1, gate2 int) *ChannelDefinition {
	for _, ch := range AllChannels {
		if (ch.Gate1 == gate1 && ch.Gate2 == gate2) ||
			(ch.Gate1 == gate2 && ch.Gate2 == gate1) {
			return &ch
		}
	}
	return nil
}
