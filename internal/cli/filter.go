package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/Vivirinter/sdr-parser/internal/adapters/filters"
	"github.com/Vivirinter/sdr-parser/internal/domain"
	"github.com/Vivirinter/sdr-parser/pkg/reader"
	"path/filepath"
	"strings"
)

func getFilterCmd() *cobra.Command {
	var (
		inputFile    string
		outputFile   string
		filterType   string
		windowSize   int
		cutoffFreq   float64
		order        int
		sampleRate   float64
		amplitude    float64
		snr          float64
		normalize    bool
	)

	cmd := &cobra.Command{
		Use:   "filter",
		Short: "Apply filter to signal",
		Long: `Apply filter to signal. Available filter types:
  - moving_average: Simple moving average filter
  - median: Median filter for impulse noise reduction
  - butterworth: Butterworth low-pass filter`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Read input file
			samples, fileSampleRate, err := reader.ReadWavFile(inputFile)
			if err != nil {
				return fmt.Errorf("failed to read input file: %w", err)
			}

			// Use file sample rate if not specified
			if sampleRate == 0 {
				sampleRate = fileSampleRate
			}

			// Create input signal
			signal := &domain.Signal{
				Samples:    samples,
				SampleRate: sampleRate,
			}

			// Create filter through adapter
			f, err := filters.NewFilterAdapter(filterType)
			if err != nil {
				return fmt.Errorf("failed to create filter: %w", err)
			}

			// Configure filter
			params := map[string]interface{}{
				"window_size": windowSize,
				"cutoff_freq": cutoffFreq,
				"order":      order,
				"sampleRate": sampleRate,
				"amplitude":  amplitude,
				"snr":        snr,
				"normalize": normalize,
			}

			if err := f.Configure(params); err != nil {
				return fmt.Errorf("failed to configure filter: %w", err)
			}

			// Process signal
			filtered, err := f.Filter(signal)
			if err != nil {
				return fmt.Errorf("failed to process signal: %w", err)
			}

			// Write output file
			if outputFile == "" {
				ext := filepath.Ext(inputFile)
				base := strings.TrimSuffix(inputFile, ext)
				outputFile = fmt.Sprintf("%s_%s.wav", base, filterType)
			}
			if err := reader.WriteWavFile(outputFile, filtered.Samples, filtered.SampleRate); err != nil {
				return fmt.Errorf("failed to write output file: %w", err)
			}

			fmt.Printf("Successfully filtered signal from %s to %s using %s filter\n",
				inputFile, outputFile, filterType)
			return nil
		},
	}

	// Required flags
	cmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input WAV file")
	cmd.MarkFlagRequired("input")

	// Output file
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output WAV file (if not specified, will be input_[filter_type].wav)")

	// Filter type and parameters
	cmd.Flags().StringVarP(&filterType, "type", "t", "moving_average", "Filter type (moving_average, median, butterworth)")
	cmd.Flags().IntVarP(&windowSize, "window", "w", 5, "Window size for moving average/median filter")
	cmd.Flags().Float64VarP(&cutoffFreq, "cutoff", "c", 1000.0, "Cutoff frequency for Butterworth filter (Hz)")
	cmd.Flags().IntVarP(&order, "order", "n", 4, "Filter order for Butterworth filter")

	// Signal processing parameters
	cmd.Flags().Float64VarP(&sampleRate, "rate", "r", 0, "Sample rate (Hz). If not specified, uses input file's rate")
	cmd.Flags().Float64VarP(&amplitude, "amplitude", "a", 1.0, "Signal amplitude scaling factor")
	cmd.Flags().Float64VarP(&snr, "snr", "s", 0.0, "Signal-to-noise ratio for noise reduction (dB)")
	cmd.Flags().BoolVarP(&normalize, "normalize", "N", false, "Normalize signal after filtering")

	return cmd
}
