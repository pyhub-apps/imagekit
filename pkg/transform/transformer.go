package transform

import (
	"bytes"
	"fmt"
	"io"
)

// Resize implements image resizing functionality
func (t *Transformer) Resize(input io.Reader, output io.Writer, options ResizeOptions) error {
	// Load the image
	img, format, err := LoadImage(input)
	if err != nil {
		return fmt.Errorf("failed to load image: %w", err)
	}
	
	// Get original dimensions
	info := GetImageInfo(img, format)
	
	// Calculate target dimensions
	targetWidth, targetHeight := CalculateDimensions(info.Width, info.Height, options)
	
	// Resize the image
	resizedImg, err := resizeImage(img, targetWidth, targetHeight, options.Mode)
	if err != nil {
		return fmt.Errorf("failed to resize image: %w", err)
	}
	
	// Save the resized image
	quality := options.Quality
	if quality <= 0 {
		quality = 95 // Default high quality
	}
	
	return SaveImage(output, resizedImg, format, quality)
}

// SetDPI implements DPI metadata setting functionality
func (t *Transformer) SetDPI(input io.Reader, output io.Writer, dpi int) error {
	// First detect the format
	img, format, err := LoadImage(input)
	if err != nil {
		return fmt.Errorf("failed to load image: %w", err)
	}
	
	// Save the image back to get raw data
	buf := &bytes.Buffer{}
	if err := SaveImage(buf, img, format, 95); err != nil {
		return fmt.Errorf("failed to encode image: %w", err)
	}
	
	// Process with DPI
	return ProcessImageWithDPI(bytes.NewReader(buf.Bytes()), output, format, dpi)
}

// RemoveWatermark implements watermark removal functionality
func (t *Transformer) RemoveWatermark(input io.Reader, output io.Writer, area Rectangle) error {
	// Load the image
	img, format, err := LoadImage(input)
	if err != nil {
		return fmt.Errorf("failed to load image: %w", err)
	}
	
	// Validate the rectangle
	info := GetImageInfo(img, format)
	if err := ValidateRectangle(area, info.Width, info.Height); err != nil {
		return fmt.Errorf("invalid watermark area: %w", err)
	}
	
	// For now, we'll implement the actual watermark removal in watermark.go
	// This is just the interface implementation
	processedImg, err := removeWatermarkFromArea(img, area)
	if err != nil {
		return fmt.Errorf("failed to remove watermark: %w", err)
	}
	
	// Save the processed image
	return SaveImage(output, processedImg, format, 95)
}