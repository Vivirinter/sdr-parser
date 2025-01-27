// Package filter implements various digital signal processing filters
package filter

import (
	"math"
	"sort"
)

// Helper functions for statistics and math operations

// calculateMean calculates the arithmetic mean of a slice of float64 values
func calculateMean(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	var sum float64
	for _, value := range data {
		sum += value
	}
	return sum / float64(len(data))
}

// calculateStdDev calculates the standard deviation of a slice of float64 values
func calculateStdDev(data []float64) float64 {
	if len(data) < 2 {
		return 0
	}

	mean := calculateMean(data)
	var sumSquares float64
	for _, value := range data {
		diff := value - mean
		sumSquares += diff * diff
	}
	return math.Sqrt(sumSquares / float64(len(data)-1))
}

// calculateMedian calculates the median value of a slice
func calculateMedian(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}

	// Copy values to avoid modifying the input
	n := len(data)
	temp := make([]float64, n)
	copy(temp, data)
	sort.Float64s(temp)

	if n%2 == 0 {
		return (temp[n/2-1] + temp[n/2]) / 2
	}
	return temp[n/2]
}
