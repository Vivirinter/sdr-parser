package demod

type DemodulationType int

const (
	AM DemodulationType = iota
	FM
	USB
	LSB
)

type GainMode int

const (
	Manual GainMode = iota
	AGC
)

type AGCSettings struct {
	AttackTime  float64
	ReleaseTime float64
	Target      float64
	MaxGain     float64
	MinGain     float64
}

type DemodulatorConfig struct {
	Type       DemodulationType
	GainMode   GainMode
	ManualGain float64
	AGCConfig  AGCSettings
}

type GainMetrics struct {
	CurrentGain   float64
	CompressionDB float64
	GainReduction float64
}

type Demodulator interface {
	Demodulate(samples []float64) ([]float64, GainMetrics)
}

type AMDemod struct{}
type FMDemod struct{}
type USBDemod struct{}
type LSBDemod struct{}

func (d *AMDemod) Demodulate(samples []float64) ([]float64, GainMetrics) {
	output := make([]float64, len(samples))
	maxAmp := 0.0

	for i, sample := range samples {
		if sample < 0 {
			sample = -sample
		}
		if sample > maxAmp {
			maxAmp = sample
		}
		output[i] = sample
	}

	gain := 1.0
	if maxAmp > 0 {
		gain = 0.7 / maxAmp
	}

	for i := range output {
		output[i] *= gain
	}

	metrics := GainMetrics{
		CurrentGain:   gain,
		CompressionDB: 20 * log10(gain),
		GainReduction: 1 / gain,
	}

	return output, metrics
}

func (d *FMDemod) Demodulate(samples []float64) ([]float64, GainMetrics) {
	output := make([]float64, len(samples)-1)
	maxAmp := 0.0

	for i := 0; i < len(samples)-1; i++ {
		diff := samples[i+1] - samples[i]
		if diff < 0 {
			diff = -diff
		}
		if diff > maxAmp {
			maxAmp = diff
		}
		output[i] = diff
	}

	gain := 1.0
	if maxAmp > 0 {
		gain = 0.7 / maxAmp
	}

	for i := range output {
		output[i] *= gain
	}

	metrics := GainMetrics{
		CurrentGain:   gain,
		CompressionDB: 20 * log10(gain),
		GainReduction: 1 / gain,
	}

	return output, metrics
}

func (d *USBDemod) Demodulate(samples []float64) ([]float64, GainMetrics) {
	output := make([]float64, len(samples))
	maxAmp := 0.0

	for i := 0; i < len(samples); i++ {
		val := samples[i]
		if i > 0 {
			val += samples[i-1]
		}
		if val < 0 {
			val = -val
		}
		if val > maxAmp {
			maxAmp = val
		}
		output[i] = val
	}

	gain := 1.0
	if maxAmp > 0 {
		gain = 0.7 / maxAmp
	}

	for i := range output {
		output[i] *= gain
	}

	metrics := GainMetrics{
		CurrentGain:   gain,
		CompressionDB: 20 * log10(gain),
		GainReduction: 1 / gain,
	}

	return output, metrics
}

func (d *LSBDemod) Demodulate(samples []float64) ([]float64, GainMetrics) {
	output := make([]float64, len(samples))
	maxAmp := 0.0

	for i := 0; i < len(samples); i++ {
		val := samples[i]
		if i < len(samples)-1 {
			val -= samples[i+1]
		}
		if val < 0 {
			val = -val
		}
		if val > maxAmp {
			maxAmp = val
		}
		output[i] = val
	}

	gain := 1.0
	if maxAmp > 0 {
		gain = 0.7 / maxAmp
	}

	for i := range output {
		output[i] *= gain
	}

	metrics := GainMetrics{
		CurrentGain:   gain,
		CompressionDB: 20 * log10(gain),
		GainReduction: 1 / gain,
	}

	return output, metrics
}

func log10(x float64) float64 {
	if x <= 0 {
		return -100
	}
	y := 0.0
	for x < 1 {
		x *= 10
		y--
	}
	for x >= 10 {
		x /= 10
		y++
	}
	return y + (x-1)/(x*2.302585092994046)
}

func NewDemodulator(config DemodulatorConfig) Demodulator {
	switch config.Type {
	case AM:
		return &AMDemod{}
	case FM:
		return &FMDemod{}
	case USB:
		return &USBDemod{}
	case LSB:
		return &LSBDemod{}
	default:
		return &AMDemod{}
	}
}
