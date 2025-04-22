// Package geojson provides functionality for fetching and saving GeoJSON data from the Geoapify API.
package geojson

import (
	"fmt"
	"io"
	utilities "logreason/internal/utils"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"logreason/internal/csvparser"
	"logreason/internal/secrets"
)

// Constants for the Geoapify API
const (
	// DefaultOutputDir Default directory for saving GeoJSON files
	DefaultOutputDir = "out/geojson"
)

// Manager handles fetching and saving GeoJSON data
type Manager struct {
	secretsManager *secrets.Manager
	apiKey         string
	baseURL        string
	outputDir      string
}

// NewManager creates a new GeoJSON manager
func NewManager(secretsManager *secrets.Manager) (*Manager, error) {
	// Get the API key and base URL from the secrets manager
	apiKey, exists := secretsManager.Get("GEOAPIFY_API_KEY")
	if !exists {
		return nil, fmt.Errorf("GEOAPIFY_API_KEY not found in secrets")
	}

	baseURL, exists := secretsManager.Get("GEOAPIFY_BASE_URL")
	if !exists {
		return nil, fmt.Errorf("GEOAPIFY_BASE_URL not found in secrets")
	}

	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(DefaultOutputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	return &Manager{
		secretsManager: secretsManager,
		apiKey:         apiKey,
		baseURL:        baseURL,
		outputDir:      DefaultOutputDir,
	}, nil
}

// SetOutputDir sets a custom output directory
func (m *Manager) SetOutputDir(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	m.outputDir = dir
	return nil
}

// FetchAndSaveGeoJSON fetches GeoJSON data for a location and saves it to a file
func (m *Manager) FetchAndSaveGeoJSON(location csvparser.Location, rangeValue int) error {
	// Build the URL with the location's coordinates, range, and API key
	url := strings.ReplaceAll(m.baseURL, "{LAT}", fmt.Sprintf("%f", location.Latitude))
	url = strings.ReplaceAll(url, "{LON}", fmt.Sprintf("%f", location.Longitude))
	url = strings.ReplaceAll(url, "{RANGE}", fmt.Sprintf("%d", rangeValue))
	url = strings.ReplaceAll(url, "{API}", m.apiKey)

	// Fetch the GeoJSON data
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch GeoJSON data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Create the output file
	//filename := fmt.Sprintf("%s.json", location.Name)
	//filePath := filepath.Join(m.outputDir, filename)
	filePath := m.getOutputFileName(&location)

	// Write the GeoJSON data to the file
	if err := os.WriteFile(filePath, body, 0644); err != nil {
		return fmt.Errorf("failed to write GeoJSON file: %w", err)
	}

	return nil
}

func (m *Manager) getOutputFileName(location *csvparser.Location) string {
	// Extract the station code (assuming it's before the parentheses)
	stationCode := location.Name

	// Convert city name to pascal case
	cityNamePascal := utilities.ToCamelCase(location.City)

	// Construct filename: STATIONCODE-CityName.json
	filename := fmt.Sprintf("%s-%s.json", stationCode, cityNamePascal)

	return filepath.Join(m.outputDir, filename)
}

// ProcessLocations processes all locations and saves their GeoJSON data
func (m *Manager) ProcessLocations(locations []csvparser.Location, rangeValue int) []error {
	var errors []error

	for _, location := range locations {
		if err := m.FetchAndSaveGeoJSON(location, rangeValue); err != nil {
			errors = append(errors, fmt.Errorf("error processing location %s: %w", location.Name, err))
		}
	}

	return errors
}
