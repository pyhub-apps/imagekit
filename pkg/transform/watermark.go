package transform

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	
	"github.com/disintegration/imaging"
)

// WatermarkRemovalMethod defines the method to use for watermark removal
type WatermarkRemovalMethod int

const (
	// RemovalMethodBlur applies blur to the watermark area
	RemovalMethodBlur WatermarkRemovalMethod = iota
	// RemovalMethodFill fills the area with a solid color
	RemovalMethodFill
	// RemovalMethodInpaint attempts to inpaint the area (simple version)
	RemovalMethodInpaint
	// RemovalMethodClone clones from another area
	RemovalMethodClone
)

// removeWatermarkFromArea removes watermark from the specified area
func removeWatermarkFromArea(img image.Image, area Rectangle) (image.Image, error) {
	// Default to blur method for MVP
	return RemoveWatermarkWithMethod(img, area, RemovalMethodBlur, nil)
}

// WatermarkRemovalOptions contains options for watermark removal
type WatermarkRemovalOptions struct {
	Method      WatermarkRemovalMethod
	BlurRadius  float64     // For blur method
	FillColor   color.Color // For fill method
	CloneSource *Rectangle  // For clone method (source area to copy from)
}

// RemoveWatermarkWithMethod removes watermark using the specified method
func RemoveWatermarkWithMethod(img image.Image, area Rectangle, method WatermarkRemovalMethod, options *WatermarkRemovalOptions) (image.Image, error) {
	// Validate the area
	bounds := img.Bounds()
	imgWidth := bounds.Max.X - bounds.Min.X
	imgHeight := bounds.Max.Y - bounds.Min.Y
	
	if err := ValidateRectangle(area, imgWidth, imgHeight); err != nil {
		return nil, fmt.Errorf("invalid watermark area: %w", err)
	}
	
	// Create a copy of the image
	result := imaging.Clone(img)
	
	switch method {
	case RemovalMethodBlur:
		return applyBlurToArea(result, area, options)
	case RemovalMethodFill:
		return fillArea(result, area, options)
	case RemovalMethodInpaint:
		return inpaintArea(result, area)
	case RemovalMethodClone:
		return cloneArea(result, area, options)
	default:
		return nil, fmt.Errorf("unsupported removal method: %v", method)
	}
}

// applyBlurToArea applies blur to a specific area of the image
func applyBlurToArea(img image.Image, area Rectangle, options *WatermarkRemovalOptions) (image.Image, error) {
	blurRadius := 5.0
	if options != nil && options.BlurRadius > 0 {
		blurRadius = options.BlurRadius
	}
	
	// Extract the area to blur
	subImg := imaging.Crop(img, image.Rect(area.X, area.Y, area.X+area.Width, area.Y+area.Height))
	
	// Apply blur to the extracted area
	blurred := imaging.Blur(subImg, blurRadius)
	
	// Create a new image with the blurred area
	result := imaging.Clone(img)
	
	// Paste the blurred area back
	result = imaging.Paste(result, blurred, image.Pt(area.X, area.Y))
	
	return result, nil
}

// fillArea fills a specific area with a color
func fillArea(img image.Image, area Rectangle, options *WatermarkRemovalOptions) (image.Image, error) {
	fillColor := color.RGBA{255, 255, 255, 255} // Default to white
	if options != nil && options.FillColor != nil {
		fillColor = color.RGBAModel.Convert(options.FillColor).(color.RGBA)
	}
	
	// Create a new image with the area filled
	result := image.NewRGBA(img.Bounds())
	draw.Draw(result, result.Bounds(), img, image.Point{}, draw.Src)
	
	// Fill the specified area
	fillRect := image.Rect(area.X, area.Y, area.X+area.Width, area.Y+area.Height)
	draw.Draw(result, fillRect, &image.Uniform{fillColor}, image.Point{}, draw.Src)
	
	return result, nil
}

// inpaintArea performs simple inpainting by averaging surrounding pixels
func inpaintArea(img image.Image, area Rectangle) (image.Image, error) {
	// Create a new image
	result := image.NewRGBA(img.Bounds())
	draw.Draw(result, result.Bounds(), img, image.Point{}, draw.Src)
	
	// Simple inpainting: average colors from the border of the area
	borderColors := getBorderColors(img, area)
	avgColor := averageColors(borderColors)
	
	// Fill the area with the average color
	fillRect := image.Rect(area.X, area.Y, area.X+area.Width, area.Y+area.Height)
	draw.Draw(result, fillRect, &image.Uniform{avgColor}, image.Point{}, draw.Src)
	
	// Apply slight blur to blend better
	result = imaging.Blur(result, 1.0)
	
	return result, nil
}

// cloneArea clones pixels from another area of the image
func cloneArea(img image.Image, targetArea Rectangle, options *WatermarkRemovalOptions) (image.Image, error) {
	if options == nil || options.CloneSource == nil {
		return nil, fmt.Errorf("clone source area must be specified")
	}
	
	sourceArea := *options.CloneSource
	
	// Validate source area
	bounds := img.Bounds()
	imgWidth := bounds.Max.X - bounds.Min.X
	imgHeight := bounds.Max.Y - bounds.Min.Y
	
	if err := ValidateRectangle(sourceArea, imgWidth, imgHeight); err != nil {
		return nil, fmt.Errorf("invalid clone source area: %w", err)
	}
	
	// Check that source and target areas have the same size
	if sourceArea.Width != targetArea.Width || sourceArea.Height != targetArea.Height {
		return nil, fmt.Errorf("source and target areas must have the same dimensions")
	}
	
	// Extract the source area
	sourceImg := imaging.Crop(img, image.Rect(
		sourceArea.X, sourceArea.Y,
		sourceArea.X+sourceArea.Width, sourceArea.Y+sourceArea.Height,
	))
	
	// Create result image and paste the source to target area
	result := imaging.Clone(img)
	result = imaging.Paste(result, sourceImg, image.Pt(targetArea.X, targetArea.Y))
	
	return result, nil
}

// getBorderColors extracts colors from the border of an area
func getBorderColors(img image.Image, area Rectangle) []color.Color {
	colors := []color.Color{}
	
	// Top border
	for x := area.X; x < area.X+area.Width; x++ {
		if area.Y > 0 {
			colors = append(colors, img.At(x, area.Y-1))
		}
	}
	
	// Bottom border
	for x := area.X; x < area.X+area.Width; x++ {
		if area.Y+area.Height < img.Bounds().Max.Y {
			colors = append(colors, img.At(x, area.Y+area.Height))
		}
	}
	
	// Left border
	for y := area.Y; y < area.Y+area.Height; y++ {
		if area.X > 0 {
			colors = append(colors, img.At(area.X-1, y))
		}
	}
	
	// Right border
	for y := area.Y; y < area.Y+area.Height; y++ {
		if area.X+area.Width < img.Bounds().Max.X {
			colors = append(colors, img.At(area.X+area.Width, y))
		}
	}
	
	return colors
}

// averageColors calculates the average color from a slice of colors
func averageColors(colors []color.Color) color.Color {
	if len(colors) == 0 {
		return color.RGBA{255, 255, 255, 255}
	}
	
	var rSum, gSum, bSum, aSum uint32
	
	for _, c := range colors {
		r, g, b, a := c.RGBA()
		rSum += r
		gSum += g
		bSum += b
		aSum += a
	}
	
	count := uint32(len(colors))
	return color.RGBA{
		uint8(rSum / count / 256),
		uint8(gSum / count / 256),
		uint8(bSum / count / 256),
		uint8(aSum / count / 256),
	}
}

// DetectWatermark attempts to detect watermark areas in an image (simple version)
func DetectWatermark(img image.Image) ([]Rectangle, error) {
	// This is a placeholder for watermark detection
	// In a real implementation, this would use more sophisticated techniques
	// such as edge detection, pattern matching, or machine learning
	
	// For now, return empty slice (no watermarks detected)
	return []Rectangle{}, nil
}

// RemoveMultipleWatermarks removes multiple watermark areas from an image
func RemoveMultipleWatermarks(img image.Image, areas []Rectangle, method WatermarkRemovalMethod) (image.Image, error) {
	result := img
	
	for _, area := range areas {
		var err error
		result, err = RemoveWatermarkWithMethod(result, area, method, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to remove watermark at area %v: %w", area, err)
		}
	}
	
	return result, nil
}

// SmartWatermarkRemoval uses edge detection to blend watermark removal better
func SmartWatermarkRemoval(img image.Image, area Rectangle) (image.Image, error) {
	// Extract area around the watermark (slightly larger)
	padding := 10
	expandedArea := Rectangle{
		X:      maxInt(0, area.X-padding),
		Y:      maxInt(0, area.Y-padding),
		Width:  area.Width + 2*padding,
		Height: area.Height + 2*padding,
	}
	
	bounds := img.Bounds()
	imgWidth := bounds.Max.X - bounds.Min.X
	imgHeight := bounds.Max.Y - bounds.Min.Y
	
	// Adjust expanded area to fit within image bounds
	if expandedArea.X+expandedArea.Width > imgWidth {
		expandedArea.Width = imgWidth - expandedArea.X
	}
	if expandedArea.Y+expandedArea.Height > imgHeight {
		expandedArea.Height = imgHeight - expandedArea.Y
	}
	
	// Apply graduated blur for better blending
	result := imaging.Clone(img)
	
	// Create multiple blur levels
	for i := 0; i < 3; i++ {
		blurRadius := float64(5 - i)
		currentArea := Rectangle{
			X:      area.X + i*2,
			Y:      area.Y + i*2,
			Width:  area.Width - i*4,
			Height: area.Height - i*4,
		}
		
		if currentArea.Width > 0 && currentArea.Height > 0 {
			result, _ = applyBlurToArea(result, currentArea, &WatermarkRemovalOptions{BlurRadius: blurRadius})
		}
	}
	
	return result, nil
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}