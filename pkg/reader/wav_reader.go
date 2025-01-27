package reader

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// Standard WAV format constants
var (
	riffChunkID   = [4]byte{'R', 'I', 'F', 'F'}
	waveFormat    = [4]byte{'W', 'A', 'V', 'E'}
	fmtSubchunkID = [4]byte{'f', 'm', 't', ' '}
	dataChunkID   = [4]byte{'d', 'a', 't', 'a'}
)

// WAVHeader represents the header structure of a WAV file
type WAVHeader struct {
	ChunkID       [4]byte // Contains "RIFF"
	ChunkSize     uint32  // Size of the entire file minus 8 bytes
	Format        [4]byte // Contains "WAVE"
	Subchunk1ID   [4]byte // Contains "fmt "
	Subchunk1Size uint32  // Size of the fmt chunk (16 for PCM)
	AudioFormat   uint16  // Audio format (1 for PCM)
	NumChannels   uint16  // Number of channels
	SampleRate    uint32  // Sample rate (e.g., 44100)
	ByteRate      uint32  // SampleRate * NumChannels * BitsPerSample/8
	BlockAlign    uint16  // NumChannels * BitsPerSample/8
	BitsPerSample uint16  // Bits per sample (e.g., 16)
	Subchunk2ID   [4]byte // Contains "data"
	Subchunk2Size uint32  // Size of the data chunk
}

// WAVReader provides functionality to read WAV files
type WAVReader struct {
	file   *os.File
	Header WAVHeader
}

// NewWAVReader creates a new WAV reader for the specified file
func NewWAVReader(filename string) (*WAVReader, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open WAV file: %w", err)
	}

	reader := &WAVReader{file: file}
	if err := reader.readHeader(); err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to read WAV header: %w", err)
	}

	// Validate WAV format
	if err := reader.validateHeader(); err != nil {
		file.Close()
		return nil, fmt.Errorf("invalid WAV format: %w", err)
	}

	return reader, nil
}

// validateHeader checks if the WAV header is valid
func (r *WAVReader) validateHeader() error {
	if r.Header.ChunkID != riffChunkID {
		return fmt.Errorf("not a RIFF file")
	}
	if r.Header.Format != waveFormat {
		return fmt.Errorf("not a WAVE file")
	}
	if r.Header.Subchunk1ID != fmtSubchunkID {
		return fmt.Errorf("fmt chunk not found")
	}
	if r.Header.AudioFormat != 1 { // PCM = 1
		return fmt.Errorf("unsupported audio format: %d", r.Header.AudioFormat)
	}
	if r.Header.Subchunk2ID != dataChunkID {
		return fmt.Errorf("data chunk not found")
	}
	return nil
}

func (r *WAVReader) readHeader() error {
	return binary.Read(r.file, binary.LittleEndian, &r.Header)
}

// ReadSamples reads all audio samples from the WAV file
func (r *WAVReader) ReadSamples() ([]float64, error) {
	numSamples := r.Header.Subchunk2Size / uint32(r.Header.BlockAlign)
	samples := make([]float64, numSamples)

	for i := uint32(0); i < numSamples; i++ {
		var value int16
		if err := binary.Read(r.file, binary.LittleEndian, &value); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("failed to read sample: %w", err)
		}
		samples[i] = float64(value) / 32767.0
	}

	return samples, nil
}

// Close closes the WAV file
func (r *WAVReader) Close() error {
	return r.file.Close()
}

// ReadWavFile is a convenience function to read all samples from a WAV file
func ReadWavFile(filename string) ([]float64, float64, error) {
	reader, err := NewWAVReader(filename)
	if err != nil {
		return nil, 0, err
	}
	defer reader.Close()

	samples, err := reader.ReadSamples()
	if err != nil {
		return nil, 0, err
	}

	return samples, float64(reader.Header.SampleRate), nil
}

// WriteWavFile writes samples to a WAV file
func WriteWavFile(filename string, samples []float64, sampleRate float64) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create WAV file: %w", err)
	}
	defer file.Close()

	// Create WAV header
	header := WAVHeader{
		ChunkID:       riffChunkID,
		Format:        waveFormat,
		Subchunk1ID:   fmtSubchunkID,
		Subchunk1Size: 16,
		AudioFormat:   1, // PCM
		NumChannels:   1, // Mono
		SampleRate:    uint32(sampleRate),
		BitsPerSample: 16,
	}

	// Calculate dependent fields
	header.ByteRate = header.SampleRate * uint32(header.NumChannels) * uint32(header.BitsPerSample/8)
	header.BlockAlign = header.NumChannels * header.BitsPerSample/8
	header.Subchunk2ID = dataChunkID
	header.Subchunk2Size = uint32(len(samples) * int(header.BlockAlign))
	header.ChunkSize = 36 + header.Subchunk2Size

	// Write header
	if err := binary.Write(file, binary.LittleEndian, &header); err != nil {
		return fmt.Errorf("failed to write WAV header: %w", err)
	}

	// Write samples
	for _, sample := range samples {
		value := int16(sample * 32767.0)
		if err := binary.Write(file, binary.LittleEndian, value); err != nil {
			return fmt.Errorf("failed to write sample: %w", err)
		}
	}

	return nil
}
