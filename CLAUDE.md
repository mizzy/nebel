# Nebel Static Site Generator

## Project Overview
Nebel is a static site generator written in Go that converts Markdown files into HTML websites.

## Build Commands
```bash
# Build the project
go build -o nebel cmd/nebel/main.go

# Run tests
go test ./...

# Generate a new site
go run cmd/nebel/main.go new [site-name]

# Build/generate the static site
go run cmd/nebel/main.go generate
```

## Project Structure
```
.
├── cmd/
│   └── nebel/
│       └── main.go       # Main entry point
├── generate.go           # Site generation logic
├── new.go               # New site creation logic
├── go.mod               # Go module definition
└── go.sum               # Go dependencies checksums
```

## Key Features
- Markdown to HTML conversion
- Static file copying from "static" to "public" directories
- Site scaffolding with `new` command
- HTML formatting and generation

## Development Notes
- The project uses Go modules for dependency management
- Main commands are `new` and `generate`
- Static files are copied from a "static" directory to "public" during generation
- HTML output is formatted for readability