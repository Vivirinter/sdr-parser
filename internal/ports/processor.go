package ports

import "github.com/Vivirinter/sdr-parser/internal/domain"

type SignalReader interface {
    Read() ([]float64, float64, error)
}

type SignalWriter interface {
    Write(samples []float64, sampleRate float64) error
}

// SignalProcessor defines core signal processing operations
type SignalProcessor interface {
    Process(samples []float64) ([]float64, error)
    Configure(config interface{}) error
}

type SignalAnalyzer interface {
    Analyze(samples []float64) (map[string]float64, error)
}

// SignalModulator defines the interface for signal modulation
type SignalModulator interface {
	Modulate(carrier, message *domain.Signal) (*domain.Signal, error)
}

// SignalDemodulator defines the interface for signal demodulation
type SignalDemodulator interface {
	Demodulate(signal *domain.Signal) (*domain.Signal, error)
}

// SignalFilter defines the interface for signal filtering
type SignalFilter interface {
	Filter(signal *domain.Signal) (*domain.Signal, error)
	Configure(params map[string]interface{}) error
}
