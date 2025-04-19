package main

import (
	"fmt"
	"log"
	"os"

	"logreason/internal/csvparser"
	"logreason/internal/geojson"
	"logreason/internal/secrets"
)

func main() {
	// Default range value for GeoJSON API calls (in seconds)
	// This should eventually come from an external source
	const defaultRangeValue = 600

	// Create a new parser
	parser := csvparser.NewParser()

	// Define the path to the input CSV file
	inputFilePath := "locations/input.csv"

	// Check if the file exists
	if _, err := os.Stat(inputFilePath); os.IsNotExist(err) {
		log.Fatalf("Error: Input file %s does not exist", inputFilePath)
	}

	// Parse the CSV file
	fmt.Printf("Parsing CSV file: %s\n", inputFilePath)
	result := parser.ParseFile(inputFilePath)

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
	if err := secretsManager.LoadFromFile("config/secret.json"); err != nil {
		log.Fatalf("Error loading secrets: %v", err)
	}

	// Create a GeoJSON manager
	geoJSONManager, err := geojson.NewManager(secretsManager)
	if err != nil {
		log.Fatalf("Error creating GeoJSON manager: %v", err)
	}

	// Process the locations and save their GeoJSON data
	fmt.Println("Processing locations and saving GeoJSON data...")
	errors := geoJSONManager.ProcessLocations(result.Locations, defaultRangeValue)

	// Check if there were any errors during processing
	if len(errors) > 0 {
		log.Printf("Warning: There were %d errors during GeoJSON processing", len(errors))
		for _, err := range errors {
			log.Printf("  %v", err)
		}
	}

	fmt.Println("Done!")
}
