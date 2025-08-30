# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based CLI tool for image transformation, designed for optimizing images with a focus on single-file processing of JPG and PNG images with size/DPI conversion and edge cropping capabilities.

## Key Requirements

### Supported Features
- **Image Formats**: JPG and PNG only (no GIF support)
- **Single File Processing**: Process one image at a time (no batch processing)
- **Size Conversion**: Resize to specified pixel dimensions or ratios
- **DPI Conversion**: Convert to recommended DPI (72dpi or 96dpi)
- **Edge Cropping**: Remove edges from images for margin removal
- **Library Design**: Core logic should be modular for future web service integration

### Technical Constraints
- Must minimize image quality loss during conversion
- Cross-platform support (Windows, macOS, Linux)
- CLI-only interface (no GUI)

## Development Commands

### Go Module Setup
```bash
# Initialize Go module
go mod init github.com/allieus/pyhub-imagekit

# Add dependencies
go get github.com/disintegration/imaging  # For image processing
go get github.com/spf13/cobra             # For CLI framework
```

### Build & Run
```bash
# Build the application
go build -o imagekit cmd/imagekit/main.go

# Run the application
./imagekit [command] [flags]

# Build for different platforms
GOOS=windows GOARCH=amd64 go build -o imagekit.exe cmd/imagekit/main.go
GOOS=darwin GOARCH=amd64 go build -o imagekit-mac cmd/imagekit/main.go
GOOS=linux GOARCH=amd64 go build -o imagekit-linux cmd/imagekit/main.go
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./pkg/transform/...

# Run with verbose output
go test -v ./...
```

## Project Structure

```
.
├── cmd/
│   └── imagekit/         # CLI application entry point
│       └── main.go
├── pkg/
│   ├── transform/        # Core image transformation library
│   │   ├── resize.go     # Size/resolution transformations
│   │   ├── dpi.go        # DPI conversion logic
│   │   └── crop.go       # Edge cropping logic
│   └── cli/              # CLI command implementations
│       ├── convert.go    # Convert command
│       └── root.go       # Root command setup
├── internal/             # Internal packages (not for library export)
│   └── utils/            # Utility functions
├── sample/               # Sample images for testing
├── go.mod
└── go.sum
```

## Architecture Guidelines

### Library Design
The core image processing logic should be in `pkg/transform/` as a reusable library with clear interfaces:

```go
// Example interface for the transform package
type ImageTransformer interface {
    Resize(input io.Reader, width, height int) (io.Reader, error)
    SetDPI(input io.Reader, dpi int) (io.Reader, error)
    CropEdges(input io.Reader, top, right, bottom, left int) (io.Reader, error)
}
```

### CLI Command Structure
Using Cobra, structure commands as:
```
imagekit convert --size=1920x1080 --dpi=96 input.jpg output.jpg
imagekit crop --top=10 --bottom=10 --left=10 --right=10 input.jpg output.jpg
```

### Error Handling
- Return meaningful error messages for user-facing CLI
- Use wrapped errors with context: `fmt.Errorf("failed to resize image: %w", err)`
- Validate input files exist and are valid JPG/PNG before processing

### Testing Strategy
- Unit tests for each transformation function
- Integration tests using sample images in `sample/` directory
- Benchmark tests for performance-critical operations