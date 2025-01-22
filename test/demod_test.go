package test

import (
	"testing"

	"github.com/Vivirinter/sdr-parser/pkg/demod"
)

func TestAMDemodulation(t *testing.T) {
	samples := []float64{0.5, 1.0, 0.5, 0.0, -0.5, -1.0, -0.5, 0.0}
	demodulator := &demod.AMDemod{}
	
	output, metrics := demodulator.Demodulate(samples)
	
	if len(output) != len(samples) {
		t.Errorf("Expected output length %d, got %d", len(samples), len(output))
	}
	
	if metrics.CurrentGain <= 0 {
		t.Error("Expected positive gain")
	}
}

func TestFMDemodulation(t *testing.T) {
	samples := []float64{0.0, 0.7071, 1.0, 0.7071, 0.0, -0.7071, -1.0, -0.7071}
	demodulator := &demod.FMDemod{}
	
	output, metrics := demodulator.Demodulate(samples)
	
	if len(output) != len(samples)-1 {
		t.Errorf("Expected output length %d, got %d", len(samples)-1, len(output))
	}
	
	if metrics.CurrentGain <= 0 {
		t.Error("Expected positive gain")
	}
}

func TestSSBDemodulation(t *testing.T) {
	samples := []float64{0.0, 0.5, 1.0, 0.5, 0.0, -0.5, -1.0, -0.5}
	
	tests := []struct {
		name       string
		demodulator demod.Demodulator
	}{
		{"USB", &demod.USBDemod{}},
		{"LSB", &demod.LSBDemod{}},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, metrics := tt.demodulator.Demodulate(samples)
			
			if len(output) != len(samples) {
				t.Errorf("Expected output length %d, got %d", len(samples), len(output))
			}
			
			if metrics.CurrentGain <= 0 {
				t.Error("Expected positive gain")
			}
		})
	}
}
