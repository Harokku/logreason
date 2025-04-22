// Package handlers provides HTTP request handlers for the application
package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// GetAllGeoJson returns all GeoJSON files from out/geojson directory as a combined JSON array
func GetAllGeoJson(c *fiber.Ctx) error {
	dirPath := "out/geojson"

	// Check if directory exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return c.Status(fiber.StatusNotFound).SendString("GeoJSON directory not found")
	}

	// Read all files in the directory
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Error reading directory: %v", err))
	}

	// Combine all GeoJSON files
	var result []json.RawMessage

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			filePath := filepath.Join(dirPath, entry.Name())

			// Read file content
			content, err := os.ReadFile(filePath)
			if err != nil {
				log.Printf("Error reading file %s: %v", filePath, err)
				continue
			}

			// Parse JSON
			var jsonData json.RawMessage
			if err := json.Unmarshal(content, &jsonData); err != nil {
				log.Printf("Error parsing JSON from file %s: %v", filePath, err)
				continue
			}

			result = append(result, jsonData)
		}
	}

	return c.JSON(result)
}

// GetGeoJsonByName returns a specific GeoJSON file by name as a JSON object
func GetGeoJsonByName(c *fiber.Ctx) error {
	name := c.Params("name")
	if name == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Name parameter is required")
	}

	filePath := filepath.Join("out/geojson", name+".json")

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.Status(fiber.StatusNotFound).SendString(fmt.Sprintf("GeoJSON file %s not found", name))
	}

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Error reading file: %v", err))
	}

	// Parse JSON
	var jsonData json.RawMessage
	if err := json.Unmarshal(content, &jsonData); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Error parsing JSON: %v", err))
	}

	return c.JSON(jsonData)
}

// GetFilteredGeoJson returns multiple specific GeoJSON files as a combined JSON array
func GetFilteredGeoJson(c *fiber.Ctx) error {
	namesParam := c.Query("names")
	if namesParam == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Names parameter is required")
	}

	// Split names by comma
	names := strings.Split(namesParam, ",")

	// Combine specified GeoJSON files
	var result []json.RawMessage

	for _, name := range names {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}

		filePath := filepath.Join("out/geojson", name+".json")

		// Check if file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			log.Printf("GeoJSON file %s not found", name)
			continue
		}

		// Read file content
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Error reading file %s: %v", filePath, err)
			continue
		}

		// Parse JSON
		var jsonData json.RawMessage
		if err := json.Unmarshal(content, &jsonData); err != nil {
			log.Printf("Error parsing JSON from file %s: %v", filePath, err)
			continue
		}

		result = append(result, jsonData)
	}

	if len(result) == 0 {
		return c.Status(fiber.StatusNotFound).SendString("No valid GeoJSON files found for the specified names")
	}

	return c.JSON(result)
}
