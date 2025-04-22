package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"logreason/internal/csvparser"
	"logreason/internal/geojson"
	"logreason/internal/secrets"
)

func main() {
	// Define command line flags
	csvFilePath := flag.String("csv", "locations/input.csv", "Path to the input CSV file")
	rangeValue := flag.Int("range", 600, "Range value for GeoJSON API calls (in seconds)")
	outputDir := flag.String("output", "out/geojson", "Directory to save GeoJSON files")
	secretsFilePath := flag.String("secrets", "config/secret.json", "Path to the secrets file")
	flag.Parse()

	// Create a new parser
	parser := csvparser.NewParser()

	// Check if the file exists
	if _, err := os.Stat(*csvFilePath); os.IsNotExist(err) {
		log.Fatalf("Error: Input file %s does not exist", *csvFilePath)
	}

	// Parse the CSV file
	fmt.Printf("Parsing CSV file: %s\n", *csvFilePath)
	result := parser.ParseFile(*csvFilePath)

	// Check if parsing was successful
	if !result.Success {
		log.Printf("Warning: There were %d errors during parsing", len(result.Errors))
		for _, err := range result.Errors {
			log.Printf("  %s", err.Error())
		}
		if len(result.Locations) == 0 {
			log.Fatalf("Error: No valid locations found in the CSV file")
		}
	}

	fmt.Printf("Found %d locations\n", len(result.Locations))

	// Create a secrets manager and load secrets from file
	secretsManager := secrets.NewManager()
	if err := secretsManager.LoadFromFile(*secretsFilePath); err != nil {
		log.Fatalf("Error loading secrets: %v", err)
	}

	// Create a GeoJSON manager
	geoJSONManager, err := geojson.NewManager(secretsManager)
	if err != nil {
		log.Fatalf("Error creating GeoJSON manager: %v", err)
	}

	// Set the output directory
	if *outputDir != geojson.DefaultOutputDir {
		if err := geoJSONManager.SetOutputDir(*outputDir); err != nil {
			log.Fatalf("Error setting output directory: %v", err)
		}
	}

	// Ensure the output directory exists
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		log.Fatalf("Error creating output directory: %v", err)
	}

	// Process the locations and save their GeoJSON data
	fmt.Printf("Processing locations and saving GeoJSON data to %s...\n", *outputDir)
	errors := geoJSONManager.ProcessLocations(result.Locations, *rangeValue)

	// Check if there were any errors during processing
	if len(errors) > 0 {
		log.Printf("Warning: There were %d errors during GeoJSON processing", len(errors))
		for _, err := range errors {
			log.Printf("  %v", err)
		}
	}

	fmt.Println("Done!")
}
