package transform

import (
	"image"
	"image/color"
	"testing"
)

// Helper function to create a test image with specific colors for verification
func createTestImageWithPattern(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	// Add a simple pattern to verify transformations
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if (x+y)%2 == 0 {
				img.Set(x, y, color.RGBA{255, 0, 0, 255}) // Red
			} else {
				img.Set(x, y, color.RGBA{0, 0, 255, 255}) // Blue
			}
		}
	}
	return img
}

func TestResizeByWidth(t *testing.T) {
	tests := []struct {
		name       string
		imgWidth   int
		imgHeight  int
		targetWidth int
		wantHeight int
	}{
		{
			name:       "Square image resize",
			imgWidth:   100,
			imgHeight:  100,
			targetWidth: 50,
			wantHeight: 50,
		},
		{
			name:       "Landscape image resize",
			imgWidth:   200,
			imgHeight:  100,
			targetWidth: 100,
			wantHeight: 50,
		},
		{
			name:       "Portrait image resize",
			imgWidth:   100,
			imgHeight:  200,
			targetWidth: 50,
			wantHeight: 100,
		},
		{
			name:       "Upscale image",
			imgWidth:   50,
			imgHeight:  100,
			targetWidth: 100,
			wantHeight: 200,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img := createTestImageWithPattern(tt.imgWidth, tt.imgHeight)
			result := ResizeByWidth(img, tt.targetWidth)
			
			bounds := result.Bounds()
			gotWidth := bounds.Max.X - bounds.Min.X
			gotHeight := bounds.Max.Y - bounds.Min.Y
			
			if gotWidth != tt.targetWidth {
				t.Errorf("ResizeByWidth() width = %d, want %d", gotWidth, tt.targetWidth)
			}
			if gotHeight != tt.wantHeight {
				t.Errorf("ResizeByWidth() height = %d, want %d", gotHeight, tt.wantHeight)
			}
		})
	}
}

func TestResizeByHeight(t *testing.T) {
	tests := []struct {
		name        string
		imgWidth    int
		imgHeight   int
		targetHeight int
		wantWidth   int
	}{
		{
			name:        "Square image resize",
			imgWidth:    100,
			imgHeight:   100,
			targetHeight: 50,
			wantWidth:   50,
		},
		{
			name:        "Landscape image resize",
			imgWidth:    200,
			imgHeight:   100,
			targetHeight: 50,
			wantWidth:   100,
		},
		{
			name:        "Portrait image resize",
			imgWidth:    100,
			imgHeight:   200,
			targetHeight: 100,
			wantWidth:   50,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img := createTestImageWithPattern(tt.imgWidth, tt.imgHeight)
			result := ResizeByHeight(img, tt.targetHeight)
			
			bounds := result.Bounds()
			gotWidth := bounds.Max.X - bounds.Min.X
			gotHeight := bounds.Max.Y - bounds.Min.Y
			
			if gotHeight != tt.targetHeight {
				t.Errorf("ResizeByHeight() height = %d, want %d", gotHeight, tt.targetHeight)
			}
			if gotWidth != tt.wantWidth {
				t.Errorf("ResizeByHeight() width = %d, want %d", gotWidth, tt.wantWidth)
			}
		})
	}
}

func TestThumbnail(t *testing.T) {
	tests := []struct {
		name         string
		imgWidth     int
		imgHeight    int
		thumbWidth   int
		thumbHeight  int
		maxExpWidth  int
		maxExpHeight int
	}{
		{
			name:         "Create thumbnail from large image",
			imgWidth:     1000,
			imgHeight:    1000,
			thumbWidth:   100,
			thumbHeight:  100,
			maxExpWidth:  100,
			maxExpHeight: 100,
		},
		{
			name:         "Create thumbnail maintaining aspect ratio",
			imgWidth:     200,
			imgHeight:    100,
			thumbWidth:   100,
			thumbHeight:  100,
			maxExpWidth:  100,
			maxExpHeight: 50,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img := createTestImageWithPattern(tt.imgWidth, tt.imgHeight)
			result := Thumbnail(img, tt.thumbWidth, tt.thumbHeight)
			
			bounds := result.Bounds()
			gotWidth := bounds.Max.X - bounds.Min.X
			gotHeight := bounds.Max.Y - bounds.Min.Y
			
			if gotWidth > tt.maxExpWidth {
				t.Errorf("Thumbnail() width = %d, want <= %d", gotWidth, tt.maxExpWidth)
			}
			if gotHeight > tt.maxExpHeight {
				t.Errorf("Thumbnail() height = %d, want <= %d", gotHeight, tt.maxExpHeight)
			}
		})
	}
}

func TestSmartCrop(t *testing.T) {
	tests := []struct {
		name       string
		imgWidth   int
		imgHeight  int
		cropWidth  int
		cropHeight int
	}{
		{
			name:       "Smart crop square to landscape",
			imgWidth:   100,
			imgHeight:  100,
			cropWidth:  150,
			cropHeight: 100,
		},
		{
			name:       "Smart crop landscape to portrait",
			imgWidth:   200,
			imgHeight:  100,
			cropWidth:  100,
			cropHeight: 150,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img := createTestImageWithPattern(tt.imgWidth, tt.imgHeight)
			result := SmartCrop(img, tt.cropWidth, tt.cropHeight)
			
			bounds := result.Bounds()
			gotWidth := bounds.Max.X - bounds.Min.X
			gotHeight := bounds.Max.Y - bounds.Min.Y
			
			if gotWidth != tt.cropWidth {
				t.Errorf("SmartCrop() width = %d, want %d", gotWidth, tt.cropWidth)
			}
			if gotHeight != tt.cropHeight {
				t.Errorf("SmartCrop() height = %d, want %d", gotHeight, tt.cropHeight)
			}
		})
	}
}

func TestResizeWithQuality(t *testing.T) {
	img := createTestImageWithPattern(100, 100)
	
	tests := []struct {
		name      string
		opts      ResizeWithQualityOptions
		wantWidth int
		wantHeight int
		wantErr   bool
	}{
		{
			name: "Resize with Lanczos filter",
			opts: ResizeWithQualityOptions{
				Width:   50,
				Height:  50,
				Mode:    ResizeFit,
				Filter:  FilterLanczos,
				Sharpen: false,
			},
			wantWidth: 50,
			wantHeight: 50,
			wantErr: false,
		},
		{
			name: "Resize with sharpening",
			opts: ResizeWithQualityOptions{
				Width:   75,
				Height:  75,
				Mode:    ResizeFit,
				Filter:  FilterLanczos,
				Sharpen: true,
			},
			wantWidth: 75,
			wantHeight: 75,
			wantErr: false,
		},
		{
			name: "Resize with Mitchell filter",
			opts: ResizeWithQualityOptions{
				Width:   60,
				Height:  60,
				Mode:    ResizeExact,
				Filter:  FilterMitchell,
				Sharpen: false,
			},
			wantWidth: 60,
			wantHeight: 60,
			wantErr: false,
		},
		{
			name: "Resize with Fill mode",
			opts: ResizeWithQualityOptions{
				Width:   80,
				Height:  40,
				Mode:    ResizeFill,
				Filter:  FilterCubic,
				Sharpen: false,
			},
			wantWidth: 80,
			wantHeight: 40,
			wantErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ResizeWithQuality(img, tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResizeWithQuality() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr {
				bounds := result.Bounds()
				gotWidth := bounds.Max.X - bounds.Min.X
				gotHeight := bounds.Max.Y - bounds.Min.Y
				
				if gotWidth != tt.wantWidth || gotHeight != tt.wantHeight {
					t.Errorf("ResizeWithQuality() dimensions = (%d, %d), want (%d, %d)",
						gotWidth, gotHeight, tt.wantWidth, tt.wantHeight)
				}
			}
		})
	}
}

func TestResizeImageEdgeCases(t *testing.T) {
	img := createTestImageWithPattern(100, 100)
	
	tests := []struct {
		name    string
		width   int
		height  int
		mode    ResizeMode
		wantErr bool
	}{
		{
			name:    "Zero dimensions",
			width:   0,
			height:  0,
			mode:    ResizeFit,
			wantErr: true,
		},
		{
			name:    "Negative width",
			width:   -50,
			height:  50,
			mode:    ResizeFit,
			wantErr: true,
		},
		{
			name:    "Only width specified with Fit",
			width:   50,
			height:  0,
			mode:    ResizeFit,
			wantErr: false,
		},
		{
			name:    "Only height specified with Fit",
			width:   0,
			height:  50,
			mode:    ResizeFit,
			wantErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := resizeImage(img, tt.width, tt.height, tt.mode)
			if (err != nil) != tt.wantErr {
				t.Errorf("resizeImage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestResizeByPercentageEdgeCases(t *testing.T) {
	img := createTestImageWithPattern(100, 100)
	
	tests := []struct {
		name       string
		percentage float64
		wantSame   bool // Should return same image
	}{
		{
			name:       "Zero percentage",
			percentage: 0,
			wantSame:   true,
		},
		{
			name:       "Negative percentage",
			percentage: -50,
			wantSame:   true,
		},
		{
			name:       "25 percent",
			percentage: 25,
			wantSame:   false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ResizeByPercentage(img, tt.percentage)
			
			if tt.wantSame {
				if result != img {
					t.Errorf("ResizeByPercentage() should return same image for percentage %f", tt.percentage)
				}
			} else {
				if result == img {
					t.Errorf("ResizeByPercentage() should return different image for percentage %f", tt.percentage)
				}
			}
		})
	}
}