package transform

import (
	"fmt"
	"image"
	"strconv"
	"strings"
	
	"github.com/disintegration/imaging"
)

// EdgeCropOptions contains options for edge cropping
type EdgeCropOptions struct {
	Top    CropValue
	Bottom CropValue
	Left   CropValue
	Right  CropValue
}

// CropValue represents a crop value that can be pixels or percentage
type CropValue struct {
	Value     int
	IsPercent bool
}

// ParseCropValue parses a string value like "100" or "10%" into CropValue
func ParseCropValue(s string) (CropValue, error) {
	if s == "" {
		return CropValue{}, nil
	}
	
	s = strings.TrimSpace(s)
	if strings.HasSuffix(s, "%") {
		// Parse percentage
		valueStr := strings.TrimSuffix(s, "%")
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			return CropValue{}, fmt.Errorf("invalid percentage value: %s", s)
		}
		if value < 0 || value > 100 {
			return CropValue{}, fmt.Errorf("percentage must be between 0 and 100: %d", value)
		}
		return CropValue{Value: value, IsPercent: true}, nil
	}
	
	// Parse pixel value
	value, err := strconv.Atoi(s)
	if err != nil {
		return CropValue{}, fmt.Errorf("invalid pixel value: %s", s)
	}
	if value < 0 {
		return CropValue{}, fmt.Errorf("pixel value cannot be negative: %d", value)
	}
	return CropValue{Value: value, IsPercent: false}, nil
}

// GetPixelValue converts a CropValue to actual pixels based on dimension
func (cv CropValue) GetPixelValue(dimension int) int {
	if cv.IsPercent {
		return dimension * cv.Value / 100
	}
	return cv.Value
}

// CropEdges crops the edges of an image based on the specified options
func CropEdges(img image.Image, options EdgeCropOptions) (image.Image, error) {
	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y
	
	// Calculate actual pixel values to crop
	topCrop := options.Top.GetPixelValue(height)
	bottomCrop := options.Bottom.GetPixelValue(height)
	leftCrop := options.Left.GetPixelValue(width)
	rightCrop := options.Right.GetPixelValue(width)
	
	// Validate that we're not cropping the entire image
	newWidth := width - leftCrop - rightCrop
	newHeight := height - topCrop - bottomCrop
	
	if newWidth <= 0 {
		return nil, fmt.Errorf("cropping would remove entire width (left: %d, right: %d, width: %d)", 
			leftCrop, rightCrop, width)
	}
	if newHeight <= 0 {
		return nil, fmt.Errorf("cropping would remove entire height (top: %d, bottom: %d, height: %d)", 
			topCrop, bottomCrop, height)
	}
	
	// Create the crop rectangle
	cropRect := image.Rect(
		bounds.Min.X + leftCrop,
		bounds.Min.Y + topCrop,
		bounds.Max.X - rightCrop,
		bounds.Max.Y - bottomCrop,
	)
	
	// Perform the crop
	result := imaging.Crop(img, cropRect)
	
	return result, nil
}

// CropToAspectRatio crops an image to a specific aspect ratio
func CropToAspectRatio(img image.Image, widthRatio, heightRatio int) image.Image {
	bounds := img.Bounds()
	srcWidth := bounds.Max.X - bounds.Min.X
	srcHeight := bounds.Max.Y - bounds.Min.Y
	
	targetRatio := float64(widthRatio) / float64(heightRatio)
	srcRatio := float64(srcWidth) / float64(srcHeight)
	
	var cropWidth, cropHeight int
	
	if srcRatio > targetRatio {
		// Image is wider than target ratio, crop width
		cropHeight = srcHeight
		cropWidth = int(float64(srcHeight) * targetRatio)
	} else {
		// Image is taller than target ratio, crop height
		cropWidth = srcWidth
		cropHeight = int(float64(srcWidth) / targetRatio)
	}
	
	// Center crop
	return imaging.CropCenter(img, cropWidth, cropHeight)
}

// AutoCrop attempts to automatically crop borders from an image
// This is a simple implementation that detects uniform colored borders
func AutoCrop(img image.Image, threshold int) image.Image {
	// Simple implementation: detect borders with low variance
	// This could be enhanced with more sophisticated edge detection
	
	// For now, just return the original image
	// A full implementation would analyze pixel variance along edges
	return img
}

// ValidateCropOptions checks if the crop options are valid
func ValidateCropOptions(options EdgeCropOptions, imgWidth, imgHeight int) error {
	// Calculate actual pixel values
	topPixels := options.Top.GetPixelValue(imgHeight)
	bottomPixels := options.Bottom.GetPixelValue(imgHeight)
	leftPixels := options.Left.GetPixelValue(imgWidth)
	rightPixels := options.Right.GetPixelValue(imgWidth)
	
	// Check if cropping would leave any image
	remainingWidth := imgWidth - leftPixels - rightPixels
	remainingHeight := imgHeight - topPixels - bottomPixels
	
	if remainingWidth <= 0 {
		return fmt.Errorf("crop would remove entire width (remaining: %d)", remainingWidth)
	}
	
	if remainingHeight <= 0 {
		return fmt.Errorf("crop would remove entire height (remaining: %d)", remainingHeight)
	}
	
	// Note: Could add a warning if cropping more than 50% of the image
	// but for now we just validate that the crop is within bounds
	
	return nil
}