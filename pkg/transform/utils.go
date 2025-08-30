package transform

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
)

// LoadImage loads an image from a reader
func LoadImage(r io.Reader) (image.Image, ImageFormat, error) {
	// Read all data into buffer for format detection
	buf := &bytes.Buffer{}
	if _, err := io.Copy(buf, r); err != nil {
		return nil, "", fmt.Errorf("failed to read image data: %w", err)
	}
	
	// Try to decode the image
	img, format, err := image.Decode(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode image: %w", err)
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
	if opts.Mode == ResizeExact {
		return opts.Width, opts.Height
	}
	
	// If both dimensions are specified
	if opts.Width > 0 && opts.Height > 0 {
		if opts.Mode == ResizeFit {
			// Calculate dimensions to fit within the bounds
			ratio := float64(srcWidth) / float64(srcHeight)
			targetRatio := float64(opts.Width) / float64(opts.Height)
			
			if ratio > targetRatio {
				// Image is wider, fit to width
				return opts.Width, int(float64(opts.Width) / ratio)
			}
			// Image is taller, fit to height
			return int(float64(opts.Height) * ratio), opts.Height
		}
		// For ResizeFill, return the exact dimensions (cropping will be handled)
		return opts.Width, opts.Height
	}
	
	// If only width is specified
	if opts.Width > 0 {
		ratio := float64(srcHeight) / float64(srcWidth)
		return opts.Width, int(float64(opts.Width) * ratio)
	}
	
	// If only height is specified
	if opts.Height > 0 {
		ratio := float64(srcWidth) / float64(srcHeight)
		return int(float64(opts.Height) * ratio), opts.Height
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