package main

import (
	"fmt"
	"humandesign/internal/database"
	"humandesign/internal/handlers"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize the people database
	dataDir := filepath.Join(".", "data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	db, err := database.NewDatabase(filepath.Join(dataDir, "people.json"))
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	h := handlers.NewHandler(db)

	// Page routes
	http.HandleFunc("/", h.HomePage)
	http.HandleFunc("/astrology", h.AstrologyPage)
	http.HandleFunc("/astrocartography", h.AstrocartographyPage)
	http.HandleFunc("/people", h.PeoplePage)
	http.HandleFunc("/people/", h.HandlePersonChart)

	// API routes
	http.HandleFunc("/api/reading", h.GenerateReading)
	http.HandleFunc("/api/reading/json", h.GetReadingJSON)
	http.HandleFunc("/api/astrology", h.GenerateAstrology)
	http.HandleFunc("/api/astrocartography", h.GenerateAstrocartography)
	http.HandleFunc("/api/people", h.HandlePeople)
	http.HandleFunc("/api/people/", h.HandlePerson)

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Cosmic Blueprint starting on http://localhost%s", addr)
	log.Printf("  /            - Human Design Calculator")
	log.Printf("  /astrology   - Natal Astrology Chart")
	log.Printf("  /astrocartography - Astrocartography")
	log.Printf("  /people      - People Database")

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
