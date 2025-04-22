// Package api provides API documentation and related utilities
package api

// Documentation is the API documentation string
const Documentation = `
# LogReason API Documentation

## Available Endpoints

### CSV Endpoints
- GET /api/locations/csv
  - Returns the locations/input.csv file as an attachment
  - Example: curl -O -J http://localhost:3000/api/locations/csv

- GET /api/locations/json
  - Returns the parsed content of locations/input.csv as a JSON array
  - Example: curl http://localhost:3000/api/locations/json

### GeoJSON Endpoints
- GET /api/geojson
  - Returns all GeoJSON files from out/geojson directory as a combined JSON array
  - Example: curl http://localhost:3000/api/geojson

- GET /api/geojson/:name
  - Returns a specific GeoJSON file by name (without .json extension)
  - Example: curl http://localhost:3000/api/geojson/APMPAD

- GET /api/geojson/filter?names=name1,name2,name3
  - Returns multiple specific GeoJSON files as a combined JSON array
  - Example: curl http://localhost:3000/api/geojson/filter?names=APMPAD,ARGLIM
`
