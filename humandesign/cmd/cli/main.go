package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"humandesign/internal/calculator"
	"humandesign/internal/csvreader"
	"humandesign/internal/storage"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const version = "1.0.0"

// Default data directory for person storage
var defaultDataDir = filepath.Join(".", "data", "people")

func main() {
	// Define subcommands
	calculateCmd := flag.NewFlagSet("calculate", flag.ExitOnError)
	validateCmd := flag.NewFlagSet("validate", flag.ExitOnError)
	saveCmd := flag.NewFlagSet("save", flag.ExitOnError)
	loadCmd := flag.NewFlagSet("load", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	searchCmd := flag.NewFlagSet("search", flag.ExitOnError)
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
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

	// Save command flags
	saveName := saveCmd.String("name", "", "Person's name (required)")
	saveDate := saveCmd.String("date", "", "Birth date (RFC3339 format, required)")
	saveLat := saveCmd.Float64("lat", 0, "Birth latitude (required)")
	saveLon := saveCmd.Float64("lon", 0, "Birth longitude (required)")
	saveLoc := saveCmd.String("location", "", "Birth location name (required)")
	saveEmail := saveCmd.String("email", "", "Person's email (optional)")
	saveNotes := saveCmd.String("notes", "", "Notes about the person (optional)")
	saveTags := saveCmd.String("tags", "", "Comma-separated tags (optional)")
	saveDataDir := saveCmd.String("datadir", defaultDataDir, "Data directory")

	// Load command flags
	loadName := loadCmd.String("name", "", "Person's name")
	loadID := loadCmd.String("id", "", "Person's ID")
	loadOutput := loadCmd.String("output", "text", "Output format: text or json")
	loadDataDir := loadCmd.String("datadir", defaultDataDir, "Data directory")

	// List command flags
	listTags := listCmd.String("tags", "", "Filter by tags (comma-separated)")
	listOutput := listCmd.String("output", "text", "Output format: text or json")
	listDataDir := listCmd.String("datadir", defaultDataDir, "Data directory")

	// Search command flags
	searchQuery := searchCmd.String("query", "", "Search query (required)")
	searchOutput := searchCmd.String("output", "text", "Output format: text or json")
	searchDataDir := searchCmd.String("datadir", defaultDataDir, "Data directory")

	// Delete command flags
	deleteName := deleteCmd.String("name", "", "Person's name")
	deleteID := deleteCmd.String("id", "", "Person's ID")
	deleteDataDir := deleteCmd.String("datadir", defaultDataDir, "Data directory")

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
	case "save":
		saveCmd.Parse(os.Args[2:])
		handleSave(*saveName, *saveDate, *saveLat, *saveLon, *saveLoc, *saveEmail, *saveNotes, *saveTags, *saveDataDir)
	case "load":
		loadCmd.Parse(os.Args[2:])
		handleLoad(*loadName, *loadID, *loadOutput, *loadDataDir)
	case "list":
		listCmd.Parse(os.Args[2:])
		handleList(*listTags, *listOutput, *listDataDir)
	case "search":
		searchCmd.Parse(os.Args[2:])
		handleSearch(*searchQuery, *searchOutput, *searchDataDir)
	case "delete":
		deleteCmd.Parse(os.Args[2:])
		handleDelete(*deleteName, *deleteID, *deleteDataDir)
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
	fmt.Println("  save         Save a person to the database")
	fmt.Println("  load         Load a person from the database")
	fmt.Println("  list         List all people in the database")
	fmt.Println("  search       Search for people by name")
	fmt.Println("  delete       Delete a person from the database")
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
	fmt.Println("  # Save a person to the database")
	fmt.Println("  humandesign-cli save -name \"John Doe\" -date 1990-06-15T14:30:00Z -lat 40.7128 -lon -74.0060 -location \"New York\" -tags \"family,friends\"")
	fmt.Println()
	fmt.Println("  # Load a person from the database")
	fmt.Println("  humandesign-cli load -name \"John Doe\"")
	fmt.Println()
	fmt.Println("  # List all people")
	fmt.Println("  humandesign-cli list")
	fmt.Println()
	fmt.Println("  # Search for people")
	fmt.Println("  humandesign-cli search -query \"John\"")
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

// Person management handlers

func handleSave(name, date string, lat, lon float64, location, email, notes, tags, dataDir string) {
	if name == "" || date == "" || location == "" {
		fmt.Fprintf(os.Stderr, "Error: -name, -date, and -location are required\n")
		os.Exit(1)
	}

	// Parse date
	dateTime, err := time.Parse(time.RFC3339, date)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing date: %v\n", err)
		os.Exit(1)
	}

	// Create birth data
	birthData := calculator.BirthData{
		DateTime:  dateTime,
		Latitude:  lat,
		Longitude: lon,
		Location:  location,
	}

	// Calculate reading
	calc := calculator.NewCalculator()
	reading, err := calc.Calculate(birthData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error calculating reading: %v\n", err)
		os.Exit(1)
	}

	// Parse tags
	var tagList []string
	if tags != "" {
		tagList = strings.Split(tags, ",")
		for i := range tagList {
			tagList[i] = strings.TrimSpace(tagList[i])
		}
	}

	// Create person
	person := &storage.Person{
		Name:      name,
		Email:     email,
		Notes:     notes,
		BirthData: birthData,
		Reading:   reading,
		Tags:      tagList,
	}

	// Initialize storage
	store, err := storage.NewFileStore(dataDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing storage: %v\n", err)
		os.Exit(1)
	}
	defer store.Close()

	// Save person
	if err := store.Create(person); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving person: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Saved %s (ID: %s)\n", person.Name, person.ID)
	fmt.Printf("  Type: %s\n", reading.Type)
	fmt.Printf("  Authority: %s\n", reading.Authority)
	fmt.Printf("  Profile: %d/%d - %s\n", reading.Profile.Conscious, reading.Profile.Unconscious, reading.Profile.Name)
	fmt.Printf("  Strategy: %s\n", reading.Strategy)
}

func handleLoad(name, id, output, dataDir string) {
	if name == "" && id == "" {
		fmt.Fprintf(os.Stderr, "Error: -name or -id is required\n")
		os.Exit(1)
	}

	// Initialize storage
	store, err := storage.NewFileStore(dataDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing storage: %v\n", err)
		os.Exit(1)
	}
	defer store.Close()

	// Load person
	var person *storage.Person
	if id != "" {
		person, err = store.Get(id)
	} else {
		person, err = store.GetByName(name)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading person: %v\n", err)
		os.Exit(1)
	}

	// Output person
	if output == "json" {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(person); err != nil {
			fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		}
	} else {
		printPersonText(person)
	}
}

func handleList(tags, output, dataDir string) {
	// Parse tags
	var tagList []string
	if tags != "" {
		tagList = strings.Split(tags, ",")
		for i := range tagList {
			tagList[i] = strings.TrimSpace(tagList[i])
		}
	}

	// Initialize storage
	store, err := storage.NewFileStore(dataDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing storage: %v\n", err)
		os.Exit(1)
	}
	defer store.Close()

	// List people
	summaries, err := store.List(tagList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing people: %v\n", err)
		os.Exit(1)
	}

	if len(summaries) == 0 {
		fmt.Println("No people found in database")
		return
	}

	// Output summaries
	if output == "json" {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(summaries); err != nil {
			fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		}
	} else {
		fmt.Printf("Found %d people:\n\n", len(summaries))
		for _, summary := range summaries {
			fmt.Printf("%-30s | Type: %-25s | Authority: %-25s | Profile: %s\n",
				summary.Name,
				summary.Type,
				summary.Authority,
				summary.Profile,
			)
			if len(summary.Tags) > 0 {
				fmt.Printf("  Tags: %s\n", strings.Join(summary.Tags, ", "))
			}
		}
	}
}

func handleSearch(query, output, dataDir string) {
	if query == "" {
		fmt.Fprintf(os.Stderr, "Error: -query is required\n")
		os.Exit(1)
	}

	// Initialize storage
	store, err := storage.NewFileStore(dataDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing storage: %v\n", err)
		os.Exit(1)
	}
	defer store.Close()

	// Search
	summaries, err := store.Search(query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error searching: %v\n", err)
		os.Exit(1)
	}

	if len(summaries) == 0 {
		fmt.Printf("No people found matching '%s'\n", query)
		return
	}

	// Output summaries
	if output == "json" {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(summaries); err != nil {
			fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		}
	} else {
		fmt.Printf("Found %d people matching '%s':\n\n", len(summaries), query)
		for _, summary := range summaries {
			fmt.Printf("%-30s | Type: %-25s | Authority: %-25s\n",
				summary.Name,
				summary.Type,
				summary.Authority,
			)
		}
	}
}

func handleDelete(name, id, dataDir string) {
	if name == "" && id == "" {
		fmt.Fprintf(os.Stderr, "Error: -name or -id is required\n")
		os.Exit(1)
	}

	// Initialize storage
	store, err := storage.NewFileStore(dataDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing storage: %v\n", err)
		os.Exit(1)
	}
	defer store.Close()

	// Get person to confirm deletion
	var person *storage.Person
	if id != "" {
		person, err = store.Get(id)
	} else {
		person, err = store.GetByName(name)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding person: %v\n", err)
		os.Exit(1)
	}

	// Confirm deletion
	fmt.Printf("Are you sure you want to delete %s (ID: %s)? [y/N]: ", person.Name, person.ID)
	var response string
	fmt.Scanln(&response)

	if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
		fmt.Println("Deletion cancelled")
		return
	}

	// Delete person
	if err := store.Delete(person.ID); err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting person: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Deleted %s\n", person.Name)
}

func printPersonText(person *storage.Person) {
	fmt.Printf("Person: %s\n", person.Name)
	fmt.Printf("ID: %s\n", person.ID)
	if person.Email != "" {
		fmt.Printf("Email: %s\n", person.Email)
	}
	if person.Notes != "" {
		fmt.Printf("Notes: %s\n", person.Notes)
	}
	if len(person.Tags) > 0 {
		fmt.Printf("Tags: %s\n", strings.Join(person.Tags, ", "))
	}
	fmt.Printf("Created: %s\n", person.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Updated: %s\n\n", person.UpdatedAt.Format("2006-01-02 15:04:05"))

	if person.Reading != nil {
		printReadingText(*person.Reading)
	}
}
