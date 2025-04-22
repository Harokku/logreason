// Package routes provides API route definitions for the application
package routes

import (
	"github.com/gofiber/fiber/v2"

	"logreason/internal/api"
	"logreason/internal/handlers"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(app *fiber.App) {
	// Define a route for the root endpoint with API documentation
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString(api.Documentation)
	})

	// API routes
	apiGroup := app.Group("/api")

	// CSV routes
	apiGroup.Get("/locations/csv", handlers.GetLocationsCsv)
	apiGroup.Get("/locations/json", handlers.GetLocationsJson)

	// GeoJSON routes
	apiGroup.Get("/geojson", handlers.GetAllGeoJson)
	apiGroup.Get("/geojson/filter", handlers.GetFilteredGeoJson)
	apiGroup.Get("/geojson/:name", handlers.GetGeoJsonByName)
}