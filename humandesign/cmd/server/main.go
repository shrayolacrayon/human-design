package main

import (
	"fmt"
	"humandesign/internal/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create handler
	h := handlers.NewHandler()

	// Setup routes
	http.HandleFunc("/", h.HomePage)
	http.HandleFunc("/api/reading", h.GenerateReading)
	http.HandleFunc("/api/reading/json", h.GetReadingJSON)

	// Start server
	addr := fmt.Sprintf(":%s", port)
	log.Printf("🌟 Human Design Calculator starting on http://localhost%s", addr)
	log.Printf("   Open your browser and visit the URL above")
	
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
