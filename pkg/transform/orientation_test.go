package transform

import (
	"bytes"
	"image"
	"image/jpeg"
	"testing"
)

func TestLoadImageWithOrientation(t *testing.T) {
	// Create a simple test image
	img := image.NewRGBA(image.Rect(0, 0, 100, 200))
	
	// Encode it as JPEG
	buf := &bytes.Buffer{}
	if err := jpeg.Encode(buf, img, nil); err != nil {
		t.Fatalf("Failed to encode test image: %v", err)
	}
	
	// Load the image (this should handle EXIF orientation if present)
	loadedImg, format, err := LoadImage(buf)
	if err != nil {
		t.Fatalf("Failed to load image: %v", err)
	}
	
	if format != FormatJPEG {
		t.Errorf("Expected format JPEG, got %v", format)
	}
	
	// Check dimensions
	bounds := loadedImg.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y
	
	if width != 100 || height != 200 {
		t.Errorf("Unexpected dimensions: %dx%d, expected 100x200", width, height)
	}
}

func TestLoadImageFormats(t *testing.T) {
	tests := []struct {
		name   string
		format ImageFormat
	}{
		{"JPEG", FormatJPEG},
		{"PNG", FormatPNG},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test image
			img := image.NewRGBA(image.Rect(0, 0, 50, 50))
			buf := &bytes.Buffer{}
			
			// Encode based on format
			switch tt.format {
			case FormatJPEG:
				if err := jpeg.Encode(buf, img, nil); err != nil {
					t.Fatalf("Failed to encode JPEG: %v", err)
				}
			case FormatPNG:
				if err := SaveImage(buf, img, FormatPNG, 0); err != nil {
					t.Fatalf("Failed to encode PNG: %v", err)
				}
			}
			
			// Load and verify
			_, format, err := LoadImage(buf)
			if err != nil {
				t.Fatalf("Failed to load image: %v", err)
			}
			
			if format != tt.format {
				t.Errorf("Expected format %v, got %v", tt.format, format)
			}
		})
	}
}