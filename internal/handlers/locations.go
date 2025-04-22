// Package handlers provides HTTP request handlers for the application
package handlers

import (
	"os"

	"github.com/gofiber/fiber/v2"

	"logreason/internal/csvparser"
)

// GetLocationsCsv returns the locations/input.csv file as an attachment
func GetLocationsCsv(c *fiber.Ctx) error {
	filePath := "locations/input.csv"

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.Status(fiber.StatusNotFound).SendString("CSV file not found")
	}

	return c.Download(filePath, "input.csv")
}

// GetLocationsJson returns the parsed content of locations/input.csv as a JSON array
func GetLocationsJson(c *fiber.Ctx) error {
	filePath := "locations/input.csv"

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.Status(fiber.StatusNotFound).SendString("CSV file not found")
	}

	// Create a new parser
	parser := csvparser.NewParser()

	// Parse the CSV file
	result := parser.ParseFile(filePath)

	// Check if parsing was successful
	if !result.Success && len(result.Locations) == 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"errors":  result.Errors,
		})
	}

	return c.JSON(result.Locations)
}