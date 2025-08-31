package transform

import (
	"image"
	"image/color"
	"testing"
)

func TestCropEdgesWithPixels(t *testing.T) {
	// Create a test image with distinct colors in each quadrant
	img := image.NewRGBA(image.Rect(0, 0, 200, 200))
	// Top-left: Red
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}
	// Top-right: Green
	for y := 0; y < 100; y++ {
		for x := 100; x < 200; x++ {
			img.Set(x, y, color.RGBA{0, 255, 0, 255})
		}
	}
	// Bottom-left: Blue
	for y := 100; y < 200; y++ {
		for x := 0; x < 100; x++ {
			img.Set(x, y, color.RGBA{0, 0, 255, 255})
		}
	}
	// Bottom-right: Yellow
	for y := 100; y < 200; y++ {
		for x := 100; x < 200; x++ {
			img.Set(x, y, color.RGBA{255, 255, 0, 255})
		}
	}
	
	tests := []struct {
		name       string
		opts       EdgeCropOptions
		wantWidth  int
		wantHeight int
		wantErr    bool
	}{
		{
			name: "Crop 50px from each edge",
			opts: EdgeCropOptions{
				Top:    CropValue{Value: 50, IsPercent: false},
				Bottom: CropValue{Value: 50, IsPercent: false},
				Left:   CropValue{Value: 50, IsPercent: false},
				Right:  CropValue{Value: 50, IsPercent: false},
			},
			wantWidth:  100,
			wantHeight: 100,
			wantErr:    false,
		},
		{
			name: "Crop 25% from each edge",
			opts: EdgeCropOptions{
				Top:    CropValue{Value: 25, IsPercent: true},
				Bottom: CropValue{Value: 25, IsPercent: true},
				Left:   CropValue{Value: 25, IsPercent: true},
				Right:  CropValue{Value: 25, IsPercent: true},
			},
			wantWidth:  100,
			wantHeight: 100,
			wantErr:    false,
		},
		{
			name: "No crop",
			opts: EdgeCropOptions{
				Top:    CropValue{Value: 0, IsPercent: false},
				Bottom: CropValue{Value: 0, IsPercent: false},
				Left:   CropValue{Value: 0, IsPercent: false},
				Right:  CropValue{Value: 0, IsPercent: false},
			},
			wantWidth:  200,
			wantHeight: 200,
			wantErr:    false,
		},
		{
			name: "Asymmetric crop",
			opts: EdgeCropOptions{
				Top:    CropValue{Value: 10, IsPercent: false},
				Bottom: CropValue{Value: 20, IsPercent: false},
				Left:   CropValue{Value: 30, IsPercent: false},
				Right:  CropValue{Value: 40, IsPercent: false},
			},
			wantWidth:  130,
			wantHeight: 170,
			wantErr:    false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CropEdges(img, tt.opts)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("CropRectangle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr {
				bounds := result.Bounds()
				gotWidth := bounds.Max.X - bounds.Min.X
				gotHeight := bounds.Max.Y - bounds.Min.Y
				
				if gotWidth != tt.wantWidth || gotHeight != tt.wantHeight {
					t.Errorf("CropRectangle() dimensions = (%d, %d), want (%d, %d)",
						gotWidth, gotHeight, tt.wantWidth, tt.wantHeight)
				}
			}
		})
	}
}

func TestCropEdgesWithMixedValues(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 200, 200))
	// Fill with a pattern
	for y := 0; y < 200; y++ {
		for x := 0; x < 200; x++ {
			img.Set(x, y, color.RGBA{uint8(x), uint8(y), 128, 255})
		}
	}
	
	tests := []struct {
		name       string
		opts       EdgeCropOptions
		wantWidth  int
		wantHeight int
		wantErr    bool
	}{
		{
			name: "Crop 10px from all edges",
			opts: EdgeCropOptions{
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
			name: "Crop only top and bottom",
			opts: EdgeCropOptions{
				Top:    CropValue{Value: 20, IsPercent: false},
				Right:  CropValue{Value: 0, IsPercent: false},
				Bottom: CropValue{Value: 30, IsPercent: false},
				Left:   CropValue{Value: 0, IsPercent: false},
			},
			wantWidth:  200,
			wantHeight: 150,
			wantErr:    false,
		},
		{
			name: "Crop only left and right",
			opts: EdgeCropOptions{
				Top:    CropValue{Value: 0, IsPercent: false},
				Right:  CropValue{Value: 25, IsPercent: false},
				Bottom: CropValue{Value: 0, IsPercent: false},
				Left:   CropValue{Value: 25, IsPercent: false},
			},
			wantWidth:  150,
			wantHeight: 200,
			wantErr:    false,
		},
		{
			name: "Asymmetric crop",
			opts: EdgeCropOptions{
				Top:    CropValue{Value: 5, IsPercent: false},
				Right:  CropValue{Value: 15, IsPercent: false},
				Bottom: CropValue{Value: 25, IsPercent: false},
				Left:   CropValue{Value: 35, IsPercent: false},
			},
			wantWidth:  150,
			wantHeight: 170,
			wantErr:    false,
		},
		{
			name: "No crop",
			opts: EdgeCropOptions{
				Top:    CropValue{Value: 0, IsPercent: false},
				Right:  CropValue{Value: 0, IsPercent: false},
				Bottom: CropValue{Value: 0, IsPercent: false},
				Left:   CropValue{Value: 0, IsPercent: false},
			},
			wantWidth:  200,
			wantHeight: 200,
			wantErr:    false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CropEdges(img, tt.opts)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("CropEdges() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr {
				bounds := result.Bounds()
				gotWidth := bounds.Max.X - bounds.Min.X
				gotHeight := bounds.Max.Y - bounds.Min.Y
				
				if gotWidth != tt.wantWidth || gotHeight != tt.wantHeight {
					t.Errorf("CropEdges() dimensions = (%d, %d), want (%d, %d)",
						gotWidth, gotHeight, tt.wantWidth, tt.wantHeight)
				}
			}
		})
	}
}

func TestValidateCropOptions(t *testing.T) {
	tests := []struct {
		name      string
		opts      EdgeCropOptions
		imgWidth  int
		imgHeight int
		wantErr   bool
	}{
		{
			name: "Valid crop options",
			opts: EdgeCropOptions{
				Top:    CropValue{Value: 10, IsPercent: false},
				Right:  CropValue{Value: 10, IsPercent: false},
				Bottom: CropValue{Value: 10, IsPercent: false},
				Left:   CropValue{Value: 10, IsPercent: false},
			},
			imgWidth:  100,
			imgHeight: 100,
			wantErr:   false,
		},
		{
			name: "Crop exceeds image width",
			opts: EdgeCropOptions{
				Top:    CropValue{Value: 0, IsPercent: false},
				Right:  CropValue{Value: 60, IsPercent: false},
				Bottom: CropValue{Value: 0, IsPercent: false},
				Left:   CropValue{Value: 60, IsPercent: false},
			},
			imgWidth:  100,
			imgHeight: 100,
			wantErr:   true,
		},
		{
			name: "Crop exceeds image height",
			opts: EdgeCropOptions{
				Top:    CropValue{Value: 60, IsPercent: false},
				Right:  CropValue{Value: 0, IsPercent: false},
				Bottom: CropValue{Value: 60, IsPercent: false},
				Left:   CropValue{Value: 0, IsPercent: false},
			},
			imgWidth:  100,
			imgHeight: 100,
			wantErr:   true,
		},
		{
			name: "Negative crop values",
			opts: EdgeCropOptions{
				Top:    CropValue{Value: -10, IsPercent: false},
				Right:  CropValue{Value: 10, IsPercent: false},
				Bottom: CropValue{Value: 10, IsPercent: false},
				Left:   CropValue{Value: 10, IsPercent: false},
			},
			imgWidth:  100,
			imgHeight: 100,
			wantErr:   false, // Negative values might be handled as 0
		},
		{
			name: "Exact crop to edges",
			opts: EdgeCropOptions{
				Top:    CropValue{Value: 50, IsPercent: false},
				Right:  CropValue{Value: 50, IsPercent: false},
				Bottom: CropValue{Value: 50, IsPercent: false},
				Left:   CropValue{Value: 50, IsPercent: false},
			},
			imgWidth:  100,
			imgHeight: 100,
			wantErr:   true, // This creates a 0x0 image which should error
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCropOptions(tt.opts, tt.imgWidth, tt.imgHeight)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCropOptions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseCropValue(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    CropValue
		wantErr bool
	}{
		{
			name:    "Parse pixel value",
			input:   "100",
			want:    CropValue{Value: 100, IsPercent: false},
			wantErr: false,
		},
		{
			name:    "Parse percentage value",
			input:   "25%",
			want:    CropValue{Value: 25, IsPercent: true},
			wantErr: false,
		},
		{
			name:    "Parse zero",
			input:   "0",
			want:    CropValue{Value: 0, IsPercent: false},
			wantErr: false,
		},
		{
			name:    "Parse empty string",
			input:   "",
			want:    CropValue{},
			wantErr: false,
		},
		{
			name:    "Invalid percentage",
			input:   "150%",
			want:    CropValue{},
			wantErr: true,
		},
		{
			name:    "Negative value",
			input:   "-10",
			want:    CropValue{},
			wantErr: true,
		},
		{
			name:    "Invalid format",
			input:   "abc",
			want:    CropValue{},
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCropValue(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCropValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseCropValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCropToAspectRatio(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 300, 200))
	
	tests := []struct {
		name         string
		widthRatio   int
		heightRatio  int
		wantWidth    int
		wantHeight   int
	}{
		{
			name:         "16:9 aspect ratio",
			widthRatio:   16,
			heightRatio:  9,
			wantWidth:    300,
			wantHeight:   169, // 300 * 9 / 16 ≈ 169
		},
		{
			name:         "1:1 square crop",
			widthRatio:   1,
			heightRatio:  1,
			wantWidth:    200,
			wantHeight:   200,
		},
		{
			name:         "4:3 aspect ratio",
			widthRatio:   4,
			heightRatio:  3,
			wantWidth:    267, // 200 * 4 / 3 ≈ 267
			wantHeight:   200,
		},
		{
			name:         "2:1 aspect ratio",
			widthRatio:   2,
			heightRatio:  1,
			wantWidth:    300,
			wantHeight:   150,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CropToAspectRatio(img, tt.widthRatio, tt.heightRatio)
			
			bounds := result.Bounds()
			gotWidth := bounds.Max.X - bounds.Min.X
			gotHeight := bounds.Max.Y - bounds.Min.Y
			
			// Allow small differences due to rounding
			widthDiff := gotWidth - tt.wantWidth
			heightDiff := gotHeight - tt.wantHeight
			
			if widthDiff < -2 || widthDiff > 2 || heightDiff < -2 || heightDiff > 2 {
				t.Errorf("CropAspectRatio() dimensions = (%d, %d), want approximately (%d, %d)",
					gotWidth, gotHeight, tt.wantWidth, tt.wantHeight)
			}
		})
	}
}