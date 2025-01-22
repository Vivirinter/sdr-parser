package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Vivirinter/sdr-parser/pkg/demod"
	"github.com/Vivirinter/sdr-parser/pkg/reader"
)

var rootCmd = &cobra.Command{
	Use:   "sdrparser",
	Short: "SDR Parser",
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate signal",
	RunE:  generateSignal,
}

var demodCmd = &cobra.Command{
	Use:   "demod",
	Short: "Demodulate signal",
	RunE:  demodulateSignal,
}

func init() {
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(demodCmd)

	generateCmd.Flags().StringP("output", "o", "signal.wav", "output WAV file")
	generateCmd.Flags().Float64P("freq", "f", 440.0, "frequency in Hz")
	generateCmd.Flags().Float64P("duration", "d", 5.0, "duration in seconds")
	generateCmd.Flags().StringP("mod", "m", "am", "modulation type (am, fm, usb, lsb)")

	demodCmd.Flags().StringP("input", "i", "", "input WAV file")
	demodCmd.Flags().StringP("output", "o", "audio.wav", "output WAV file")
	demodCmd.Flags().StringP("type", "t", "am", "demodulation type (am, fm, usb, lsb)")

	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
}

var cfgFile string

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".sdrparser")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func generateSignal(cmd *cobra.Command, args []string) error {
	outputFile, _ := cmd.Flags().GetString("output")
	freq, _ := cmd.Flags().GetFloat64("freq")
	duration, _ := cmd.Flags().GetFloat64("duration")
	modType, _ := cmd.Flags().GetString("mod")

	sampleRate := 44100.0
	numSamples := int(duration * sampleRate)
	samples := make([]float64, numSamples)

	for i := range samples {
		t := float64(i) / sampleRate
		switch modType {
		case "am":
			samples[i] = math.Sin(2*math.Pi*freq*t) * (1 + math.Sin(2*math.Pi*2*t))
		case "fm":
			samples[i] = math.Sin(2*math.Pi*freq*t + 3*math.Sin(2*math.Pi*2*t))
		case "usb":
			samples[i] = math.Sin(2*math.Pi*freq*t) + math.Sin(2*math.Pi*(freq+100)*t)
		case "lsb":
			samples[i] = math.Sin(2*math.Pi*freq*t) + math.Sin(2*math.Pi*(freq-100)*t)
		}
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	header := reader.WAVHeader{
		ChunkID:       [4]byte{'R', 'I', 'F', 'F'},
		ChunkSize:     36 + uint32(len(samples)*2),
		Format:        [4]byte{'W', 'A', 'V', 'E'},
		Subchunk1ID:   [4]byte{'f', 'm', 't', ' '},
		Subchunk1Size: 16,
		AudioFormat:   1,
		NumChannels:   1,
		SampleRate:    uint32(sampleRate),
		ByteRate:      uint32(sampleRate * 2),
		BlockAlign:    2,
		BitsPerSample: 16,
		Subchunk2ID:   [4]byte{'d', 'a', 't', 'a'},
		Subchunk2Size: uint32(len(samples) * 2),
	}

	binary.Write(file, binary.LittleEndian, &header)

	for _, sample := range samples {
		binary.Write(file, binary.LittleEndian, int16(sample*32767))
	}

	fmt.Printf("Generated WAV file with %d samples at %d Hz\n", len(samples), int(sampleRate))
	return nil
}

func demodulateSignal(cmd *cobra.Command, args []string) error {
	inputFile, _ := cmd.Flags().GetString("input")
	outputFile, _ := cmd.Flags().GetString("output")
	demodType, _ := cmd.Flags().GetString("type")

	wavReader, err := reader.NewWAVReader(inputFile)
	if err != nil {
		return err
	}

	samples, err := wavReader.ReadSamples()
	if err != nil {
		return err
	}

	var demodulator demod.Demodulator
	switch demodType {
	case "am":
		demodulator = &demod.AMDemod{}
	case "fm":
		demodulator = &demod.FMDemod{}
	case "usb":
		demodulator = &demod.USBDemod{}
	case "lsb":
		demodulator = &demod.LSBDemod{}
	default:
		return fmt.Errorf("unknown demodulation type: %s", demodType)
	}

	demodulated, metrics := demodulator.Demodulate(samples)

	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	header := reader.WAVHeader{
		ChunkID:       [4]byte{'R', 'I', 'F', 'F'},
		ChunkSize:     36 + uint32(len(demodulated)*2),
		Format:        [4]byte{'W', 'A', 'V', 'E'},
		Subchunk1ID:   [4]byte{'f', 'm', 't', ' '},
		Subchunk1Size: 16,
		AudioFormat:   1,
		NumChannels:   1,
		SampleRate:    wavReader.Header.SampleRate,
		ByteRate:      wavReader.Header.SampleRate * 2,
		BlockAlign:    2,
		BitsPerSample: 16,
		Subchunk2ID:   [4]byte{'d', 'a', 't', 'a'},
		Subchunk2Size: uint32(len(demodulated) * 2),
	}

	binary.Write(file, binary.LittleEndian, &header)

	for _, sample := range demodulated {
		binary.Write(file, binary.LittleEndian, int16(sample*32767))
	}

	fmt.Printf("Demodulated %s signal, saved to %s\n", demodType, outputFile)
	fmt.Printf("Gain metrics: %+v\n", metrics)
	return nil
}

// SDR Parser - tool for generating and demodulating AM/FM/SSB signals
func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
