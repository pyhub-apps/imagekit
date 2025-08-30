package transform

import (
	"image"
	"io"
)

// ImageTransformer defines the interface for image transformation operations
type ImageTransformer interface {
	// Resize resizes an image to the specified dimensions
	Resize(input io.Reader, output io.Writer, options ResizeOptions) error
	
	// SetDPI sets the DPI metadata for an image
	SetDPI(input io.Reader, output io.Writer, dpi int) error
}

// ResizeOptions contains options for resizing an image
type ResizeOptions struct {
	Width  int     // Target width in pixels (0 = auto) - deprecated, use WidthDim
	Height int     // Target height in pixels (0 = auto) - deprecated, use HeightDim
	WidthDim  DimensionValue // Target width (can be pixels or multiplier)
	HeightDim DimensionValue // Target height (can be pixels or multiplier)
	Mode   ResizeMode // Resize mode
	Quality int    // JPEG quality (1-100)
}

// ResizeMode defines how the image should be resized
type ResizeMode int

const (
	// ResizeFit fits the image within the specified dimensions, maintaining aspect ratio
	ResizeFit ResizeMode = iota
	// ResizeFill fills the specified dimensions, cropping if necessary
	ResizeFill
	// ResizeExact resizes to exact dimensions, may distort aspect ratio
	ResizeExact
)

// Rectangle defines a rectangular area in an image
type Rectangle struct {
	X      int // X coordinate of top-left corner
	Y      int // Y coordinate of top-left corner
	Width  int // Width of the rectangle
	Height int // Height of the rectangle
}

// ImageFormat represents the format of an image
type ImageFormat string

const (
	FormatJPEG ImageFormat = "jpeg"
	FormatPNG  ImageFormat = "png"
)

// ImageInfo contains metadata about an image
type ImageInfo struct {
	Width  int
	Height int
	Format ImageFormat
	DPI    int
}

// Transformer implements the ImageTransformer interface
type Transformer struct {
	// preserveMetadata indicates whether to preserve EXIF data
	preserveMetadata bool
}

// NewTransformer creates a new image transformer
func NewTransformer() *Transformer {
	return &Transformer{
		preserveMetadata: false,
	}
}

// detectFormat detects the format of an image from its data
func detectFormat(img image.Image) (ImageFormat, error) {
	// This is a simplified detection - in real implementation,
	// we would check the actual image headers
	return FormatJPEG, nil
}