// Package secrets provides secure handling of API keys and other sensitive information.
// It supports loading secrets from environment variables or from encrypted files.
package secrets

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Manager handles the loading and retrieval of secrets
type Manager struct {
	secrets map[string]string
	mu      sync.RWMutex
}

// NewManager creates a new secrets manager
func NewManager() *Manager {
	return &Manager{
		secrets: make(map[string]string),
	}
}

// LoadFromEnv loads secrets from environment variables with the given prefix
// For example, if prefix is "API_KEY_", then environment variables like
// API_KEY_SERVICE1, API_KEY_SERVICE2 will be loaded as secrets with keys
// "SERVICE1", "SERVICE2"
func (m *Manager) LoadFromEnv(prefix string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, env := range os.Environ() {
		if !strings.HasPrefix(env, prefix) {
			continue
		}

		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimPrefix(parts[0], prefix)
		value := parts[1]

		m.secrets[key] = value
	}

	return nil
}

// LoadFromFile loads secrets from a JSON file
// The file should contain a JSON object where keys are secret names and values are secret values
func (m *Manager) LoadFromFile(filePath string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Ensure the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("secrets file does not exist: %s", filePath)
	}

	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read secrets file: %w", err)
	}

	// Parse the JSON
	var fileSecrets map[string]string
	if err := json.Unmarshal(data, &fileSecrets); err != nil {
		return fmt.Errorf("failed to parse secrets file: %w", err)
	}

	// Add the secrets to our map
	for k, v := range fileSecrets {
		m.secrets[k] = v
	}

	return nil
}

// Get retrieves a secret by its key
// Returns the secret value and a boolean indicating if the secret was found
func (m *Manager) Get(key string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	value, exists := m.secrets[key]
	return value, exists
}

// GetOrDefault retrieves a secret by its key, returning a default value if not found
func (m *Manager) GetOrDefault(key, defaultValue string) string {
	value, exists := m.Get(key)
	if !exists {
		return defaultValue
	}
	return value
}

// Set adds or updates a secret
func (m *Manager) Set(key, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.secrets[key] = value
}

// GetAll returns a copy of all secrets
func (m *Manager) GetAll() map[string]string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Create a copy to avoid exposing the internal map
	result := make(map[string]string, len(m.secrets))
	for k, v := range m.secrets {
		result[k] = v
	}

	return result
}

// LoadFromEnvVar loads a single environment variable as a secret
// If the environment variable doesn't exist, it returns false
func (m *Manager) LoadFromEnvVar(key string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	value, exists := os.LookupEnv(key)
	if !exists {
		return false
	}

	m.secrets[key] = value
	return true
}

// LoadFromDotEnvFile loads secrets from a .env file
// The file should contain lines in the format KEY=VALUE
func (m *Manager) LoadFromDotEnvFile(filePath string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Ensure the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("env file does not exist: %s", filePath)
	}

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open env file: %w", err)
	}
	defer file.Close()

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split the line into key and value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		if len(value) > 1 && (value[0] == '"' || value[0] == '\'') && value[0] == value[len(value)-1] {
			value = value[1 : len(value)-1]
		}

		m.secrets[key] = value
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading env file: %w", err)
	}

	return nil
}

// SaveToFile saves the current secrets to a JSON file
func (m *Manager) SaveToFile(filePath string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Create the directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Marshal the secrets to JSON
	data, err := json.MarshalIndent(m.secrets, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal secrets: %w", err)
	}

	// Write to the file
	if err := os.WriteFile(filePath, data, 0600); err != nil {
		return fmt.Errorf("failed to write secrets file: %w", err)
	}

	return nil
}
