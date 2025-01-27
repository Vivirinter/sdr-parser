// Package filter implements various digital signal processing filters
package filter

import (
	"fmt"
	"math"
	"math/cmplx"
)

// ButterworthFilter implements a Butterworth low-pass filter.
// Butterworth filters are maximally flat in the passband and roll off
// towards zero in the stopband. The roll-off rate is determined by the filter order.
type ButterworthFilter struct {
	order      int       // Filter order
	cutoffFreq float64   // Cutoff frequency in Hz
	sampleRate float64   // Sample rate in Hz
	a          []float64 // Feedback coefficients
	b          []float64 // Feedforward coefficients
	x          []float64 // Input history buffer
	y          []float64 // Output history buffer
	stats      FilterStats  // Filter statistics
}

// NewButterworthFilter creates a new Butterworth filter instance
func NewButterworthFilter() *ButterworthFilter {
	return &ButterworthFilter{}
}

// Configure sets up the Butterworth filter with the provided configuration
func (f *ButterworthFilter) Configure(config FilterConfig) error {
	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Check Nyquist criterion
	nyquistFreq := config.SampleRate / 2
	if config.CutoffFreq >= nyquistFreq {
		return fmt.Errorf("cutoff frequency (%f Hz) must be less than Nyquist frequency (%f Hz)",
			config.CutoffFreq, nyquistFreq)
	}

	f.order = config.Order
	f.cutoffFreq = config.CutoffFreq
	f.sampleRate = config.SampleRate

	// Calculate filter coefficients
	if err := f.calculateCoefficients(); err != nil {
		return fmt.Errorf("failed to calculate coefficients: %w", err)
	}

	// Initialize history buffers
	f.x = make([]float64, len(f.b))
	f.y = make([]float64, len(f.a))

	return nil
}

// Process applies the Butterworth filter to the input samples
func (f *ButterworthFilter) Process(samples []float64) ([]float64, error) {
	if f.a == nil || f.b == nil {
		return nil, fmt.Errorf("filter not configured")
	}
	if len(samples) == 0 {
		return []float64{}, nil
	}

	result := make([]float64, len(samples))

	// Process each sample
	for i, sample := range samples {
		// Shift input history
		copy(f.x[1:], f.x)
		f.x[0] = sample

		// Calculate new output
		var sum float64
		for j := 0; j < len(f.b); j++ {
			sum += f.b[j] * f.x[j]
		}
		for j := 1; j < len(f.a); j++ {
			sum -= f.a[j] * f.y[j-1]
		}
		sum /= f.a[0]

		// Update output history and result
		copy(f.y[1:], f.y)
		f.y[0] = sum
		result[i] = sum
	}

	// Calculate statistics
	f.stats = FilterStats{
		InputMean:      calculateMean(samples),
		InputStdDev:    calculateStdDev(samples),
		OutputMean:     calculateMean(result),
		OutputStdDev:   calculateStdDev(result),
	}
	f.stats.NoiseReduction = 1 - (f.stats.OutputStdDev / f.stats.InputStdDev)

	return result, nil
}

// calculateCoefficients computes the filter coefficients using the bilinear transform
func (f *ButterworthFilter) calculateCoefficients() error {
	// Pre-warp the cutoff frequency
	omega := 2.0 * math.Pi * f.cutoffFreq
	warpedOmega := 2.0 * f.sampleRate * math.Tan(omega/(2.0*f.sampleRate))

	// Generate analog prototype poles
	poles := make([]complex128, f.order)
	for k := 0; k < f.order; k++ {
		theta := math.Pi * float64(2*k+1) / float64(2*f.order)
		poles[k] = complex(-math.Sin(theta), math.Cos(theta))
		poles[k] *= complex(warpedOmega, 0)
	}

	// Apply bilinear transform
	f.a = make([]float64, f.order+1)
	f.b = make([]float64, f.order+1)
	f.a[0] = 1.0
	f.b[0] = 1.0

	for _, pole := range poles {
		// Convert analog pole to digital using bilinear transform
		z := (1.0 + pole/(2.0*complex(f.sampleRate, 0))) / (1.0 - pole/(2.0*complex(f.sampleRate, 0)))

		// Check stability
		if cmplx.Abs(z) > MaxPoleRadius {
			return fmt.Errorf("unstable filter: pole magnitude %.3f exceeds maximum allowed %.3f",
				cmplx.Abs(z), MaxPoleRadius)
		}

		// Update coefficients using polynomial multiplication
		for i := f.order; i > 0; i-- {
			f.a[i] += real(z) * f.a[i-1]
			f.b[i] += real(z) * f.b[i-1]
		}
	}

	// Normalize coefficients for unity gain at DC
	gain := 0.0
	for i := 0; i <= f.order; i++ {
		gain += f.b[i]
	}
	if math.Abs(gain) < 1e-10 {
		return fmt.Errorf("filter gain too small: %e", gain)
	}
	for i := 0; i <= f.order; i++ {
		f.b[i] /= gain
	}

	return nil
}

// GetFrequencyResponse returns the filter's magnitude response at the specified frequencies
func (f *ButterworthFilter) GetFrequencyResponse(frequencies []float64) ([]float64, error) {
	if f.a == nil || f.b == nil {
		return nil, fmt.Errorf("filter not configured")
	}

	response := make([]float64, len(frequencies))
	for i, freq := range frequencies {
		omega := 2 * math.Pi * freq / f.sampleRate
		z := cmplx.Rect(1, omega)

		// Calculate H(z) = B(z)/A(z)
		num := complex(f.b[0], 0)
		den := complex(f.a[0], 0)
		zPow := complex(1, 0)
		
		for n := 1; n < len(f.b); n++ {
			zPow *= z
			num += complex(f.b[n], 0) * zPow
			den += complex(f.a[n], 0) * zPow
		}

		if cmplx.Abs(den) < 1e-10 {
			return nil, fmt.Errorf("division by zero in frequency response at %.1f Hz", freq)
		}

		response[i] = cmplx.Abs(num / den)
	}

	return response, nil
}

// GetStats returns the filter's statistics
func (f *ButterworthFilter) GetStats() FilterStats {
	return f.stats
}

// GetConfig returns the filter's configuration
func (f *ButterworthFilter) GetConfig() FilterConfig {
	return FilterConfig{
		Type:       Butterworth,
		Order:      f.order,
		CutoffFreq: f.cutoffFreq,
		SampleRate: f.sampleRate,
	}
}
