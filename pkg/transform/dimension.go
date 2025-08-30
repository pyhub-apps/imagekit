package transform

import (
	"fmt"
	"strconv"
	"strings"
)

// DimensionValue represents a dimension that can be either pixels or a multiplier
type DimensionValue struct {
	Value        int     // Pixel value when IsMultiplier is false
	IsMultiplier bool    // Whether this is a multiplier
	Multiplier   float64 // Multiplier value when IsMultiplier is true
}

// ParseDimension parses a dimension string like "1920", "2x", "x2", "0.5x"
func ParseDimension(s string) (DimensionValue, error) {
	if s == "" || s == "0" {
		return DimensionValue{Value: 0}, nil
	}
	
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	
	// Check for multiplier formats: "2x", "x2", "2.5x", "x2.5"
	if strings.Contains(s, "x") {
		// Remove 'x' and get the number part
		numStr := strings.ReplaceAll(s, "x", "")
		
		if numStr == "" {
			return DimensionValue{}, fmt.Errorf("invalid multiplier format: %s", s)
		}
		
		multiplier, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			return DimensionValue{}, fmt.Errorf("invalid multiplier value: %s", s)
		}
		
		// Validate multiplier range
		if multiplier <= 0 {
			return DimensionValue{}, fmt.Errorf("multiplier must be positive: %f", multiplier)
		}
		if multiplier > 10 {
			return DimensionValue{}, fmt.Errorf("multiplier too large (max 10x): %f", multiplier)
		}
		
		return DimensionValue{
			IsMultiplier: true,
			Multiplier:   multiplier,
		}, nil
	}
	
	// Parse as pixel value
	value, err := strconv.Atoi(s)
	if err != nil {
		return DimensionValue{}, fmt.Errorf("invalid pixel value: %s", s)
	}
	
	if value < 0 {
		return DimensionValue{}, fmt.Errorf("pixel value cannot be negative: %d", value)
	}
	
	return DimensionValue{
		Value:        value,
		IsMultiplier: false,
	}, nil
}

// Calculate returns the final pixel value based on the original size
func (d DimensionValue) Calculate(originalSize int) int {
	if d.IsMultiplier {
		result := float64(originalSize) * d.Multiplier
		// Round to nearest integer
		return int(result + 0.5)
	}
	return d.Value
}

// IsZero returns true if the dimension is not set (zero value)
func (d DimensionValue) IsZero() bool {
	if d.IsMultiplier {
		return d.Multiplier == 0
	}
	return d.Value == 0
}

// String returns the string representation of the dimension
func (d DimensionValue) String() string {
	if d.IsMultiplier {
		// Format multiplier with minimal decimal places
		if d.Multiplier == float64(int(d.Multiplier)) {
			return fmt.Sprintf("%dx", int(d.Multiplier))
		}
		return fmt.Sprintf("%.2fx", d.Multiplier)
	}
	return fmt.Sprintf("%d", d.Value)
}