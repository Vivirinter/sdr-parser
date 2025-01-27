package filter

import (
	"fmt"
	"sort"
)

// MedianFilter implements a median filter for impulse noise reduction.
// The median filter is particularly effective at removing impulse noise
// (salt and pepper noise) while preserving edges in the signal.
type MedianFilter struct {
	windowSize int       // Size of the sliding window (must be odd)
	window     []float64 // Buffer for window values
	stats      FilterStats
}

// NewMedianFilter creates a new median filter instance
func NewMedianFilter() *MedianFilter {
	return &MedianFilter{}
}

// Configure sets up the median filter with the provided configuration
func (f *MedianFilter) Configure(config FilterConfig) error {
	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	if config.WindowSize < MinMedianWindowSize {
		return fmt.Errorf("window size must be at least %d, got %d",
			MinMedianWindowSize, config.WindowSize)
	}

	// Ensure odd window size for symmetry
	f.windowSize = config.WindowSize
	if f.windowSize%2 == 0 {
		f.windowSize++
	}

	// Allocate buffer
	f.window = make([]float64, f.windowSize)
	f.stats = FilterStats{}

	return nil
}

// Process applies the median filter to the input samples.
// For each sample, it takes a window of neighboring samples,
// sorts them, and selects the median value as the output.
func (f *MedianFilter) Process(samples []float64) ([]float64, error) {
	if f.windowSize == 0 {
		return nil, fmt.Errorf("filter not configured")
	}
	if len(samples) == 0 {
		return []float64{}, nil
	}

	result := make([]float64, len(samples))
	halfWindow := f.windowSize / 2

	// Process each sample
	for i := 0; i < len(samples); i++ {
		// Determine window boundaries
		windowStart := max(0, i-halfWindow)
		windowEnd := min(len(samples), i+halfWindow+1)
		windowSize := windowEnd - windowStart

		// Copy window values and sort
		window := f.window[:windowSize]
		copy(window, samples[windowStart:windowEnd])
		sort.Float64s(window)

		// Select median value
		medianIdx := windowSize / 2
		if windowSize%2 == 0 {
			result[i] = (window[medianIdx-1] + window[medianIdx]) / 2
		} else {
			result[i] = window[medianIdx]
		}
	}

	// Calculate statistics
	f.stats = FilterStats{
		InputMean:      calculateMean(samples),
		InputStdDev:    calculateStdDev(samples),
		OutputMean:     calculateMean(result),
		OutputStdDev:   calculateStdDev(result),
		InputMedian:    calculateMedian(samples),
		OutputMedian:   calculateMedian(result),
	}
	f.stats.NoiseReduction = 1 - (f.stats.OutputStdDev / f.stats.InputStdDev)
	f.stats.MedianShift = f.stats.OutputMedian - f.stats.InputMedian

	return result, nil
}

// GetStats returns the filter's statistics
func (f *MedianFilter) GetStats() FilterStats {
	return f.stats
}

// GetConfig returns the filter's configuration
func (f *MedianFilter) GetConfig() FilterConfig {
	return FilterConfig{
		Type:       Median,
		WindowSize: f.windowSize,
	}
}
