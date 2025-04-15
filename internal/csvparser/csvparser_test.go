package csvparser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseNameAndCity(t *testing.T) {
	tests := []struct {
		input      string
		expectName string
		expectCity string
	}{
		{"APMPAD (PADERNO DUGNANO)", "APMPAD", "PADERNO DUGNANO"},
		{"ARGLIM (LIMBIATE)", "ARGLIM", "LIMBIATE"},
		{"SIMPLE", "SIMPLE", ""},
		{"", "", ""},
		{"(CITY)", "", "CITY"},
		{"NAME (CITY) EXTRA", "NAME (CITY) EXTRA", ""},
		{"NAME (CITY", "NAME (CITY", ""},
		{"NAME CITY)", "NAME CITY)", ""},
	}

	for _, tc := range tests {
		name, city := parseNameAndCity(tc.input)
		if name != tc.expectName {
			t.Errorf("parseNameAndCity(%q) name = %q, want %q", tc.input, name, tc.expectName)
		}
		if city != tc.expectCity {
			t.Errorf("parseNameAndCity(%q) city = %q, want %q", tc.input, city, tc.expectCity)
		}
	}
}

func TestParseFloat(t *testing.T) {
	tests := []struct {
		input       string
		expectValue float64
		expectError bool
	}{
		{"45.57520", 45.57520, false},
		{"9.15325", 9.15325, false},
		{" 45.57520 ", 45.57520, false},
		{"", 0, true},
		{"invalid", 0, true},
	}

	for _, tc := range tests {
		value, err := parseFloat(tc.input)
		if (err != nil) != tc.expectError {
			t.Errorf("parseFloat(%q) error = %v, wantError = %v", tc.input, err, tc.expectError)
		}
		if !tc.expectError && value != tc.expectValue {
			t.Errorf("parseFloat(%q) = %v, want %v", tc.input, value, tc.expectValue)
		}
	}
}

func TestParseLocation(t *testing.T) {
	tests := []struct {
		name        string
		row         []string
		rowNum      int
		expectLoc   Location
		expectError bool
	}{
		{
			name:   "valid location",
			row:    []string{"APMPAD (PADERNO DUGNANO)", "45.57520", "9.15325"},
			rowNum: 1,
			expectLoc: Location{
				Name:      "APMPAD",
				City:      "PADERNO DUGNANO",
				Latitude:  45.57520,
				Longitude: 9.15325,
			},
			expectError: false,
		},
		{
			name:        "too few columns",
			row:         []string{"APMPAD (PADERNO DUGNANO)", "45.57520"},
			rowNum:      2,
			expectLoc:   Location{},
			expectError: true,
		},
		{
			name:        "invalid latitude",
			row:         []string{"APMPAD (PADERNO DUGNANO)", "invalid", "9.15325"},
			rowNum:      3,
			expectLoc:   Location{},
			expectError: true,
		},
		{
			name:        "invalid longitude",
			row:         []string{"APMPAD (PADERNO DUGNANO)", "45.57520", "invalid"},
			rowNum:      4,
			expectLoc:   Location{},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			loc, errs := parseLocation(tc.row, tc.rowNum)
			if (len(errs) > 0) != tc.expectError {
				t.Errorf("parseLocation() error = %v, wantError = %v", errs, tc.expectError)
			}
			if !tc.expectError {
				if loc.Name != tc.expectLoc.Name {
					t.Errorf("parseLocation() name = %q, want %q", loc.Name, tc.expectLoc.Name)
				}
				if loc.City != tc.expectLoc.City {
					t.Errorf("parseLocation() city = %q, want %q", loc.City, tc.expectLoc.City)
				}
				if loc.Latitude != tc.expectLoc.Latitude {
					t.Errorf("parseLocation() latitude = %v, want %v", loc.Latitude, tc.expectLoc.Latitude)
				}
				if loc.Longitude != tc.expectLoc.Longitude {
					t.Errorf("parseLocation() longitude = %v, want %v", loc.Longitude, tc.expectLoc.Longitude)
				}
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name           string
		csv            string
		expectSuccess  bool
		expectLocCount int
		expectErrCount int
	}{
		{
			name: "valid csv",
			csv: `STAZIONAMENTO,LAT,LON
APMPAD (PADERNO DUGNANO),45.57520,9.15325
ARGLIM (LIMBIATE),45.61493,9.12310`,
			expectSuccess:  true,
			expectLocCount: 2,
			expectErrCount: 0,
		},
		{
			name:           "invalid header",
			csv:            `STAZIONAMENTO,LAT`,
			expectSuccess:  false,
			expectLocCount: 0,
			expectErrCount: 1,
		},
		{
			name: "invalid row",
			csv: `STAZIONAMENTO,LAT,LON
APMPAD (PADERNO DUGNANO),invalid,9.15325`,
			expectSuccess:  false,
			expectLocCount: 0,
			expectErrCount: 1,
		},
		{
			name: "mixed valid and invalid rows",
			csv: `STAZIONAMENTO,LAT,LON
APMPAD (PADERNO DUGNANO),45.57520,9.15325
ARGLIM (LIMBIATE),invalid,9.12310`,
			expectSuccess:  false,
			expectLocCount: 1,
			expectErrCount: 1,
		},
	}

	parser := NewParser()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			reader := strings.NewReader(tc.csv)
			result := parser.Parse(reader)

			if result.Success != tc.expectSuccess {
				t.Errorf("Parse() success = %v, want %v", result.Success, tc.expectSuccess)
			}
			if len(result.Locations) != tc.expectLocCount {
				t.Errorf("Parse() location count = %v, want %v", len(result.Locations), tc.expectLocCount)
			}
			if len(result.Errors) != tc.expectErrCount {
				t.Errorf("Parse() error count = %v, want %v", len(result.Errors), tc.expectErrCount)
			}
		})
	}
}

func TestUpdateFile(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "csvparser_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test locations
	locations := []Location{
		{
			Name:      "APMPAD",
			City:      "PADERNO DUGNANO",
			Latitude:  45.57520,
			Longitude: 9.15325,
		},
		{
			Name:      "ARGLIM",
			City:      "LIMBIATE",
			Latitude:  45.61493,
			Longitude: 9.12310,
		},
	}

	// Create a test file path
	testFilePath := filepath.Join(tempDir, "test.csv")

	// Test UpdateFile
	parser := NewParser()
	err = parser.UpdateFile(testFilePath, locations)
	if err != nil {
		t.Fatalf("UpdateFile() error = %v", err)
	}

	// Read the file back and verify its contents
	data, err := os.ReadFile(testFilePath)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	expected := `STAZIONAMENTO,LAT,LON
APMPAD (PADERNO DUGNANO),45.57520,9.15325
ARGLIM (LIMBIATE),45.61493,9.12310
`

	if string(data) != expected {
		t.Errorf("UpdateFile() produced incorrect file content:\nGot:\n%s\nWant:\n%s", string(data), expected)
	}

	// Test ParseFile
	result := parser.ParseFile(testFilePath)
	if !result.Success {
		t.Errorf("ParseFile() success = %v, want true", result.Success)
	}
	if len(result.Locations) != 2 {
		t.Errorf("ParseFile() location count = %v, want 2", len(result.Locations))
	}
	if len(result.Errors) != 0 {
		t.Errorf("ParseFile() error count = %v, want 0", len(result.Errors))
	}
}

func TestRecoverFromPanic(t *testing.T) {
	// Create a function that will panic
	panicFunc := func() (result ParseResult) {
		// Set up recovery
		defer func() {
			if r := recover(); r != nil {
				result = ParseResult{
					Success: false,
					Errors: []ParseError{
						{Row: 0, Column: 0, Message: fmt.Sprintf("recovered from panic: %v", r)},
					},
				}
			}
		}()

		// This will panic
		var nilSlice []string
		_ = nilSlice[0] // This will cause a panic

		// This won't be reached
		return ParseResult{Success: true}
	}

	// Call the function and check the result
	result := panicFunc()
	if result.Success {
		t.Errorf("Expected failure after panic, got success")
	}
	if len(result.Errors) != 1 {
		t.Errorf("Expected 1 error after panic, got %d", len(result.Errors))
	}
}
