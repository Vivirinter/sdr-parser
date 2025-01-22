package processing

import (
	"math"
)

type Signal struct {
	Samples    []float64
	SampleRate float64
}

func NewSignal(sampleRate float64) *Signal {
	return &Signal{
		SampleRate: sampleRate,
		Samples:    make([]float64, 0),
	}
}

func (s *Signal) AddSample(sample float64) {
	s.Samples = append(s.Samples, sample)
}

func (s *Signal) GetSamples() []float64 {
	return s.Samples
}

func (s *Signal) GetSampleRate() float64 {
	return s.SampleRate
}

func (s *Signal) GetDuration() float64 {
	return float64(len(s.Samples)) / s.SampleRate
}

func (s *Signal) Normalize() {
	maxAmp := 0.0
	for _, sample := range s.Samples {
		if abs := math.Abs(sample); abs > maxAmp {
			maxAmp = abs
		}
	}

	if maxAmp > 0 {
		for i := range s.Samples {
			s.Samples[i] /= maxAmp
		}
	}
}

func (s *Signal) Filter(coeffs []float64) {
	filtered := make([]float64, len(s.Samples))
	for i := range s.Samples {
		sum := 0.0
		for j, coeff := range coeffs {
			if i-j >= 0 {
				sum += coeff * s.Samples[i-j]
			}
		}
		filtered[i] = sum
	}
	s.Samples = filtered
}

func (s *Signal) GenerateCarrier(freq float64, duration float64) {
	numSamples := int(duration * s.SampleRate)
	s.Samples = make([]float64, numSamples)
	
	for i := range s.Samples {
		t := float64(i) / s.SampleRate
		s.Samples[i] = math.Sin(2 * math.Pi * freq * t)
	}
}

func (s *Signal) ApplyAM(freq float64, depth float64) {
	carrier := make([]float64, len(s.Samples))
	for i := range carrier {
		t := float64(i) / s.SampleRate
		carrier[i] = math.Sin(2 * math.Pi * freq * t)
	}

	for i := range s.Samples {
		s.Samples[i] = carrier[i] * (1 + depth*s.Samples[i])
	}
}

func (s *Signal) ApplyFM(freq float64, deviation float64) {
	phase := 0.0
	modulated := make([]float64, len(s.Samples))
	
	for i := range s.Samples {
		phase += 2 * math.Pi * (freq + deviation*s.Samples[i]) / s.SampleRate
		modulated[i] = math.Sin(phase)
	}
	
	s.Samples = modulated
}

func (s *Signal) ApplySSB(freq float64, upperSideband bool) {
	hilbert := make([]float64, len(s.Samples))
	copy(hilbert, s.Samples)
	
	for i := range hilbert {
		var sum float64
		for j := range s.Samples {
			if i != j {
				sum += s.Samples[j] / (math.Pi * float64(i-j))
			}
		}
		hilbert[i] = sum
	}
	
	for i := range s.Samples {
		t := float64(i) / s.SampleRate
		carrier := math.Sin(2 * math.Pi * freq * t)
		carrierHilbert := -math.Cos(2 * math.Pi * freq * t)
		
		if upperSideband {
			s.Samples[i] = s.Samples[i]*carrier - hilbert[i]*carrierHilbert
		} else {
			s.Samples[i] = s.Samples[i]*carrier + hilbert[i]*carrierHilbert
		}
	}
}
