package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/Vivirinter/sdr-parser/pkg/demod"
	"github.com/Vivirinter/sdr-parser/pkg/reader"
)

func getDemodCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "demod",
		Short: "Demodulate signal",
		RunE:  demodulateSignal,
	}

	cmd.Flags().StringP("input", "i", "", "input WAV file")
	cmd.Flags().StringP("output", "o", "audio.wav", "output WAV file")
	cmd.Flags().StringP("type", "t", "am", "demodulation type (am, fm, usb, lsb)")

	cmd.MarkFlagRequired("input")
	return cmd
}

func demodulateSignal(cmd *cobra.Command, args []string) error {
	input, _ := cmd.Flags().GetString("input")
	output, _ := cmd.Flags().GetString("output")
	demodType, _ := cmd.Flags().GetString("type")

	// Read input signal
	samples, sampleRate, err := reader.ReadWavFile(input)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// Apply demodulation
	var demodulated []float64
	switch demodType {
	case "am":
		demodulated = demod.AmDemodulate(samples)
	case "fm":
		demodulated = demod.FmDemodulate(samples)
	case "usb":
		demodulated = demod.UsbDemodulate(samples)
	case "lsb":
		demodulated = demod.LsbDemodulate(samples)
	default:
		return fmt.Errorf("unknown demodulation type: %s", demodType)
	}

	// Write to WAV file
	return reader.WriteWavFile(output, demodulated, sampleRate)
}
