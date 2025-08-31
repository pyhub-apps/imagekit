package transform

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	
	"github.com/disintegration/imaging"
)

// LoadImage loads an image from a reader and automatically corrects EXIF orientation
func LoadImage(r io.Reader) (image.Image, ImageFormat, error) {
	// Read all data into buffer for format detection
	buf := &bytes.Buffer{}
	if _, err := io.Copy(buf, r); err != nil {
		return nil, "", fmt.Errorf("failed to read image data: %w", err)
	}
	
	// Create a new reader from buffer
	reader := bytes.NewReader(buf.Bytes())
	
	// Try to decode the image with EXIF orientation support
	// The imaging library's Open function handles EXIF orientation automatically,
	// but we need to decode from a reader, so we'll use DecodeWithOrientation
	img, err := imaging.Decode(reader, imaging.AutoOrientation(true))
	if err != nil {
		// If imaging.Decode fails, try standard image.Decode as fallback
		_, _ = reader.Seek(0, 0)
		standardImg, format, decodeErr := image.Decode(reader)
		if decodeErr != nil {
			return nil, "", fmt.Errorf("failed to decode image: %w", decodeErr)
		}
		
		// Convert format string to our ImageFormat type
		var imgFormat ImageFormat
		switch format {
		case "jpeg", "jpg":
			imgFormat = FormatJPEG
		case "png":
			imgFormat = FormatPNG
		default:
			return nil, "", fmt.Errorf("unsupported image format: %s", format)
		}
		
		return standardImg, imgFormat, nil
	}
	
	// Detect format from the buffer
	_, _ = reader.Seek(0, 0)
	_, format, err := image.DecodeConfig(reader)
	if err != nil {
		return nil, "", fmt.Errorf("failed to detect image format: %w", err)
	}
	
	// Convert format string to our ImageFormat type
	var imgFormat ImageFormat
	switch format {
	case "jpeg", "jpg":
		imgFormat = FormatJPEG
	case "png":
		imgFormat = FormatPNG
	default:
		return nil, "", fmt.Errorf("unsupported image format: %s", format)
	}
	
	return img, imgFormat, nil
}

// SaveImage saves an image to a writer
func SaveImage(w io.Writer, img image.Image, format ImageFormat, quality int) error {
	switch format {
	case FormatJPEG:
		opts := &jpeg.Options{
			Quality: quality,
		}
		if quality <= 0 {
			opts.Quality = 95 // Default high quality
		}
		return jpeg.Encode(w, img, opts)
	case FormatPNG:
		return png.Encode(w, img)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

// GetImageInfo extracts information about an image
func GetImageInfo(img image.Image, format ImageFormat) ImageInfo {
	bounds := img.Bounds()
	return ImageInfo{
		Width:  bounds.Max.X - bounds.Min.X,
		Height: bounds.Max.Y - bounds.Min.Y,
		Format: format,
		DPI:    96, // Default DPI
	}
}

// CalculateDimensions calculates target dimensions based on resize options
func CalculateDimensions(srcWidth, srcHeight int, opts ResizeOptions) (int, int) {
	// Calculate actual pixel values from DimensionValue
	var targetWidth, targetHeight int
	
	// Use new DimensionValue fields if available, otherwise fallback to old fields
	if !opts.WidthDim.IsZero() || !opts.HeightDim.IsZero() {
		targetWidth = opts.WidthDim.Calculate(srcWidth)
		targetHeight = opts.HeightDim.Calculate(srcHeight)
	} else {
		// Fallback to old Width/Height fields for backward compatibility
		targetWidth = opts.Width
		targetHeight = opts.Height
	}
	
	if opts.Mode == ResizeExact {
		return targetWidth, targetHeight
	}
	
	// If both dimensions are specified
	if targetWidth > 0 && targetHeight > 0 {
		if opts.Mode == ResizeFit {
			// Calculate dimensions to fit within the bounds
			ratio := float64(srcWidth) / float64(srcHeight)
			targetRatio := float64(targetWidth) / float64(targetHeight)
			
			if ratio > targetRatio {
				// Image is wider, fit to width
				return targetWidth, int(float64(targetWidth) / ratio)
			}
			// Image is taller, fit to height
			return int(float64(targetHeight) * ratio), targetHeight
		}
		// For ResizeFill, return the exact dimensions (cropping will be handled)
		return targetWidth, targetHeight
	}
	
	// If only width is specified
	if targetWidth > 0 {
		ratio := float64(srcHeight) / float64(srcWidth)
		return targetWidth, int(float64(targetWidth) * ratio)
	}
	
	// If only height is specified
	if targetHeight > 0 {
		ratio := float64(srcWidth) / float64(srcHeight)
		return int(float64(targetHeight) * ratio), targetHeight
	}
	
	// No dimensions specified, return original
	return srcWidth, srcHeight
}

// ValidateRectangle validates that a rectangle is within image bounds
func ValidateRectangle(rect Rectangle, imgWidth, imgHeight int) error {
	if rect.X < 0 || rect.Y < 0 {
		return fmt.Errorf("rectangle coordinates cannot be negative")
	}
	if rect.Width <= 0 || rect.Height <= 0 {
		return fmt.Errorf("rectangle dimensions must be positive")
	}
	if rect.X+rect.Width > imgWidth || rect.Y+rect.Height > imgHeight {
		return fmt.Errorf("rectangle exceeds image bounds")
	}
	return nil
}