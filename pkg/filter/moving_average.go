// Package filter implements various digital signal processing filters
package filter

import (
	"fmt"
)

// MovingAverageFilter implements sliding window averaging for noise reduction
type MovingAverageFilter struct {
	BaseFilter
	buffer   []float64
	sum      float64
	position int
}

// NewMovingAverageFilter creates a new moving average filter
func NewMovingAverageFilter() *MovingAverageFilter {
	return &MovingAverageFilter{}
}

// Configure sets up the moving average filter with the given parameters
func (f *MovingAverageFilter) Configure(config FilterConfig) error {
	if err := f.BaseFilter.Configure(config); err != nil {
		return err
	}

	f.buffer = make([]float64, config.WindowSize)
	f.position = 0
	f.sum = 0

	return nil
}

// Process applies the moving average filter to the input signal
func (f *MovingAverageFilter) Process(input []float64) ([]float64, error) {
	if f.config.WindowSize == 0 {
		return nil, fmt.Errorf("filter not configured")
	}
	if len(input) == 0 {
		return []float64{}, nil
	}

	output := make([]float64, len(input))
	
	for i, sample := range input {
		f.sum -= f.buffer[f.position]
		f.buffer[f.position] = sample
		f.sum += sample
		
		output[i] = f.sum / float64(f.config.WindowSize)
		
		f.position = (f.position + 1) % f.config.WindowSize
	}

	f.updateStats(input, output)
	return output, nil
}

// GetStats returns the filter's statistics
func (f *MovingAverageFilter) GetStats() FilterStats {
	return f.BaseFilter.stats
}

// GetConfig returns the filter's configuration
func (f *MovingAverageFilter) GetConfig() FilterConfig {
	return f.BaseFilter.config
}
