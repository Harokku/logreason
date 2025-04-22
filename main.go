package main

import (
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"logreason/internal/routes"
)

func main() {
	// Create a new Fiber app
	app := fiber.New(fiber.Config{
		AppName: "LogReason API",
	})

	// Add CORS middleware
	app.Use(cors.New())

	// Add logger middleware
	app.Use(logger.New())

	// Setup routes
	routes.SetupRoutes(app)

	// Start the server
	log.Println("Starting server on :3000")
	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
