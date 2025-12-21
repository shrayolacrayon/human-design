package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"humandesign/internal/calculator"
	"humandesign/internal/csvreader"
	"os"
	"strings"
	"time"
)

const version = "1.0.0"

func main() {
	// Define subcommands
	calculateCmd := flag.NewFlagSet("calculate", flag.ExitOnError)
	validateCmd := flag.NewFlagSet("validate", flag.ExitOnError)
	versionCmd := flag.NewFlagSet("version", flag.ExitOnError)

	// Calculate command flags
	calcCSV := calculateCmd.String("csv", "", "CSV file with birth data")
	calcName := calculateCmd.String("name", "", "Person's name")
	calcDate := calculateCmd.String("date", "", "Birth date (RFC3339 format, e.g., 1990-06-15T14:30:00Z)")
	calcLat := calculateCmd.Float64("lat", 0, "Birth latitude")
	calcLon := calculateCmd.Float64("lon", 0, "Birth longitude")
	calcLoc := calculateCmd.String("location", "", "Birth location name")
	calcOutput := calculateCmd.String("output", "text", "Output format: text, json, or csv")

	// Validate command flags
	valCSV := validateCmd.String("csv", "", "CSV file with test cases (required)")
	valOutput := validateCmd.String("output", "text", "Output format: text or json")
	valVerbose := validateCmd.Bool("verbose", false, "Verbose output")

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "calculate":
		calculateCmd.Parse(os.Args[2:])
		handleCalculate(*calcCSV, *calcName, *calcDate, *calcLat, *calcLon, *calcLoc, *calcOutput)
	case "validate":
		validateCmd.Parse(os.Args[2:])
		handleValidate(*valCSV, *valOutput, *valVerbose)
	case "version":
		versionCmd.Parse(os.Args[2:])
		fmt.Printf("Human Design CLI v%s\n", version)
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Human Design CLI - Calculate and validate Human Design charts")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  humandesign-cli <command> [flags]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  calculate    Calculate Human Design reading")
	fmt.Println("  validate     Validate readings against CSV test cases")
	fmt.Println("  version      Show version information")
	fmt.Println("  help         Show this help message")
	fmt.Println()
	fmt.Println("Calculate flags:")
	fmt.Println("  -csv <file>        CSV file with birth data")
	fmt.Println("  -name <name>       Person's name")
	fmt.Println("  -date <datetime>   Birth date (RFC3339 format)")
	fmt.Println("  -lat <latitude>    Birth latitude")
	fmt.Println("  -lon <longitude>   Birth longitude")
	fmt.Println("  -location <loc>    Birth location name")
	fmt.Println("  -output <format>   Output format: text, json, csv (default: text)")
	fmt.Println()
	fmt.Println("Validate flags:")
	fmt.Println("  -csv <file>        CSV file with test cases (required)")
	fmt.Println("  -output <format>   Output format: text or json (default: text)")
	fmt.Println("  -verbose           Show detailed output")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  # Calculate from individual parameters")
	fmt.Println("  humandesign-cli calculate -name \"John Doe\" -date 1990-06-15T14:30:00Z -lat 40.7128 -lon -74.0060 -location \"New York\"")
	fmt.Println()
	fmt.Println("  # Calculate from CSV file")
	fmt.Println("  humandesign-cli calculate -csv testdata/birth_data.csv -output json")
	fmt.Println()
	fmt.Println("  # Validate against test cases")
	fmt.Println("  humandesign-cli validate -csv testdata/birth_data.csv -verbose")
}

func handleCalculate(csvFile, name, date string, lat, lon float64, location, output string) {
	calc := calculator.NewCalculator()

	if csvFile != "" {
		// Read from CSV and calculate for all entries
		testCases, err := csvreader.ReadBirthDataCSV(csvFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading CSV: %v\n", err)
			os.Exit(1)
		}

		results := make([]calculator.Reading, 0, len(testCases))
		for _, tc := range testCases {
			reading, err := calc.Calculate(tc.BirthData)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error calculating for %s: %v\n", tc.Name, err)
				continue
			}
			results = append(results, *reading)
		}

		printResults(results, output)
	} else {
		// Calculate from individual parameters
		if name == "" || date == "" {
			fmt.Fprintf(os.Stderr, "Error: -name and -date are required when not using -csv\n")
			os.Exit(1)
		}

		dateTime, err := time.Parse(time.RFC3339, date)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing date: %v\n", err)
			os.Exit(1)
		}

		birthData := calculator.BirthData{
			DateTime:  dateTime,
			Latitude:  lat,
			Longitude: lon,
			Location:  location,
		}

		reading, err := calc.Calculate(birthData)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error calculating: %v\n", err)
			os.Exit(1)
		}

		printResults([]calculator.Reading{*reading}, output)
	}
}

func handleValidate(csvFile, output string, verbose bool) {
	if csvFile == "" {
		fmt.Fprintf(os.Stderr, "Error: -csv flag is required for validate command\n")
		os.Exit(1)
	}

	testCases, err := csvreader.ReadBirthDataCSV(csvFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading CSV: %v\n", err)
		os.Exit(1)
	}

	calc := calculator.NewCalculator()
	results := make([]ValidationResult, 0, len(testCases))

	for _, tc := range testCases {
		reading, err := calc.Calculate(tc.BirthData)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error calculating for %s: %v\n", tc.Name, err)
			continue
		}

		result := validateReading(tc, reading)
		results = append(results, result)

		if output == "text" {
			printValidationResult(result, verbose)
		}
	}

	if output == "json" {
		printValidationResultsJSON(results)
	}

	// Print summary
	if output == "text" {
		printValidationSummary(results)
	}
}

type ValidationResult struct {
	Name              string   `json:"name"`
	Passed            bool     `json:"passed"`
	TypeMatch         bool     `json:"type_match"`
	AuthorityMatch    bool     `json:"authority_match"`
	StrategyMatch     bool     `json:"strategy_match"`
	ProfileMatch      bool     `json:"profile_match"`
	ExpectedType      string   `json:"expected_type,omitempty"`
	ActualType        string   `json:"actual_type,omitempty"`
	ExpectedAuthority string   `json:"expected_authority,omitempty"`
	ActualAuthority   string   `json:"actual_authority,omitempty"`
	ExpectedStrategy  string   `json:"expected_strategy,omitempty"`
	ActualStrategy    string   `json:"actual_strategy,omitempty"`
	ExpectedProfile   string   `json:"expected_profile,omitempty"`
	ActualProfile     string   `json:"actual_profile,omitempty"`
	Messages          []string `json:"messages,omitempty"`
}

func validateReading(tc csvreader.TestCase, reading *calculator.Reading) ValidationResult {
	result := ValidationResult{
		Name:   tc.Name,
		Passed: true,
	}

	// Validate Type
	if tc.ExpectedType != "" {
		result.ExpectedType = tc.ExpectedType
		result.ActualType = string(reading.Type)
		result.TypeMatch = string(reading.Type) == tc.ExpectedType
		if !result.TypeMatch {
			result.Passed = false
			result.Messages = append(result.Messages,
				fmt.Sprintf("Type mismatch: expected %s, got %s", tc.ExpectedType, reading.Type))
		}
	}

	// Validate Authority
	if tc.ExpectedAuthority != "" {
		result.ExpectedAuthority = tc.ExpectedAuthority
		result.ActualAuthority = string(reading.Authority)
		result.AuthorityMatch = string(reading.Authority) == tc.ExpectedAuthority
		if !result.AuthorityMatch {
			result.Passed = false
			result.Messages = append(result.Messages,
				fmt.Sprintf("Authority mismatch: expected %s, got %s", tc.ExpectedAuthority, reading.Authority))
		}
	}

	// Validate Strategy
	if tc.ExpectedStrategy != "" {
		result.ExpectedStrategy = tc.ExpectedStrategy
		result.ActualStrategy = reading.Strategy
		result.StrategyMatch = reading.Strategy == tc.ExpectedStrategy
		if !result.StrategyMatch {
			result.Passed = false
			result.Messages = append(result.Messages,
				fmt.Sprintf("Strategy mismatch: expected %s, got %s", tc.ExpectedStrategy, reading.Strategy))
		}
	}

	// Validate Profile
	if tc.ExpectedProfileConscious > 0 && tc.ExpectedProfileUnconscious > 0 {
		expectedProfile := fmt.Sprintf("%d/%d", tc.ExpectedProfileConscious, tc.ExpectedProfileUnconscious)
		actualProfile := fmt.Sprintf("%d/%d", reading.Profile.Conscious, reading.Profile.Unconscious)
		result.ExpectedProfile = expectedProfile
		result.ActualProfile = actualProfile
		result.ProfileMatch = expectedProfile == actualProfile
		if !result.ProfileMatch {
			result.Passed = false
			result.Messages = append(result.Messages,
				fmt.Sprintf("Profile mismatch: expected %s, got %s", expectedProfile, actualProfile))
		}
	}

	return result
}

func printValidationResult(result ValidationResult, verbose bool) {
	status := "✓ PASS"
	if !result.Passed {
		status = "✗ FAIL"
	}

	fmt.Printf("%s %s\n", status, result.Name)

	if !result.Passed || verbose {
		for _, msg := range result.Messages {
			fmt.Printf("  - %s\n", msg)
		}
	}
}

func printValidationResultsJSON(results []ValidationResult) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(results); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
	}
}

func printValidationSummary(results []ValidationResult) {
	passed := 0
	failed := 0

	for _, result := range results {
		if result.Passed {
			passed++
		} else {
			failed++
		}
	}

	total := passed + failed
	fmt.Println()
	fmt.Println(strings.Repeat("-", 50))
	fmt.Printf("Total: %d | Passed: %d | Failed: %d\n", total, passed, failed)
	if total > 0 {
		fmt.Printf("Success Rate: %.1f%%\n", float64(passed)/float64(total)*100)
	}
}

func printResults(readings []calculator.Reading, format string) {
	switch format {
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(readings); err != nil {
			fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		}
	case "csv":
		printResultsCSV(readings)
	default: // text
		for i, reading := range readings {
			if i > 0 {
				fmt.Println()
				fmt.Println(strings.Repeat("-", 70))
				fmt.Println()
			}
			printReadingText(reading)
		}
	}
}

func printReadingText(reading calculator.Reading) {
	fmt.Printf("Human Design Reading\n")
	fmt.Printf("Birth Date: %s\n", reading.BirthData.DateTime.Format("2006-01-02 15:04:05 MST"))
	fmt.Printf("Location: %s (%.4f, %.4f)\n\n",
		reading.BirthData.Location, reading.BirthData.Latitude, reading.BirthData.Longitude)

	fmt.Printf("Type: %s\n", reading.Type)
	fmt.Printf("Strategy: %s\n", reading.Strategy)
	fmt.Printf("Authority: %s\n", reading.Authority)
	fmt.Printf("Profile: %d/%d - %s\n",
		reading.Profile.Conscious, reading.Profile.Unconscious, reading.Profile.Name)
	fmt.Printf("Definition: %s\n", reading.Definition)
	fmt.Printf("Signature: %s\n", reading.Signature)
	fmt.Printf("Not-Self Theme: %s\n", reading.NotSelfTheme)
	fmt.Printf("Incarnation Cross: %s\n\n", reading.IncarnationCross)

	// Print Centers
	fmt.Println("Centers:")
	for _, centerName := range calculator.CenterNames {
		if center, ok := reading.Centers[centerName]; ok {
			status := "Undefined"
			if center.Defined {
				status = "Defined"
			}
			fmt.Printf("  %-15s: %s", centerName, status)
			if len(center.Gates) > 0 {
				fmt.Printf(" (Gates: %v)", center.Gates)
			}
			fmt.Println()
		}
	}

	// Print Defined Channels
	fmt.Println("\nDefined Channels:")
	hasDefinedChannels := false
	for _, channel := range reading.Channels {
		if channel.Defined {
			hasDefinedChannels = true
			fmt.Printf("  %s (%d-%d): %s ↔ %s\n",
				channel.Name, channel.Gate1, channel.Gate2, channel.Center1, channel.Center2)
		}
	}
	if !hasDefinedChannels {
		fmt.Println("  None")
	}

	// Print Personality Gates
	fmt.Println("\nPersonality Gates (Conscious):")
	for _, gate := range reading.PersonalityGates {
		fmt.Printf("  Gate %d.%d - %s (%s)\n",
			gate.Number, gate.Line, gate.Name, gate.Planet)
	}

	// Print Design Gates
	fmt.Println("\nDesign Gates (Unconscious):")
	for _, gate := range reading.DesignGates {
		fmt.Printf("  Gate %d.%d - %s (%s)\n",
			gate.Number, gate.Line, gate.Name, gate.Planet)
	}
}

func printResultsCSV(readings []calculator.Reading) {
	// Print header
	fmt.Println("name,type,strategy,authority,profile,definition,signature,not_self_theme,incarnation_cross")

	// Print each reading
	for _, reading := range readings {
		profile := fmt.Sprintf("%d/%d", reading.Profile.Conscious, reading.Profile.Unconscious)
		fmt.Printf("%s,%s,%s,%s,%s,%s,%s,%s,%s\n",
			reading.BirthData.Location,
			reading.Type,
			reading.Strategy,
			reading.Authority,
			profile,
			reading.Definition,
			reading.Signature,
			reading.NotSelfTheme,
			reading.IncarnationCross,
		)
	}
}
