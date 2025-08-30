package transform

import (
	"image"
	"testing"
)

func TestCalculateDimensions(t *testing.T) {
	tests := []struct {
		name       string
		srcWidth   int
		srcHeight  int
		opts       ResizeOptions
		wantWidth  int
		wantHeight int
	}{
		{
			name:      "ResizeFit with both dimensions",
			srcWidth:  1000,
			srcHeight: 500,
			opts: ResizeOptions{
				Width:  400,
				Height: 400,
				Mode:   ResizeFit,
			},
			wantWidth:  400,
			wantHeight: 200,
		},
		{
			name:      "ResizeFit with width only",
			srcWidth:  1000,
			srcHeight: 500,
			opts: ResizeOptions{
				Width: 400,
				Mode:  ResizeFit,
			},
			wantWidth:  400,
			wantHeight: 200,
		},
		{
			name:      "ResizeFit with height only",
			srcWidth:  1000,
			srcHeight: 500,
			opts: ResizeOptions{
				Height: 250,
				Mode:   ResizeFit,
			},
			wantWidth:  500,
			wantHeight: 250,
		},
		{
			name:      "ResizeExact",
			srcWidth:  1000,
			srcHeight: 500,
			opts: ResizeOptions{
				Width:  300,
				Height: 300,
				Mode:   ResizeExact,
			},
			wantWidth:  300,
			wantHeight: 300,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWidth, gotHeight := CalculateDimensions(tt.srcWidth, tt.srcHeight, tt.opts)
			if gotWidth != tt.wantWidth || gotHeight != tt.wantHeight {
				t.Errorf("CalculateDimensions() = (%d, %d), want (%d, %d)",
					gotWidth, gotHeight, tt.wantWidth, tt.wantHeight)
			}
		})
	}
}

func TestResizeImage(t *testing.T) {
	// Create a test image
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	
	tests := []struct {
		name       string
		width      int
		height     int
		mode       ResizeMode
		wantWidth  int
		wantHeight int
	}{
		{
			name:       "ResizeFit",
			width:      50,
			height:     50,
			mode:       ResizeFit,
			wantWidth:  50,
			wantHeight: 50,
		},
		{
			name:       "ResizeExact",
			width:      50,
			height:     25,
			mode:       ResizeExact,
			wantWidth:  50,
			wantHeight: 25,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := resizeImage(img, tt.width, tt.height, tt.mode)
			if err != nil {
				t.Fatalf("resizeImage() error = %v", err)
			}
			
			bounds := result.Bounds()
			gotWidth := bounds.Max.X - bounds.Min.X
			gotHeight := bounds.Max.Y - bounds.Min.Y
			
			if gotWidth != tt.wantWidth || gotHeight != tt.wantHeight {
				t.Errorf("resizeImage() dimensions = (%d, %d), want (%d, %d)",
					gotWidth, gotHeight, tt.wantWidth, tt.wantHeight)
			}
		})
	}
}

func TestResizeByPercentage(t *testing.T) {
	// Create a test image
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	
	tests := []struct {
		name       string
		percentage float64
		wantWidth  int
		wantHeight int
	}{
		{
			name:       "50 percent",
			percentage: 50,
			wantWidth:  50,
			wantHeight: 50,
		},
		{
			name:       "200 percent",
			percentage: 200,
			wantWidth:  200,
			wantHeight: 200,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ResizeByPercentage(img, tt.percentage)
			bounds := result.Bounds()
			gotWidth := bounds.Max.X - bounds.Min.X
			gotHeight := bounds.Max.Y - bounds.Min.Y
			
			if gotWidth != tt.wantWidth || gotHeight != tt.wantHeight {
				t.Errorf("ResizeByPercentage() dimensions = (%d, %d), want (%d, %d)",
					gotWidth, gotHeight, tt.wantWidth, tt.wantHeight)
			}
		})
	}
}