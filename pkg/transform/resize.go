package transform

import (
	"fmt"
	"image"
	
	"github.com/disintegration/imaging"
)

// resizeImage resizes an image to the specified dimensions using the given mode
func resizeImage(img image.Image, width, height int, mode ResizeMode) (image.Image, error) {
	if width <= 0 && height <= 0 {
		return nil, fmt.Errorf("at least one dimension must be specified")
	}
	
	bounds := img.Bounds()
	srcWidth := bounds.Max.X - bounds.Min.X
	srcHeight := bounds.Max.Y - bounds.Min.Y
	
	switch mode {
	case ResizeFit:
		// Resize maintaining aspect ratio
		// If only one dimension is specified, use Resize with auto-calculation
		if width > 0 && height <= 0 {
			return imaging.Resize(img, width, 0, imaging.Lanczos), nil
		} else if height > 0 && width <= 0 {
			return imaging.Resize(img, 0, height, imaging.Lanczos), nil
		} else {
			// Both dimensions specified - resize to fit within bounds while maintaining aspect ratio
			// Calculate which dimension is the limiting factor
			ratio := float64(srcWidth) / float64(srcHeight)
			targetRatio := float64(width) / float64(height)
			
			if ratio > targetRatio {
				// Image is wider - fit to width
				return imaging.Resize(img, width, 0, imaging.Lanczos), nil
			} else {
				// Image is taller - fit to height
				return imaging.Resize(img, 0, height, imaging.Lanczos), nil
			}
		}
		
	case ResizeFill:
		// Fill the specified dimensions, cropping if necessary
		return imaging.Fill(img, width, height, imaging.Center, imaging.Lanczos), nil
		
	case ResizeExact:
		// Resize to exact dimensions, may distort aspect ratio
		if width <= 0 {
			width = srcWidth
		}
		if height <= 0 {
			height = srcHeight
		}
		return imaging.Resize(img, width, height, imaging.Lanczos), nil
		
	default:
		return nil, fmt.Errorf("unsupported resize mode: %v", mode)
	}
}

// ResizeByWidth resizes an image to the specified width, maintaining aspect ratio
func ResizeByWidth(img image.Image, width int) image.Image {
	return imaging.Resize(img, width, 0, imaging.Lanczos)
}

// ResizeByHeight resizes an image to the specified height, maintaining aspect ratio
func ResizeByHeight(img image.Image, height int) image.Image {
	return imaging.Resize(img, 0, height, imaging.Lanczos)
}

// ResizeByPercentage resizes an image by a percentage of its original size
func ResizeByPercentage(img image.Image, percentage float64) image.Image {
	if percentage <= 0 {
		return img
	}
	
	bounds := img.Bounds()
	width := int(float64(bounds.Max.X-bounds.Min.X) * percentage / 100)
	height := int(float64(bounds.Max.Y-bounds.Min.Y) * percentage / 100)
	
	return imaging.Resize(img, width, height, imaging.Lanczos)
}

// Thumbnail creates a thumbnail of the specified size
func Thumbnail(img image.Image, width, height int) image.Image {
	return imaging.Thumbnail(img, width, height, imaging.Lanczos)
}

// SmartCrop crops an image to the specified aspect ratio using smart cropping
func SmartCrop(img image.Image, width, height int) image.Image {
	// First, calculate the target aspect ratio
	targetRatio := float64(width) / float64(height)
	
	bounds := img.Bounds()
	srcWidth := float64(bounds.Max.X - bounds.Min.X)
	srcHeight := float64(bounds.Max.Y - bounds.Min.Y)
	srcRatio := srcWidth / srcHeight
	
	var cropWidth, cropHeight int
	
	if srcRatio > targetRatio {
		// Image is wider than target ratio, crop width
		cropHeight = int(srcHeight)
		cropWidth = int(srcHeight * targetRatio)
	} else {
		// Image is taller than target ratio, crop height
		cropWidth = int(srcWidth)
		cropHeight = int(srcWidth / targetRatio)
	}
	
	// Center crop
	cropped := imaging.CropCenter(img, cropWidth, cropHeight)
	
	// Resize to final dimensions
	return imaging.Resize(cropped, width, height, imaging.Lanczos)
}

// ResizeWithQuality resizes an image with specific quality settings
type ResizeWithQualityOptions struct {
	Width   int
	Height  int
	Mode    ResizeMode
	Filter  ResampleFilter
	Sharpen bool
}

// ResampleFilter represents the resampling filter to use
type ResampleFilter int

const (
	FilterNearest ResampleFilter = iota
	FilterLinear
	FilterCubic
	FilterLanczos
	FilterMitchell
)

// getImagingFilter converts our ResampleFilter to imaging.ResampleFilter
func getImagingFilter(filter ResampleFilter) imaging.ResampleFilter {
	switch filter {
	case FilterNearest:
		return imaging.NearestNeighbor
	case FilterLinear:
		return imaging.Linear
	case FilterCubic:
		return imaging.CatmullRom
	case FilterMitchell:
		return imaging.MitchellNetravali
	case FilterLanczos:
		fallthrough
	default:
		return imaging.Lanczos
	}
}

// ResizeWithQuality performs high-quality resizing with optional sharpening
func ResizeWithQuality(img image.Image, opts ResizeWithQualityOptions) (image.Image, error) {
	filter := getImagingFilter(opts.Filter)
	
	var result image.Image
	
	switch opts.Mode {
	case ResizeFit:
		result = imaging.Fit(img, opts.Width, opts.Height, filter)
	case ResizeFill:
		result = imaging.Fill(img, opts.Width, opts.Height, imaging.Center, filter)
	case ResizeExact:
		result = imaging.Resize(img, opts.Width, opts.Height, filter)
	default:
		return nil, fmt.Errorf("unsupported resize mode: %v", opts.Mode)
	}
	
	// Apply sharpening if requested
	if opts.Sharpen {
		result = imaging.Sharpen(result, 0.5)
	}
	
	return result, nil
}