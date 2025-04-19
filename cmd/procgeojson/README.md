# procgeojson

A command-line tool for processing location data from CSV files and converting it to GeoJSON format.

## Overview

The `procgeojson` tool reads location data from a CSV file, processes it using GeoJSON API calls, and saves the resulting GeoJSON data to an output directory. It's designed to simplify the workflow of converting location coordinates to rich GeoJSON representations.

## Installation

### Prerequisites

- Go 1.16 or higher
- Access to the required GeoJSON API (credentials stored in a secrets file)

### Building from Source

```bash
# Clone the repository
git clone https://github.com/harokku/logreason.git
cd logreason

# Build the procgeojson tool
go build -o procgeojson ./cmd/procgeojson
```

## Usage

```bash
./procgeojson [flags]
```

### Command-line Flags

| Flag | Description | Default Value |
|------|-------------|---------------|
| `-csv` | Path to the input CSV file | `/locations/input.csv` |
| `-range` | Range value for GeoJSON API calls (in seconds) | `600` |
| `-output` | Directory to save GeoJSON files | `/out/geojson` |
| `-secrets` | Path to the secrets file | `/config/secret.json` |

### Examples

Basic usage with default settings:
```bash
./procgeojson
```

Specifying custom file paths:
```bash
./procgeojson -csv ./my_locations.csv -output ./my_geojson_output -secrets ./my_secrets.json
```

Adjusting the range value:
```bash
./procgeojson -range 1200
```

## Input CSV Format

The input CSV file should contain location data with the following columns:
- Latitude
- Longitude
- (Additional columns as required by your implementation)

## Secrets File

The secrets file should be a JSON file containing the necessary API keys and credentials for accessing the GeoJSON API. The format depends on the specific requirements of your implementation.

## Error Handling

The tool provides detailed error messages for common issues:
- Missing input files
- CSV parsing errors
- API authentication failures
- Output directory access problems

If errors occur during processing, the tool will log warnings but continue processing other locations when possible.

## Troubleshooting

### Common Issues

1. **"Input file does not exist"**
   - Ensure the CSV file path is correct and the file exists

2. **"Error loading secrets"**
   - Verify the secrets file path and format

3. **"Error creating output directory"**
   - Check if the application has write permissions to create the specified output directory

4. **"No valid locations found in the CSV file"**
   - Ensure your CSV file contains valid location data in the expected format

## License

[Specify your license here]

## Contributing

[Add contribution guidelines if applicable]