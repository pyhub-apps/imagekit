package transform

import (
	"image"
	"image/color"
	"testing"
)

func TestAverageColors(t *testing.T) {
	colors := []color.Color{
		color.RGBA{R: 100, G: 100, B: 100, A: 255},
		color.RGBA{R: 200, G: 200, B: 200, A: 255},
	}
	
	avg := averageColors(colors)
	avgRGBA := avg.(color.RGBA)
	
	// Expected average: (150, 150, 150, 255)
	if avgRGBA.R != 150 || avgRGBA.G != 150 || avgRGBA.B != 150 {
		t.Errorf("averageColors() = %v, want approximately (150, 150, 150, 255)", avgRGBA)
	}
}

func TestGetBorderColors(t *testing.T) {
	// Create a test image with known border colors
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	
	// Set some specific colors on the border
	for x := 0; x < 10; x++ {
		img.Set(x, 0, color.RGBA{R: 255, G: 0, B: 0, A: 255}) // Top edge - red
		img.Set(x, 9, color.RGBA{R: 0, G: 255, B: 0, A: 255}) // Bottom edge - green
	}
	for y := 0; y < 10; y++ {
		img.Set(0, y, color.RGBA{R: 0, G: 0, B: 255, A: 255}) // Left edge - blue
		img.Set(9, y, color.RGBA{R: 255, G: 255, B: 0, A: 255}) // Right edge - yellow
	}
	
	area := Rectangle{X: 2, Y: 2, Width: 6, Height: 6}
	colors := getBorderColors(img, area)
	
	// Should have collected colors from around the area
	if len(colors) == 0 {
		t.Error("getBorderColors() returned no colors")
	}
}

func TestRemoveWatermarkMethods(t *testing.T) {
	// Create a test image
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	// Fill with a pattern
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			img.Set(x, y, color.RGBA{R: uint8(x), G: uint8(y), B: 128, A: 255})
		}
	}
	
	area := Rectangle{X: 10, Y: 10, Width: 20, Height: 20}
	
	tests := []struct {
		name   string
		method WatermarkRemovalMethod
	}{
		{
			name:   "Blur method",
			method: RemovalMethodBlur,
		},
		{
			name:   "Fill method",
			method: RemovalMethodFill,
		},
		{
			name:   "Inpaint method",
			method: RemovalMethodInpaint,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := RemoveWatermarkWithMethod(img, area, tt.method, nil)
			if err != nil {
				t.Fatalf("RemoveWatermarkWithMethod() error = %v", err)
			}
			if result == nil {
				t.Error("RemoveWatermarkWithMethod() returned nil image")
			}
			
			// Check that the image dimensions remain the same
			bounds := result.Bounds()
			if bounds.Max.X != 100 || bounds.Max.Y != 100 {
				t.Errorf("Result image dimensions changed: got %dx%d, want 100x100",
					bounds.Max.X, bounds.Max.Y)
			}
		})
	}
}

func TestRemoveMultipleWatermarks(t *testing.T) {
	// Create a test image
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	
	areas := []Rectangle{
		{X: 10, Y: 10, Width: 20, Height: 20},
		{X: 50, Y: 50, Width: 20, Height: 20},
	}
	
	result, err := RemoveMultipleWatermarks(img, areas, RemovalMethodBlur)
	if err != nil {
		t.Fatalf("RemoveMultipleWatermarks() error = %v", err)
	}
	if result == nil {
		t.Error("RemoveMultipleWatermarks() returned nil image")
	}
}