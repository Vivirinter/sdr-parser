package domain

// Signal represents a digital signal in the time domain
type Signal struct {
	Samples    []float64
	SampleRate float64
}

// NewSignal creates a new Signal instance
func NewSignal(samples []float64, sampleRate float64) *Signal {
	return &Signal{
		Samples:    samples,
		SampleRate: sampleRate,
	}
}

// Duration returns the duration of the signal in seconds
func (s *Signal) Duration() float64 {
	return float64(len(s.Samples)) / s.SampleRate
}

// Clone creates a deep copy of the signal
func (s *Signal) Clone() *Signal {
	samples := make([]float64, len(s.Samples))
	copy(samples, s.Samples)
	return NewSignal(samples, s.SampleRate)
}
