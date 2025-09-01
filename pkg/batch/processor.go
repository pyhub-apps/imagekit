package batch

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	
	"github.com/allieus/imagekit/pkg/transform"
)

// Processor handles batch image processing
type Processor struct {
	transformer *transform.Transformer
	verbose     bool
}

// NewProcessor creates a new batch processor
func NewProcessor(transformer *transform.Transformer) *Processor {
	return &Processor{
		transformer: transformer,
		verbose:     true,
	}
}

// ProcessFiles processes multiple files matching the pattern
func (p *Processor) ProcessFiles(pattern string, options ProcessOptions, progressCallback func(current int, total int, fileName string, success bool)) (*BatchResult, error) {
	// Find matching files
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid glob pattern: %w", err)
	}
	
	if len(matches) == 0 {
		return nil, fmt.Errorf("no files matching pattern: %s", pattern)
	}
	
	// Filter out already converted files and non-image files
	var filesToProcess []string
	for _, match := range matches {
		// Skip if it's already a converted file
		if IsConvertedFile(match) {
			continue
		}
		
		// Check if it's a supported image format
		ext := strings.ToLower(filepath.Ext(match))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			continue
		}
		
		filesToProcess = append(filesToProcess, match)
	}
	
	if len(filesToProcess) == 0 {
		return nil, fmt.Errorf("no valid image files found to process")
	}
	
	// Process files
	result := &BatchResult{
		TotalFiles:   len(filesToProcess),
		SuccessCount: 0,
		FailedFiles:  []FailedFile{},
	}
	
	for i, inputPath := range filesToProcess {
		outputPath := GenerateOutputPath(inputPath)
		
		// Process single file
		err := p.processSingleFile(inputPath, outputPath, options)
		
		if err != nil {
			result.FailedFiles = append(result.FailedFiles, FailedFile{
				Path:  inputPath,
				Error: err,
			})
			if progressCallback != nil {
				progressCallback(i+1, result.TotalFiles, filepath.Base(inputPath), false)
			}
		} else {
			result.SuccessCount++
			if progressCallback != nil {
				progressCallback(i+1, result.TotalFiles, filepath.Base(inputPath), true)
			}
		}
	}
	
	return result, nil
}

// ProcessSingleFile processes a single file (public for single file mode)
func (p *Processor) ProcessSingleFile(inputPath, outputPath string, options ProcessOptions) error {
	return p.processSingleFile(inputPath, outputPath, options)
}

// processSingleFile handles the actual file processing
func (p *Processor) processSingleFile(inputPath, outputPath string, options ProcessOptions) error {
	// Check if input file exists
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return fmt.Errorf("input file does not exist: %s", inputPath)
	}
	
	// Process based on options
	if options.ResizeOptions != nil {
		// Open input file for resize
		inputFile, err := os.Open(inputPath)
		if err != nil {
			return fmt.Errorf("failed to open input file: %w", err)
		}
		defer func() { _ = inputFile.Close() }()
		
		// Create output file
		outputFile, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer func() { _ = outputFile.Close() }()
		
		// Resize operation
		if err := p.transformer.Resize(inputFile, outputFile, *options.ResizeOptions); err != nil {
			_ = os.Remove(outputPath) // Clean up on failure
			return fmt.Errorf("resize failed: %w", err)
		}
		
		// If DPI is also specified, process it
		if options.DPI > 0 {
			// Close files first
			_ = inputFile.Close()
			_ = outputFile.Close()
			
			// Use the resized output as input for DPI change
			tempFile, err := os.Open(outputPath)
			if err != nil {
				return fmt.Errorf("failed to open temp file: %w", err)
			}
			defer func() { _ = tempFile.Close() }()
			
			tempOutputPath := outputPath + ".tmp"
			tempOutput, err := os.Create(tempOutputPath)
			if err != nil {
				return fmt.Errorf("failed to create temp output: %w", err)
			}
			defer func() { _ = tempOutput.Close() }()
			
			if err := p.transformer.SetDPI(tempFile, tempOutput, options.DPI); err != nil {
				_ = os.Remove(tempOutputPath)
				return fmt.Errorf("DPI setting failed: %w", err)
			}
			
			// Replace original with DPI-adjusted version
			_ = tempFile.Close()
			_ = tempOutput.Close()
			if err := os.Rename(tempOutputPath, outputPath); err != nil {
				return fmt.Errorf("failed to replace file: %w", err)
			}
		}
	} else if options.DPI > 0 {
		// DPI only operation
		inputFile, err := os.Open(inputPath)
		if err != nil {
			return fmt.Errorf("failed to open input file: %w", err)
		}
		defer func() { _ = inputFile.Close() }()
		
		outputFile, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer func() { _ = outputFile.Close() }()
		
		if err := p.transformer.SetDPI(inputFile, outputFile, options.DPI); err != nil {
			_ = os.Remove(outputPath) // Clean up on failure
			return fmt.Errorf("DPI setting failed: %w", err)
		}
	} else {
		return fmt.Errorf("no conversion options specified")
	}
	
	return nil
}