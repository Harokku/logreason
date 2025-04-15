# CSV Parser Package

This package provides functionality for reading, parsing, and manipulating CSV files containing location data.

## Features

- Read and parse CSV files with location data
- Convert CSV data to structured Go objects
- Handle errors gracefully with detailed reporting
- Update CSV files with modified data
- Recover from panics to prevent application crashes

## Usage

### Parsing a CSV File

```go
import "logreason/internal/csvparser"

// Create a new parser
parser := csvparser.NewParser()

// Parse a CSV file
result := parser.ParseFile("path/to/file.csv")

// Check if parsing was successful
if result.Success {
    // Use the parsed locations
    for _, location := range result.Locations {
        fmt.Printf("Name: %s, City: %s, Coordinates: %.5f, %.5f\n",
            location.Name, location.City, location.Latitude, location.Longitude)
    }
} else {
    // Handle errors
    for _, err := range result.Errors {
        fmt.Println(err.Error())
    }
}
```

### Parsing from a Reader

```go
import (
    "logreason/internal/csvparser"
    "strings"
)

// Create a CSV string
csvData := `STAZIONAMENTO,LAT,LON
APMPAD (PADERNO DUGNANO),45.57520,9.15325
ARGLIM (LIMBIATE),45.61493,9.12310`

// Create a reader
reader := strings.NewReader(csvData)

// Parse the CSV data
parser := csvparser.NewParser()
result := parser.Parse(reader)

// Process the result as shown above
```

### Updating a CSV File

```go
import "logreason/internal/csvparser"

// Create a new parser
parser := csvparser.NewParser()

// Parse a CSV file
result := parser.ParseFile("path/to/file.csv")

// Modify the locations
if len(result.Locations) > 0 {
    result.Locations[0].City = "NEW CITY"
}

// Update the file with the modified locations
err := parser.UpdateFile("path/to/file.csv", result.Locations)
if err != nil {
    fmt.Printf("Error updating file: %v\n", err)
}
```

## Design Decisions

### Error Handling

The package uses a structured approach to error handling:

- Errors are collected rather than causing immediate failure
- Each error includes the row and column where it occurred
- The `ParseResult` struct includes both successful results and errors
- The `Success` flag indicates whether parsing was completely successful

### Data Structure

The `Location` struct represents a location with:

- Name: The name of the location
- City: The city where the location is (optional)
- Latitude: The latitude coordinate
- Longitude: The longitude coordinate

### Encoder/Decoder Pattern

The package follows an encoder/decoder pattern:

- The `Parser` interface defines methods for parsing and updating CSV files
- The `DefaultParser` implements this interface
- Helper functions handle specific parsing tasks
- This design makes it easy to reason about, debug, and modify the system

### Recovery from Panics

The package includes a test demonstrating how to recover from panics:

- Use defer and recover to catch panics
- Convert panics to structured errors
- Return a meaningful error message

## Future Improvements

- Add support for different CSV formats
- Implement validation for location data
- Add support for batch processing of multiple files
- Implement concurrent processing for large files