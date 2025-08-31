package transform

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"testing"
)

func TestTransformerResize(t *testing.T) {
	transformer := NewTransformer()
	
	// Create test image
	img := image.NewRGBA(image.Rect(0, 0, 200, 100))
	
	tests := []struct {
		name      string
		format    ImageFormat
		options   ResizeOptions
		wantWidth int
		wantHeight int
		wantErr   bool
	}{
		{
			name:   "Resize JPEG with Fit mode",
			format: FormatJPEG,
			options: ResizeOptions{
				Width:   100,
				Height:  100,
				Mode:    ResizeFit,
				Quality: 90,
			},
			wantWidth: 100,
			wantHeight: 50,
			wantErr: false,
		},
		{
			name:   "Resize PNG with Exact mode",
			format: FormatPNG,
			options: ResizeOptions{
				Width:   150,
				Height:  150,
				Mode:    ResizeExact,
				Quality: 0, // PNG ignores quality
			},
			wantWidth: 150,
			wantHeight: 150,
			wantErr: false,
		},
		{
			name:   "Resize with Fill mode",
			format: FormatJPEG,
			options: ResizeOptions{
				Width:   100,
				Height:  100,
				Mode:    ResizeFill,
				Quality: 85,
			},
			wantWidth: 100,
			wantHeight: 100,
			wantErr: false,
		},
		{
			name:   "Resize with only width",
			format: FormatJPEG,
			options: ResizeOptions{
				Width:   50,
				Height:  0,
				Mode:    ResizeFit,
				Quality: 95,
			},
			wantWidth: 50,
			wantHeight: 25,
			wantErr: false,
		},
		{
			name:   "Resize with only height",
			format: FormatPNG,
			options: ResizeOptions{
				Width:   0,
				Height:  50,
				Mode:    ResizeFit,
				Quality: 0,
			},
			wantWidth: 100,
			wantHeight: 50,
			wantErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode input image
			inputBuf := &bytes.Buffer{}
			switch tt.format {
			case FormatJPEG:
				if err := jpeg.Encode(inputBuf, img, nil); err != nil {
					t.Fatalf("Failed to encode JPEG: %v", err)
				}
			case FormatPNG:
				if err := png.Encode(inputBuf, img); err != nil {
					t.Fatalf("Failed to encode PNG: %v", err)
				}
			}
			
			// Resize
			outputBuf := &bytes.Buffer{}
			err := transformer.Resize(bytes.NewReader(inputBuf.Bytes()), outputBuf, tt.options)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("Transformer.Resize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr {
				// Decode output and check dimensions
				outputImg, _, err := image.Decode(outputBuf)
				if err != nil {
					t.Fatalf("Failed to decode output: %v", err)
				}
				
				bounds := outputImg.Bounds()
				gotWidth := bounds.Max.X - bounds.Min.X
				gotHeight := bounds.Max.Y - bounds.Min.Y
				
				if gotWidth != tt.wantWidth || gotHeight != tt.wantHeight {
					t.Errorf("Transformer.Resize() dimensions = (%d, %d), want (%d, %d)",
						gotWidth, gotHeight, tt.wantWidth, tt.wantHeight)
				}
			}
		})
	}
}

func TestTransformerSetDPI(t *testing.T) {
	transformer := NewTransformer()
	
	// Create test image
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	
	tests := []struct {
		name    string
		format  ImageFormat
		dpi     int
		wantErr bool
	}{
		{
			name:    "Set DPI for JPEG",
			format:  FormatJPEG,
			dpi:     300,
			wantErr: false,
		},
		{
			name:    "Set DPI for PNG",
			format:  FormatPNG,
			dpi:     96,
			wantErr: false,
		},
		{
			name:    "Set low DPI",
			format:  FormatJPEG,
			dpi:     72,
			wantErr: false,
		},
		{
			name:    "Set high DPI",
			format:  FormatPNG,
			dpi:     600,
			wantErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode input image
			inputBuf := &bytes.Buffer{}
			switch tt.format {
			case FormatJPEG:
				if err := jpeg.Encode(inputBuf, img, nil); err != nil {
					t.Fatalf("Failed to encode JPEG: %v", err)
				}
			case FormatPNG:
				if err := png.Encode(inputBuf, img); err != nil {
					t.Fatalf("Failed to encode PNG: %v", err)
				}
			}
			
			// Set DPI
			outputBuf := &bytes.Buffer{}
			err := transformer.SetDPI(bytes.NewReader(inputBuf.Bytes()), outputBuf, tt.dpi)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("Transformer.SetDPI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr {
				// Verify output is not empty
				if outputBuf.Len() == 0 {
					t.Errorf("Transformer.SetDPI() produced empty output")
				}
				
				// Verify output is valid image
				_, _, err := image.Decode(outputBuf)
				if err != nil {
					t.Errorf("Transformer.SetDPI() produced invalid image: %v", err)
				}
			}
		})
	}
}

func TestTransformerCropEdges(t *testing.T) {
	transformer := NewTransformer()
	
	// Create test image
	img := image.NewRGBA(image.Rect(0, 0, 200, 200))
	
	tests := []struct {
		name       string
		format     ImageFormat
		options    EdgeCropOptions
		wantWidth  int
		wantHeight int
		wantErr    bool
	}{
		{
			name:   "Crop edges from JPEG",
			format: FormatJPEG,
			options: EdgeCropOptions{
				Top:    CropValue{Value: 20, IsPercent: false},
				Right:  CropValue{Value: 30, IsPercent: false},
				Bottom: CropValue{Value: 20, IsPercent: false},
				Left:   CropValue{Value: 30, IsPercent: false},
			},
			wantWidth:  140,
			wantHeight: 160,
			wantErr:    false,
		},
		{
			name:   "Crop edges from PNG",
			format: FormatPNG,
			options: EdgeCropOptions{
				Top:    CropValue{Value: 10, IsPercent: false},
				Right:  CropValue{Value: 10, IsPercent: false},
				Bottom: CropValue{Value: 10, IsPercent: false},
				Left:   CropValue{Value: 10, IsPercent: false},
			},
			wantWidth:  180,
			wantHeight: 180,
			wantErr:    false,
		},
		{
			name:   "Invalid crop - exceeds bounds",
			format: FormatJPEG,
			options: EdgeCropOptions{
				Top:    CropValue{Value: 100, IsPercent: false},
				Right:  CropValue{Value: 100, IsPercent: false},
				Bottom: CropValue{Value: 100, IsPercent: false},
				Left:   CropValue{Value: 100, IsPercent: false},
			},
			wantWidth:  0,
			wantHeight: 0,
			wantErr:    true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode input image
			inputBuf := &bytes.Buffer{}
			switch tt.format {
			case FormatJPEG:
				if err := jpeg.Encode(inputBuf, img, nil); err != nil {
					t.Fatalf("Failed to encode JPEG: %v", err)
				}
			case FormatPNG:
				if err := png.Encode(inputBuf, img); err != nil {
					t.Fatalf("Failed to encode PNG: %v", err)
				}
			}
			
			// Crop edges
			outputBuf := &bytes.Buffer{}
			err := transformer.CropEdges(bytes.NewReader(inputBuf.Bytes()), outputBuf, tt.options)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("Transformer.CropEdges() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr {
				// Decode output and check dimensions
				outputImg, _, err := image.Decode(outputBuf)
				if err != nil {
					t.Fatalf("Failed to decode output: %v", err)
				}
				
				bounds := outputImg.Bounds()
				gotWidth := bounds.Max.X - bounds.Min.X
				gotHeight := bounds.Max.Y - bounds.Min.Y
				
				if gotWidth != tt.wantWidth || gotHeight != tt.wantHeight {
					t.Errorf("Transformer.CropEdges() dimensions = (%d, %d), want (%d, %d)",
						gotWidth, gotHeight, tt.wantWidth, tt.wantHeight)
				}
			}
		})
	}
}

func TestNewTransformer(t *testing.T) {
	transformer := NewTransformer()
	if transformer == nil {
		t.Errorf("NewTransformer() returned nil")
	}
}

func TestTransformerWithInvalidInput(t *testing.T) {
	transformer := NewTransformer()
	
	tests := []struct {
		name    string
		input   []byte
		wantErr bool
	}{
		{
			name:    "Empty input",
			input:   []byte{},
			wantErr: true,
		},
		{
			name:    "Invalid image data",
			input:   []byte("not an image"),
			wantErr: true,
		},
		{
			name:    "Corrupted JPEG header",
			input:   []byte{0xFF, 0xD8, 0xFF, 0x00, 0x00}, // Invalid JPEG
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test Resize with invalid input
			outputBuf := &bytes.Buffer{}
			err := transformer.Resize(bytes.NewReader(tt.input), outputBuf, ResizeOptions{
				Width:  100,
				Height: 100,
				Mode:   ResizeFit,
			})
			
			if (err == nil) != !tt.wantErr {
				t.Errorf("Transformer.Resize() with invalid input: error = %v, wantErr %v", err, tt.wantErr)
			}
			
			// Test SetDPI with invalid input
			outputBuf.Reset()
			err = transformer.SetDPI(bytes.NewReader(tt.input), outputBuf, 96)
			
			if (err == nil) != !tt.wantErr {
				t.Errorf("Transformer.SetDPI() with invalid input: error = %v, wantErr %v", err, tt.wantErr)
			}
			
			// Test CropEdges with invalid input
			outputBuf.Reset()
			err = transformer.CropEdges(bytes.NewReader(tt.input), outputBuf, EdgeCropOptions{
				Top:    CropValue{Value: 10, IsPercent: false},
				Right:  CropValue{Value: 10, IsPercent: false},
				Bottom: CropValue{Value: 10, IsPercent: false},
				Left:   CropValue{Value: 10, IsPercent: false},
			})
			
			if (err == nil) != !tt.wantErr {
				t.Errorf("Transformer.CropEdges() with invalid input: error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func BenchmarkTransformerResize(b *testing.B) {
	transformer := NewTransformer()
	
	// Create a test image
	img := image.NewRGBA(image.Rect(0, 0, 1920, 1080))
	inputBuf := &bytes.Buffer{}
	_ = jpeg.Encode(inputBuf, img, nil)
	inputData := inputBuf.Bytes()
	
	options := ResizeOptions{
		Width:   800,
		Height:  600,
		Mode:    ResizeFit,
		Quality: 85,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		outputBuf := &bytes.Buffer{}
		_ = transformer.Resize(bytes.NewReader(inputData), outputBuf, options)
	}
}

func BenchmarkTransformerSetDPI(b *testing.B) {
	transformer := NewTransformer()
	
	// Create a test image
	img := image.NewRGBA(image.Rect(0, 0, 1000, 1000))
	inputBuf := &bytes.Buffer{}
	_ = jpeg.Encode(inputBuf, img, nil)
	inputData := inputBuf.Bytes()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		outputBuf := &bytes.Buffer{}
		_ = transformer.SetDPI(bytes.NewReader(inputData), outputBuf, 300)
	}
}