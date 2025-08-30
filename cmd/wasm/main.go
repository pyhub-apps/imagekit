// +build js,wasm

package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"strings"
	"syscall/js"

	"github.com/allieus/pyhub-imagekit/pkg/transform"
	"github.com/disintegration/imaging"
)

// Version is set at build time or defaults to the current version
var Version = "1.2535.17"

func main() {
	// Register JavaScript functions
	js.Global().Set("resizeImage", js.FuncOf(resizeImage))
	js.Global().Set("cropImage", js.FuncOf(cropImage))
	js.Global().Set("convertDPI", js.FuncOf(convertDPI))
	js.Global().Set("processImage", js.FuncOf(processImage))
	js.Global().Set("getImageKitVersion", js.FuncOf(getVersion))

	// Keep the Go program running
	select {}
}

// getVersion returns the ImageKit version
func getVersion(this js.Value, args []js.Value) interface{} {
	return Version
}

// resizeImage resizes an image
func resizeImage(this js.Value, args []js.Value) interface{} {
	if len(args) != 3 {
		return createErrorResult("resizeImage requires 3 arguments: imageData, width, height")
	}

	imageData := args[0].String()
	width := args[1].Int()
	height := args[2].Int()

	result, err := processImageData(imageData, func(img image.Image, format string) (image.Image, error) {
		// Handle different resize scenarios
		if width > 0 && height > 0 {
			// Both dimensions specified
			return imaging.Resize(img, width, height, imaging.Lanczos), nil
		} else if width > 0 {
			// Width only
			return transform.ResizeByWidth(img, width), nil
		} else if height > 0 {
			// Height only
			return transform.ResizeByHeight(img, height), nil
		}
		return img, nil
	})

	if err != nil {
		return createErrorResult(err.Error())
	}

	return createSuccessResult(result)
}

// cropImage crops edges from an image
func cropImage(this js.Value, args []js.Value) interface{} {
	if len(args) != 5 {
		return createErrorResult("cropImage requires 5 arguments: imageData, top, right, bottom, left")
	}

	imageData := args[0].String()
	top := args[1].String()
	right := args[2].String()
	bottom := args[3].String()
	left := args[4].String()

	result, err := processImageData(imageData, func(img image.Image, format string) (image.Image, error) {
		topVal, _ := transform.ParseCropValue(top)
		rightVal, _ := transform.ParseCropValue(right)
		bottomVal, _ := transform.ParseCropValue(bottom)
		leftVal, _ := transform.ParseCropValue(left)

		options := transform.EdgeCropOptions{
			Top:    topVal,
			Right:  rightVal,
			Bottom: bottomVal,
			Left:   leftVal,
		}
		return transform.CropEdges(img, options)
	})

	if err != nil {
		return createErrorResult(err.Error())
	}

	return createSuccessResult(result)
}

// convertDPI changes the DPI of an image
func convertDPI(this js.Value, args []js.Value) interface{} {
	if len(args) != 2 {
		return createErrorResult("convertDPI requires 2 arguments: imageData, dpi")
	}

	imageData := args[0].String()
	dpi := args[1].Int()

	// For WebAssembly, we'll just pass through the image since DPI is metadata
	// The actual DPI conversion will be handled by the download process
	result := map[string]interface{}{
		"success": true,
		"data":    imageData,
		"dpi":     dpi,
	}

	return result
}

// processImage is a combined function that can apply multiple transformations
func processImage(this js.Value, args []js.Value) interface{} {
	if len(args) != 2 {
		return createErrorResult("processImage requires 2 arguments: imageData, options")
	}

	imageData := args[0].String()
	options := args[1]

	result, err := processImageData(imageData, func(img image.Image, format string) (image.Image, error) {
		// Apply resize if specified
		if options.Get("resize").Bool() {
			widthVal := options.Get("width")
			heightVal := options.Get("height")
			
			width := 0
			height := 0
			
			if !widthVal.IsUndefined() && !widthVal.IsNull() {
				width = widthVal.Int()
			}
			if !heightVal.IsUndefined() && !heightVal.IsNull() {
				height = heightVal.Int()
			}
			
			if width > 0 && height > 0 {
				// Both dimensions specified
				img = imaging.Resize(img, width, height, imaging.Lanczos)
			} else if width > 0 {
				// Width only
				img = transform.ResizeByWidth(img, width)
			} else if height > 0 {
				// Height only
				img = transform.ResizeByHeight(img, height)
			}
		}

		// Apply crop if specified
		if options.Get("crop").Bool() {
			topVal := options.Get("cropTop")
			rightVal := options.Get("cropRight")
			bottomVal := options.Get("cropBottom")
			leftVal := options.Get("cropLeft")
			
			top := ""
			right := ""
			bottom := ""
			left := ""
			
			if !topVal.IsUndefined() && !topVal.IsNull() {
				top = topVal.String()
			}
			if !rightVal.IsUndefined() && !rightVal.IsNull() {
				right = rightVal.String()
			}
			if !bottomVal.IsUndefined() && !bottomVal.IsNull() {
				bottom = bottomVal.String()
			}
			if !leftVal.IsUndefined() && !leftVal.IsNull() {
				left = leftVal.String()
			}

			if top != "" || right != "" || bottom != "" || left != "" {
				topCrop, _ := transform.ParseCropValue(top)
				rightCrop, _ := transform.ParseCropValue(right)
				bottomCrop, _ := transform.ParseCropValue(bottom)
				leftCrop, _ := transform.ParseCropValue(left)

				cropOpts := transform.EdgeCropOptions{
					Top:    topCrop,
					Right:  rightCrop,
					Bottom: bottomCrop,
					Left:   leftCrop,
				}
				img, _ = transform.CropEdges(img, cropOpts)
			}
		}

		return img, nil
	})

	if err != nil {
		return createErrorResult(err.Error())
	}

	// Add DPI if specified
	resultMap := createSuccessResult(result).(map[string]interface{})
	dpiVal := options.Get("dpi")
	if !dpiVal.IsUndefined() && !dpiVal.IsNull() && !dpiVal.IsNaN() {
		// Check if dpi is a truthy value (not 0, false, null, undefined)
		if dpiVal.Type() == js.TypeNumber && dpiVal.Int() > 0 {
			resultMap["dpi"] = dpiVal.Int()
		} else if dpiVal.Type() == js.TypeBoolean && dpiVal.Bool() {
			// If dpi is true (boolean), use default 96
			resultMap["dpi"] = 96
		}
	}

	return resultMap
}

// processImageData handles the common image processing workflow
func processImageData(base64Data string, processor func(image.Image, string) (image.Image, error)) (string, error) {
	// Remove data URL prefix if present
	if strings.HasPrefix(base64Data, "data:") {
		parts := strings.SplitN(base64Data, ",", 2)
		if len(parts) == 2 {
			base64Data = parts[1]
		} else {
			return "", fmt.Errorf("invalid data URL format")
		}
	}

	// Clean up base64 string (remove any whitespace)
	base64Data = strings.TrimSpace(base64Data)
	
	// Decode base64
	imageBytes, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		// Try URL encoding if standard encoding fails
		imageBytes, err = base64.URLEncoding.DecodeString(base64Data)
		if err != nil {
			// Try RawStdEncoding as last resort
			imageBytes, err = base64.RawStdEncoding.DecodeString(base64Data)
			if err != nil {
				return "", fmt.Errorf("failed to decode base64: %w (data length: %d)", err, len(base64Data))
			}
		}
	}

	// Decode image
	reader := bytes.NewReader(imageBytes)
	img, format, err := image.Decode(reader)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	// Process image
	processedImg, err := processor(img, format)
	if err != nil {
		return "", err
	}

	// Encode result
	var buf bytes.Buffer
	switch format {
	case "jpeg":
		err = jpeg.Encode(&buf, processedImg, &jpeg.Options{Quality: 95})
	case "png":
		err = png.Encode(&buf, processedImg)
	default:
		// Default to PNG for unknown formats
		err = png.Encode(&buf, processedImg)
		format = "png"
	}

	if err != nil {
		return "", fmt.Errorf("failed to encode image: %w", err)
	}

	// Convert to base64
	result := base64.StdEncoding.EncodeToString(buf.Bytes())
	
	// Add data URL prefix
	mimeType := "image/" + format
	if format == "jpeg" {
		mimeType = "image/jpeg"
	}
	result = fmt.Sprintf("data:%s;base64,%s", mimeType, result)

	return result, nil
}

// createSuccessResult creates a success result object
func createSuccessResult(data string) interface{} {
	return map[string]interface{}{
		"success": true,
		"data":    data,
	}
}

// createErrorResult creates an error result object
func createErrorResult(message string) interface{} {
	return map[string]interface{}{
		"success": false,
		"error":   message,
	}
}