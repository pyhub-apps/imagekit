package batch

import (
	"path/filepath"
	"strings"
)

// GenerateOutputPath generates the output file path by adding "_converted" suffix
// Example: "image.jpg" -> "image_converted.jpg"
// Example: "photos/pic.png" -> "photos/pic_converted.png"
func GenerateOutputPath(inputPath string) string {
	dir := filepath.Dir(inputPath)
	base := filepath.Base(inputPath)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	
	return filepath.Join(dir, name+"_converted"+ext)
}

// IsConvertedFile checks if a file already has the "_converted" suffix
func IsConvertedFile(path string) bool {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	
	return strings.HasSuffix(name, "_converted")
}