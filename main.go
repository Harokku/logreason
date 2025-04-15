package main

import (
	"fmt"
	"log"
	"os"

	"logreason/internal/secrets"
)

func main() {
	// Create a new secrets manager
	secretsManager := secrets.NewManager()

	// Example 1: Load from environment variable
	// You can set this in your terminal with: export API_ENDPOINT="https://api.example.com/v1"
	if secretsManager.LoadFromEnvVar("API_ENDPOINT") {
		fmt.Println("Loaded API_ENDPOINT from environment variable")
	} else {
		fmt.Println("API_ENDPOINT not found in environment variables")
	}

	// Example 2: Load from .env file
	// Create a sample .env file if it doesn't exist
	envFilePath := "config/.env"
	if _, err := os.Stat(envFilePath); os.IsNotExist(err) {
		envContent := `# Sample .env file
API_KEY=your-api-key-here
DATABASE_URL="postgres://user:password@localhost:5432/dbname"
COMPLEX_URL="https://api.example.com/v2/resource?param1=value1&param2=value2"
`
		if err := os.WriteFile(envFilePath, []byte(envContent), 0600); err != nil {
			log.Fatalf("Failed to create sample .env file: %v", err)
		}
		fmt.Printf("Created sample .env file at %s\n", envFilePath)
	}

	// Load secrets from .env file
	if err := secretsManager.LoadFromDotEnvFile(envFilePath); err != nil {
		log.Fatalf("Failed to load secrets from .env file: %v", err)
	}
	fmt.Println("Loaded secrets from .env file")

	// Example 3: Load from JSON file
	if err := secretsManager.LoadFromFile("config/secret.json"); err != nil {
		fmt.Printf("Note: Failed to load secrets from JSON file: %v\n", err)
	} else {
		fmt.Println("Loaded secrets from JSON file")
	}

	// Display all loaded secrets
	fmt.Println("\nLoaded Secrets:")
	for key, value := range secretsManager.GetAll() {
		fmt.Printf("%s: %s\n", key, value)
	}

	// Example of how to use a secret
	if apiEndpoint, exists := secretsManager.Get("API_ENDPOINT"); exists {
		fmt.Printf("\nUsing API endpoint: %s\n", apiEndpoint)
	}

	if complexURL, exists := secretsManager.Get("COMPLEX_URL"); exists {
		fmt.Printf("Using complex URL: %s\n", complexURL)
	}
}
