package transform

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
)

const (
	// Standard DPI values
	DPI72  = 72
	DPI96  = 96
	DPI150 = 150
	DPI300 = 300
	
	// Conversion factor from DPI to pixels per meter
	inchesToMeters = 0.0254
)

// setImageDPI sets the DPI metadata for an image
func setImageDPI(img image.Image, format ImageFormat, dpi int) (image.Image, error) {
	// For now, we return the image as-is since Go's standard library
	// doesn't directly support DPI metadata manipulation
	// In a real implementation, we would need to handle the raw image data
	return img, nil
}

// SetJPEGDPI sets DPI for JPEG images by modifying JFIF header
func SetJPEGDPI(data []byte, dpi int) ([]byte, error) {
	// JPEG files start with SOI marker (0xFFD8)
	if len(data) < 2 || data[0] != 0xFF || data[1] != 0xD8 {
		return nil, fmt.Errorf("not a valid JPEG file")
	}
	
	// Look for JFIF APP0 marker (0xFFE0)
	index := 2
	for index < len(data)-10 {
		if data[index] == 0xFF && data[index+1] == 0xE0 {
			// Found APP0 marker
			length := int(data[index+2])<<8 | int(data[index+3])
			
			// Check for JFIF identifier
			if index+4+5 < len(data) &&
				data[index+4] == 'J' &&
				data[index+5] == 'F' &&
				data[index+6] == 'I' &&
				data[index+7] == 'F' &&
				data[index+8] == 0 {
				
				// JFIF header found
				// Density units: 0=no units, 1=dots/inch, 2=dots/cm
				data[index+11] = 1 // Set units to dots/inch
				
				// Set X and Y density (DPI)
				binary.BigEndian.PutUint16(data[index+12:], uint16(dpi))
				binary.BigEndian.PutUint16(data[index+14:], uint16(dpi))
				
				return data, nil
			}
			index += 2 + length
		} else if data[index] == 0xFF {
			index += 2
		} else {
			index++
		}
	}
	
	// If no JFIF header found, we need to insert one
	return insertJFIFHeader(data, dpi)
}

// insertJFIFHeader inserts a JFIF APP0 segment with DPI information
func insertJFIFHeader(data []byte, dpi int) ([]byte, error) {
	// Create JFIF APP0 segment
	jfif := []byte{
		0xFF, 0xE0, // APP0 marker
		0x00, 0x10, // Length (16 bytes)
		'J', 'F', 'I', 'F', 0x00, // JFIF identifier
		0x01, 0x01, // Version 1.1
		0x01,       // Units: dots/inch
		byte(dpi >> 8), byte(dpi), // X density
		byte(dpi >> 8), byte(dpi), // Y density
		0x00, 0x00, // No thumbnail
	}
	
	// Insert after SOI marker
	result := make([]byte, 0, len(data)+len(jfif))
	result = append(result, data[:2]...) // SOI marker
	result = append(result, jfif...)      // JFIF header
	result = append(result, data[2:]...)  // Rest of the file
	
	return result, nil
}

// SetPNGDPI sets DPI for PNG images by modifying pHYs chunk
func SetPNGDPI(data []byte, dpi int) ([]byte, error) {
	// PNG signature
	pngSignature := []byte{137, 80, 78, 71, 13, 10, 26, 10}
	
	if len(data) < 8 || !bytes.Equal(data[:8], pngSignature) {
		return nil, fmt.Errorf("not a valid PNG file")
	}
	
	// Calculate pixels per meter from DPI
	pixelsPerMeter := uint32(float64(dpi) / inchesToMeters)
	
	// Create pHYs chunk
	phys := createPHYsChunk(pixelsPerMeter, pixelsPerMeter)
	
	// Find IDAT chunk and insert pHYs before it
	index := 8
	result := make([]byte, 0, len(data)+len(phys))
	result = append(result, data[:8]...) // PNG signature
	
	inserted := false
	for index < len(data) {
		// Read chunk length
		if index+8 > len(data) {
			break
		}
		
		length := binary.BigEndian.Uint32(data[index : index+4])
		chunkType := string(data[index+4 : index+8])
		
		// Insert pHYs before IDAT
		if !inserted && chunkType == "IDAT" {
			result = append(result, phys...)
			inserted = true
		}
		
		// Skip existing pHYs chunks
		if chunkType != "pHYs" {
			chunkEnd := index + 12 + int(length) // 4 (length) + 4 (type) + length + 4 (CRC)
			if chunkEnd > len(data) {
				break
			}
			result = append(result, data[index:chunkEnd]...)
		}
		
		index += 12 + int(length)
	}
	
	// If no IDAT found, append pHYs at the end (before IEND)
	if !inserted && len(result) > 0 {
		// Find IEND and insert before it
		iendIndex := bytes.Index(result, []byte("IEND"))
		if iendIndex > 0 {
			newResult := make([]byte, 0, len(result)+len(phys))
			newResult = append(newResult, result[:iendIndex-4]...)
			newResult = append(newResult, phys...)
			newResult = append(newResult, result[iendIndex-4:]...)
			return newResult, nil
		}
	}
	
	return result, nil
}

// createPHYsChunk creates a PNG pHYs chunk with the given resolution
func createPHYsChunk(xPixelsPerMeter, yPixelsPerMeter uint32) []byte {
	data := make([]byte, 9)
	binary.BigEndian.PutUint32(data[0:4], xPixelsPerMeter)
	binary.BigEndian.PutUint32(data[4:8], yPixelsPerMeter)
	data[8] = 1 // Unit: meter
	
	// Create chunk
	chunk := make([]byte, 0, 21)
	chunk = append(chunk, 0, 0, 0, 9) // Length
	chunk = append(chunk, 'p', 'H', 'Y', 's') // Type
	chunk = append(chunk, data...) // Data
	
	// Calculate CRC
	crc := calculateCRC(append([]byte("pHYs"), data...))
	crcBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(crcBytes, crc)
	chunk = append(chunk, crcBytes...)
	
	return chunk
}

// calculateCRC calculates CRC32 for PNG chunk
func calculateCRC(data []byte) uint32 {
	// CRC32 table for PNG
	var crcTable [256]uint32
	for i := 0; i < 256; i++ {
		c := uint32(i)
		for j := 0; j < 8; j++ {
			if c&1 == 1 {
				c = 0xEDB88320 ^ (c >> 1)
			} else {
				c = c >> 1
			}
		}
		crcTable[i] = c
	}
	
	crc := uint32(0xFFFFFFFF)
	for _, b := range data {
		crc = crcTable[(crc^uint32(b))&0xFF] ^ (crc >> 8)
	}
	return ^crc
}

// GetImageDPI extracts DPI information from an image
func GetImageDPI(r io.Reader, format ImageFormat) (int, error) {
	// Read all data
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, r); err != nil {
		return 0, err
	}
	data := buf.Bytes()
	
	switch format {
	case FormatJPEG:
		return getJPEGDPI(data)
	case FormatPNG:
		return getPNGDPI(data)
	default:
		return 96, nil // Default DPI
	}
}

// getJPEGDPI extracts DPI from JPEG JFIF header
func getJPEGDPI(data []byte) (int, error) {
	if len(data) < 2 || data[0] != 0xFF || data[1] != 0xD8 {
		return 0, fmt.Errorf("not a valid JPEG file")
	}
	
	index := 2
	for index < len(data)-10 {
		if data[index] == 0xFF && data[index+1] == 0xE0 {
			// Found APP0 marker
			if index+14 < len(data) &&
				data[index+4] == 'J' &&
				data[index+5] == 'F' &&
				data[index+6] == 'I' &&
				data[index+7] == 'F' &&
				data[index+8] == 0 {
				
				// Check density units
				if data[index+11] == 1 { // dots/inch
					xDPI := binary.BigEndian.Uint16(data[index+12:])
					return int(xDPI), nil
				}
			}
		}
		index++
	}
	
	return 96, nil // Default DPI if not found
}

// getPNGDPI extracts DPI from PNG pHYs chunk
func getPNGDPI(data []byte) (int, error) {
	pngSignature := []byte{137, 80, 78, 71, 13, 10, 26, 10}
	
	if len(data) < 8 || !bytes.Equal(data[:8], pngSignature) {
		return 0, fmt.Errorf("not a valid PNG file")
	}
	
	index := 8
	for index < len(data) {
		if index+8 > len(data) {
			break
		}
		
		length := binary.BigEndian.Uint32(data[index : index+4])
		chunkType := string(data[index+4 : index+8])
		
		if chunkType == "pHYs" && index+12+9 <= len(data) {
			// Found pHYs chunk
			xPixelsPerMeter := binary.BigEndian.Uint32(data[index+8 : index+12])
			unit := data[index+16]
			
			if unit == 1 { // meter
				dpi := int(float64(xPixelsPerMeter) * inchesToMeters)
				return dpi, nil
			}
		}
		
		index += 12 + int(length)
	}
	
	return 96, nil // Default DPI if not found
}

// ConvertDPIValue converts between different DPI units
func ConvertDPIValue(value float64, fromUnit, toUnit string) (float64, error) {
	// Convert to dots per inch first
	var dpi float64
	switch fromUnit {
	case "dpi", "dots/inch":
		dpi = value
	case "dpcm", "dots/cm":
		dpi = value * 2.54
	case "pixels/meter":
		dpi = value * inchesToMeters
	default:
		return 0, fmt.Errorf("unsupported unit: %s", fromUnit)
	}
	
	// Convert to target unit
	switch toUnit {
	case "dpi", "dots/inch":
		return dpi, nil
	case "dpcm", "dots/cm":
		return dpi / 2.54, nil
	case "pixels/meter":
		return dpi / inchesToMeters, nil
	default:
		return 0, fmt.Errorf("unsupported unit: %s", toUnit)
	}
}

// ProcessImageWithDPI processes an image and sets its DPI
func ProcessImageWithDPI(r io.Reader, w io.Writer, format ImageFormat, dpi int) error {
	// Read image
	img, _, err := image.Decode(r)
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}
	
	// Encode to buffer
	buf := new(bytes.Buffer)
	switch format {
	case FormatJPEG:
		err = jpeg.Encode(buf, img, &jpeg.Options{Quality: 95})
	case FormatPNG:
		err = png.Encode(buf, img)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
	
	if err != nil {
		return fmt.Errorf("failed to encode image: %w", err)
	}
	
	// Modify DPI in encoded data
	data := buf.Bytes()
	switch format {
	case FormatJPEG:
		data, err = SetJPEGDPI(data, dpi)
	case FormatPNG:
		data, err = SetPNGDPI(data, dpi)
	}
	
	if err != nil {
		return fmt.Errorf("failed to set DPI: %w", err)
	}
	
	// Write to output
	_, err = w.Write(data)
	return err
}