# Nebel Static Site Generator

## Project Overview

Nebel is a static site generator written in Go that converts Markdown
files into HTML websites.

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

# Serve public/ and regenerate on changes to posts/, layouts/, static/
go run cmd/nebel/main.go serve
```

## Project Structure

```text
.
├── cmd/
│   └── nebel/
│       └── main.go       # Main entry point
├── generate.go           # Site generation logic
├── serve.go              # Local server with auto-regeneration
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
- Main commands are `new`, `generate` and `serve`
- Static files are copied from a "static" directory to "public" during generation
- HTML output is formatted for readability

## Git Workflow

- **Never commit directly to main branch**
- Always create a feature branch and submit a Pull Request
- PRs should be reviewed before merging
