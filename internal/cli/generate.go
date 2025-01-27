package cli

import (
	"fmt"
	"math"

	"github.com/spf13/cobra"
	"github.com/Vivirinter/sdr-parser/pkg/demod"
	"github.com/Vivirinter/sdr-parser/pkg/reader"
)

func getGenerateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate signal",
		RunE:  generateSignal,
	}

	cmd.Flags().StringP("output", "o", "signal.wav", "output WAV file")
	cmd.Flags().Float64P("freq", "f", 440.0, "frequency in Hz")
	cmd.Flags().Float64P("duration", "d", 5.0, "duration in seconds")
	cmd.Flags().StringP("mod", "m", "am", "modulation type (am, fm, usb, lsb)")

	return cmd
}

func generateSignal(cmd *cobra.Command, args []string) error {
	output, _ := cmd.Flags().GetString("output")
	freq, _ := cmd.Flags().GetFloat64("freq")
	duration, _ := cmd.Flags().GetFloat64("duration")
	modType, _ := cmd.Flags().GetString("mod")

	// Generate carrier signal
	sampleRate := 44100.0
	numSamples := int(duration * sampleRate)
	carrier := make([]float64, numSamples)
	for i := 0; i < numSamples; i++ {
		t := float64(i) / sampleRate
		carrier[i] = math.Sin(2 * math.Pi * freq * t)
	}

	// Generate modulating signal (message)
	message := make([]float64, numSamples)
	msgFreq := freq / 10 // Modulating frequency is 1/10 of carrier
	for i := 0; i < numSamples; i++ {
		t := float64(i) / sampleRate
		message[i] = math.Sin(2 * math.Pi * msgFreq * t)
	}

	// Apply modulation
	var modulated []float64
	switch modType {
	case "am":
		modulated = demod.AmModulate(carrier, message)
	case "fm":
		modulated = demod.FmModulate(carrier, message)
	case "usb":
		modulated = demod.UsbModulate(carrier, message)
	case "lsb":
		modulated = demod.LsbModulate(carrier, message)
	default:
		return fmt.Errorf("unknown modulation type: %s", modType)
	}

	// Write to WAV file
	return reader.WriteWavFile(output, modulated, sampleRate)
}
