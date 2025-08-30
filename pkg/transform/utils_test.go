package transform

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"testing"
)

func TestLoadImage(t *testing.T) {
	// Create test JPEG image
	jpegImg := image.NewRGBA(image.Rect(0, 0, 10, 10))
	jpegBuf := &bytes.Buffer{}
	if err := jpeg.Encode(jpegBuf, jpegImg, nil); err != nil {
		t.Fatalf("Failed to create test JPEG: %v", err)
	}
	
	// Create test PNG image
	pngImg := image.NewRGBA(image.Rect(0, 0, 10, 10))
	pngBuf := &bytes.Buffer{}
	if err := png.Encode(pngBuf, pngImg); err != nil {
		t.Fatalf("Failed to create test PNG: %v", err)
	}
	
	tests := []struct {
		name       string
		input      *bytes.Buffer
		wantFormat ImageFormat
		wantErr    bool
	}{
		{
			name:       "Load JPEG",
			input:      jpegBuf,
			wantFormat: FormatJPEG,
			wantErr:    false,
		},
		{
			name:       "Load PNG",
			input:      pngBuf,
			wantFormat: FormatPNG,
			wantErr:    false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img, format, err := LoadImage(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if img == nil {
					t.Error("LoadImage() returned nil image")
				}
				if format != tt.wantFormat {
					t.Errorf("LoadImage() format = %v, want %v", format, tt.wantFormat)
				}
			}
		})
	}
}

func TestSaveImage(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	
	tests := []struct {
		name    string
		format  ImageFormat
		quality int
		wantErr bool
	}{
		{
			name:    "Save as JPEG",
			format:  FormatJPEG,
			quality: 95,
			wantErr: false,
		},
		{
			name:    "Save as PNG",
			format:  FormatPNG,
			quality: 0, // Quality is ignored for PNG
			wantErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			err := SaveImage(buf, img, tt.format, tt.quality)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveImage() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && buf.Len() == 0 {
				t.Error("SaveImage() produced empty output")
			}
		})
	}
}

func TestGetImageInfo(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 50))
	
	info := GetImageInfo(img, FormatJPEG)
	
	if info.Width != 100 {
		t.Errorf("GetImageInfo() Width = %d, want 100", info.Width)
	}
	if info.Height != 50 {
		t.Errorf("GetImageInfo() Height = %d, want 50", info.Height)
	}
	if info.Format != FormatJPEG {
		t.Errorf("GetImageInfo() Format = %v, want %v", info.Format, FormatJPEG)
	}
}

func TestValidateRectangle(t *testing.T) {
	tests := []struct {
		name      string
		rect      Rectangle
		imgWidth  int
		imgHeight int
		wantErr   bool
	}{
		{
			name: "Valid rectangle",
			rect: Rectangle{X: 10, Y: 10, Width: 50, Height: 50},
			imgWidth:  100,
			imgHeight: 100,
			wantErr:   false,
		},
		{
			name: "Negative coordinates",
			rect: Rectangle{X: -10, Y: 10, Width: 50, Height: 50},
			imgWidth:  100,
			imgHeight: 100,
			wantErr:   true,
		},
		{
			name: "Zero dimensions",
			rect: Rectangle{X: 10, Y: 10, Width: 0, Height: 50},
			imgWidth:  100,
			imgHeight: 100,
			wantErr:   true,
		},
		{
			name: "Exceeds bounds",
			rect: Rectangle{X: 60, Y: 10, Width: 50, Height: 50},
			imgWidth:  100,
			imgHeight: 100,
			wantErr:   true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRectangle(tt.rect, tt.imgWidth, tt.imgHeight)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRectangle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}