package csvparser

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

// Location represents a location with a name, latitude, and longitude
type Location struct {
	Name      string  `json:"name"`
	City      string  `json:"city,omitempty"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// ParseError represents an error that occurred during parsing
type ParseError struct {
	Row     int
	Column  int
	Message string
}

func (e ParseError) Error() string {
	return fmt.Sprintf("error at row %d, column %d: %s", e.Row, e.Column, e.Message)
}

// ParseResult represents the result of parsing a CSV file
type ParseResult struct {
	Locations []Location
	Errors    []ParseError
	Success   bool
}

// Parser defines the interface for CSV parsers
type Parser interface {
	Parse(reader io.Reader) ParseResult
	ParseFile(filePath string) ParseResult
	UpdateFile(filePath string, locations []Location) error
}

// DefaultParser is the default implementation of Parser
type DefaultParser struct{}

// NewParser creates a new DefaultParser
func NewParser() Parser {
	return &DefaultParser{}
}

// Parse parses a CSV from an io.Reader
func (p *DefaultParser) Parse(reader io.Reader) ParseResult {
	csvReader := csv.NewReader(reader)

	// Read header
	header, err := csvReader.Read()
	if err != nil {
		return ParseResult{
			Success: false,
			Errors: []ParseError{
				{Row: 0, Column: 0, Message: fmt.Sprintf("failed to read header: %v", err)},
			},
		}
	}

	// Validate header
	if len(header) < 3 {
		return ParseResult{
			Success: false,
			Errors: []ParseError{
				{Row: 0, Column: 0, Message: "header must contain at least 3 columns"},
			},
		}
	}

	// Process rows
	var locations []Location
	var errors []ParseError
	rowNum := 1 // Start from 1 because header is row 0

	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			errors = append(errors, ParseError{
				Row:     rowNum,
				Column:  0,
				Message: fmt.Sprintf("failed to read row: %v", err),
			})
			rowNum++
			continue
		}

		// Parse location
		location, parseErrors := parseLocation(row, rowNum)
		if len(parseErrors) > 0 {
			errors = append(errors, parseErrors...)
		} else {
			locations = append(locations, location)
		}

		rowNum++
	}

	return ParseResult{
		Locations: locations,
		Errors:    errors,
		Success:   len(errors) == 0,
	}
}

// ParseFile parses a CSV file
func (p *DefaultParser) ParseFile(filePath string) ParseResult {
	file, err := os.Open(filePath)
	if err != nil {
		return ParseResult{
			Success: false,
			Errors: []ParseError{
				{Row: 0, Column: 0, Message: fmt.Sprintf("failed to open file: %v", err)},
			},
		}
	}
	defer file.Close()

	return p.Parse(file)
}

// UpdateFile updates a CSV file with new location data
func (p *DefaultParser) UpdateFile(filePath string, locations []Location) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	err = writer.Write([]string{"STAZIONAMENTO", "LAT", "LON"})
	if err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write locations
	for _, loc := range locations {
		var name string
		if loc.City != "" {
			name = fmt.Sprintf("%s (%s)", loc.Name, loc.City)
		} else {
			name = loc.Name
		}

		err = writer.Write([]string{
			name,
			fmt.Sprintf("%.5f", loc.Latitude),
			fmt.Sprintf("%.5f", loc.Longitude),
		})
		if err != nil {
			return fmt.Errorf("failed to write location: %w", err)
		}
	}

	return nil
}

// parseLocation parses a location from a CSV row
func parseLocation(row []string, rowNum int) (Location, []ParseError) {
	var errors []ParseError

	if len(row) < 3 {
		errors = append(errors, ParseError{
			Row:     rowNum,
			Column:  0,
			Message: "row must contain at least 3 columns",
		})
		return Location{}, errors
	}

	// Parse name and city
	name, city := parseNameAndCity(row[0])

	// Parse latitude
	lat, err := parseFloat(row[1])
	if err != nil {
		errors = append(errors, ParseError{
			Row:     rowNum,
			Column:  1,
			Message: fmt.Sprintf("invalid latitude: %v", err),
		})
	}

	// Parse longitude
	lon, err := parseFloat(row[2])
	if err != nil {
		errors = append(errors, ParseError{
			Row:     rowNum,
			Column:  2,
			Message: fmt.Sprintf("invalid longitude: %v", err),
		})
	}

	if len(errors) > 0 {
		return Location{}, errors
	}

	return Location{
		Name:      name,
		City:      city,
		Latitude:  lat,
		Longitude: lon,
	}, nil
}

// parseNameAndCity parses a name and city from a string like "NAME (CITY)"
func parseNameAndCity(s string) (string, string) {
	s = strings.TrimSpace(s)

	// Check if the string contains a city in parentheses
	openParen := strings.Index(s, "(")
	closeParen := strings.LastIndex(s, ")")

	// Only extract city if the pattern is exactly "NAME (CITY)" with nothing after the closing parenthesis
	if openParen != -1 && closeParen != -1 && openParen < closeParen && closeParen == len(s)-1 {
		name := strings.TrimSpace(s[:openParen])
		city := strings.TrimSpace(s[openParen+1 : closeParen])
		return name, city
	}

	return s, ""
}

// parseFloat parses a float from a string
func parseFloat(s string) (float64, error) {
	var result float64
	_, err := fmt.Sscanf(strings.TrimSpace(s), "%f", &result)
	return result, err
}
