package batch

import "github.com/allieus/pyhub-imagekit/pkg/transform"

// BatchResult contains the results of batch processing
type BatchResult struct {
	TotalFiles   int
	SuccessCount int
	FailedFiles  []FailedFile
}

// FailedFile represents a file that failed to process
type FailedFile struct {
	Path  string
	Error error
}

// ProcessOptions contains options for batch processing
type ProcessOptions struct {
	ResizeOptions *transform.ResizeOptions
	DPI           int
}

// HasErrors returns true if there were any failures
func (r *BatchResult) HasErrors() bool {
	return len(r.FailedFiles) > 0
}

// GetFailureRate returns the failure rate as a percentage
func (r *BatchResult) GetFailureRate() float64 {
	if r.TotalFiles == 0 {
		return 0
	}
	return float64(len(r.FailedFiles)) / float64(r.TotalFiles) * 100
}