package processing

import (
	"math"
)

type FilterType int

const (
	LowPass FilterType = iota
	HighPass
	BandPass
	BandStop
)

type WindowType int

const (
	Rectangular WindowType = iota
	Hamming
	Hanning
	Blackman
)

func DesignFIR(filterType FilterType, cutoffFreq, sampleRate float64, numTaps int, windowType WindowType) []float64 {
	coeffs := make([]float64, numTaps)
	
	fc := cutoffFreq / sampleRate
	center := float64(numTaps-1) / 2

	for i := range coeffs {
		n := float64(i) - center
		if n == 0 {
			coeffs[i] = 2 * fc
		} else {
			coeffs[i] = math.Sin(2*math.Pi*fc*n) / (math.Pi * n)
		}
		
		switch windowType {
		case Hamming:
			coeffs[i] *= 0.54 - 0.46*math.Cos(2*math.Pi*float64(i)/float64(numTaps-1))
		case Hanning:
			coeffs[i] *= 0.5 * (1 - math.Cos(2*math.Pi*float64(i)/float64(numTaps-1)))
		case Blackman:
			coeffs[i] *= 0.42 - 0.5*math.Cos(2*math.Pi*float64(i)/float64(numTaps-1)) +
				0.08*math.Cos(4*math.Pi*float64(i)/float64(numTaps-1))
		}
	}

	switch filterType {
	case HighPass:
		for i := range coeffs {
			coeffs[i] = -coeffs[i]
		}
		coeffs[int(center)] += 1
	case BandPass:
		fc2 := cutoffFreq * 1.5 / sampleRate
		for i := range coeffs {
			n := float64(i) - center
			if n == 0 {
				coeffs[i] = 2 * (fc2 - fc)
			} else {
				coeffs[i] = (math.Sin(2*math.Pi*fc2*n) - math.Sin(2*math.Pi*fc*n)) / (math.Pi * n)
			}
		}
	case BandStop:
		fc2 := cutoffFreq * 1.5 / sampleRate
		for i := range coeffs {
			n := float64(i) - center
			if n == 0 {
				coeffs[i] = 1 - 2*(fc2-fc)
			} else {
				coeffs[i] = (math.Sin(2*math.Pi*fc*n) - math.Sin(2*math.Pi*fc2*n)) / (math.Pi * n)
			}
		}
	}

	return coeffs
}

func ApplyFilter(samples []float64, coeffs []float64) []float64 {
	filtered := make([]float64, len(samples))
	
	for i := range samples {
		sum := 0.0
		for j, coeff := range coeffs {
			if i-j >= 0 {
				sum += coeff * samples[i-j]
			}
		}
		filtered[i] = sum
	}
	
	return filtered
}

func CreateDecimationFilter(factor int, sampleRate float64) []float64 {
	cutoffFreq := sampleRate / float64(2*factor)
	numTaps := factor*10 + 1
	return DesignFIR(LowPass, cutoffFreq, sampleRate, numTaps, Hamming)
}

func CreateInterpolationFilter(factor int, sampleRate float64) []float64 {
	cutoffFreq := sampleRate / float64(2*factor)
	numTaps := factor*10 + 1
	coeffs := DesignFIR(LowPass, cutoffFreq, sampleRate, numTaps, Hamming)
	
	for i := range coeffs {
		coeffs[i] *= float64(factor)
	}
	
	return coeffs
}

func Decimate(samples []float64, factor int) []float64 {
	decimated := make([]float64, len(samples)/factor)
	for i := range decimated {
		decimated[i] = samples[i*factor]
	}
	return decimated
}

func Interpolate(samples []float64, factor int) []float64 {
	interpolated := make([]float64, len(samples)*factor)
	for i := range samples {
		interpolated[i*factor] = samples[i]
	}
	return interpolated
}
