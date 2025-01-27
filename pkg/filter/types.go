package filter

import "fmt"

const (
	MovingAverage FilterType = "moving_average"
	Median       FilterType = "median"
	Butterworth  FilterType = "butterworth"

	MinOrder           = 1
	MinFrequency       = 0.0
	MaxPoleRadius      = 0.99
	MinMedianWindowSize = 3
)

type FilterType string

func (ft FilterType) String() string {
	return string(ft)
}

func (ft FilterType) Validate() error {
	switch ft {
	case MovingAverage, Median, Butterworth:
		return nil
	default:
		return fmt.Errorf("unsupported filter type: %s", ft)
	}
}

// Filter interface defines the common operations for all filters in the system
type Filter interface {
	Configure(config FilterConfig) error
	Process(samples []float64) ([]float64, error)
	GetStats() FilterStats
	GetConfig() FilterConfig
}

type FilterConfig struct {
	Type       FilterType `json:"type"`
	WindowSize int       `json:"window_size,omitempty"`
	CutoffFreq float64   `json:"cutoff_freq,omitempty"`
	Order      int       `json:"order,omitempty"`
	SampleRate float64   `json:"sample_rate,omitempty"`
	Amplitude  float64   `json:"amplitude,omitempty"`
	SNR        float64   `json:"snr,omitempty"`
	Normalize  bool      `json:"normalize,omitempty"`
}

type FilterStats struct {
	InputSamples   int     `json:"input_samples"`
	OutputSamples  int     `json:"output_samples"`
	InputMean      float64 `json:"input_mean"`
	OutputMean     float64 `json:"output_mean"`
	InputStdDev    float64 `json:"input_std_dev"`
	OutputStdDev   float64 `json:"output_std_dev"`
	NoiseReduction float64 `json:"noise_reduction"`
	InputMedian    float64 `json:"input_median,omitempty"`
	OutputMedian   float64 `json:"output_median,omitempty"`
	MedianShift    float64 `json:"median_shift,omitempty"`
}

func (c FilterConfig) Validate() error {
	if err := c.Type.Validate(); err != nil {
		return err
	}

	switch c.Type {
	case MovingAverage, Median:
		if c.WindowSize <= 0 {
			return fmt.Errorf("window size must be positive")
		}
	case Butterworth:
		if c.Order <= 0 {
			return fmt.Errorf("filter order must be positive")
		}
		if c.CutoffFreq <= 0 {
			return fmt.Errorf("cutoff frequency must be positive")
		}
		if c.SampleRate <= 0 {
			return fmt.Errorf("sample rate must be positive")
		}
	}
	return nil
}
