package transform

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"testing"
)

func TestProcessImageWithDPI(t *testing.T) {
	// Create test images
	jpegImg := image.NewRGBA(image.Rect(0, 0, 100, 100))
	pngImg := image.NewRGBA(image.Rect(0, 0, 100, 100))
	
	tests := []struct {
		name   string
		format ImageFormat
		dpi    int
		img    image.Image
	}{
		{
			name:   "JPEG with 72 DPI",
			format: FormatJPEG,
			dpi:    72,
			img:    jpegImg,
		},
		{
			name:   "JPEG with 96 DPI",
			format: FormatJPEG,
			dpi:    96,
			img:    jpegImg,
		},
		{
			name:   "JPEG with 300 DPI",
			format: FormatJPEG,
			dpi:    300,
			img:    jpegImg,
		},
		{
			name:   "PNG with 72 DPI",
			format: FormatPNG,
			dpi:    72,
			img:    pngImg,
		},
		{
			name:   "PNG with 96 DPI",
			format: FormatPNG,
			dpi:    96,
			img:    pngImg,
		},
		{
			name:   "PNG with 300 DPI",
			format: FormatPNG,
			dpi:    300,
			img:    pngImg,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode image to bytes
			inputBuf := &bytes.Buffer{}
			switch tt.format {
			case FormatJPEG:
				if err := jpeg.Encode(inputBuf, tt.img, nil); err != nil {
					t.Fatalf("Failed to encode JPEG: %v", err)
				}
			case FormatPNG:
				if err := png.Encode(inputBuf, tt.img); err != nil {
					t.Fatalf("Failed to encode PNG: %v", err)
				}
			}
			
			// Process with DPI
			outputBuf := &bytes.Buffer{}
			err := ProcessImageWithDPI(bytes.NewReader(inputBuf.Bytes()), outputBuf, tt.format, tt.dpi)
			
			if err != nil {
				t.Errorf("ProcessImageWithDPI() error = %v", err)
			}
			
			// Verify output is not empty
			if outputBuf.Len() == 0 {
				t.Errorf("ProcessImageWithDPI() produced empty output")
			}
			
			// For JPEG, verify the output is valid
			if tt.format == FormatJPEG {
				_, err := jpeg.Decode(outputBuf)
				if err != nil {
					t.Errorf("ProcessImageWithDPI() produced invalid JPEG: %v", err)
				}
			}
			
			// For PNG, verify the output is valid
			if tt.format == FormatPNG {
				_, err := png.Decode(outputBuf)
				if err != nil {
					t.Errorf("ProcessImageWithDPI() produced invalid PNG: %v", err)
				}
			}
		})
	}
}

func TestSetJPEGDPI(t *testing.T) {
	tests := []struct {
		name    string
		dpi     int
		wantErr bool
	}{
		{
			name:    "Valid DPI 72",
			dpi:     72,
			wantErr: false,
		},
		{
			name:    "Valid DPI 96",
			dpi:     96,
			wantErr: false,
		},
		{
			name:    "Valid DPI 150",
			dpi:     150,
			wantErr: false,
		},
		{
			name:    "Valid DPI 300",
			dpi:     300,
			wantErr: false,
		},
		{
			name:    "Valid DPI 600",
			dpi:     600,
			wantErr: false,
		},
		{
			name:    "High DPI 1200",
			dpi:     1200,
			wantErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a JPEG image
			img := image.NewRGBA(image.Rect(0, 0, 100, 100))
			inputBuf := &bytes.Buffer{}
			if err := jpeg.Encode(inputBuf, img, nil); err != nil {
				t.Fatalf("Failed to encode JPEG: %v", err)
			}
			
			// Set DPI using SetJPEGDPI (exported function)
			result, err := SetJPEGDPI(inputBuf.Bytes(), tt.dpi)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("SetJPEGDPI() error = %v, wantErr %v", err, tt.wantErr)
			}
			
			if !tt.wantErr && len(result) == 0 {
				t.Errorf("SetJPEGDPI() produced empty output")
			}
		})
	}
}

// TestDPIConstants tests the DPI constants
func TestDPIConstants(t *testing.T) {
	tests := []struct {
		name string
		dpi  int
		want int
	}{
		{"DPI72", DPI72, 72},
		{"DPI96", DPI96, 96},
		{"DPI150", DPI150, 150},
		{"DPI300", DPI300, 300},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.dpi != tt.want {
				t.Errorf("DPI constant = %d, want %d", tt.dpi, tt.want)
			}
		})
	}
}

func TestProcessImageWithDPIEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		format  ImageFormat
		dpi     int
		wantErr bool
	}{
		{
			name:    "Zero DPI",
			format:  FormatJPEG,
			dpi:     0,
			wantErr: false, // Should handle gracefully
		},
		{
			name:    "Negative DPI",
			format:  FormatJPEG,
			dpi:     -100,
			wantErr: false, // Should handle gracefully
		},
		{
			name:    "Very high DPI",
			format:  FormatPNG,
			dpi:     10000,
			wantErr: false,
		},
		{
			name:    "Invalid format",
			format:  ImageFormat("invalid"),
			dpi:     96,
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a simple image
			img := image.NewRGBA(image.Rect(0, 0, 10, 10))
			inputBuf := &bytes.Buffer{}
			
			// Encode based on format (for valid formats)
			if tt.format == FormatJPEG {
				jpeg.Encode(inputBuf, img, nil)
			} else if tt.format == FormatPNG {
				png.Encode(inputBuf, img)
			} else {
				// For invalid format, just write some bytes
				inputBuf.Write([]byte("invalid"))
			}
			
			outputBuf := &bytes.Buffer{}
			err := ProcessImageWithDPI(bytes.NewReader(inputBuf.Bytes()), outputBuf, tt.format, tt.dpi)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessImageWithDPI() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}