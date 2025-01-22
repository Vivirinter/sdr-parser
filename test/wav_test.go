package test

import (
	"os"
	"testing"

	"github.com/Vivirinter/sdr-parser/pkg/reader"
)

func TestWAVReader(t *testing.T) {
	tempFile := "test.wav"
	defer os.Remove(tempFile)

	file, err := os.Create(tempFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer file.Close()

	header := reader.WAVHeader{
		ChunkID:       [4]byte{'R', 'I', 'F', 'F'},
		ChunkSize:     36,
		Format:        [4]byte{'W', 'A', 'V', 'E'},
		Subchunk1ID:   [4]byte{'f', 'm', 't', ' '},
		Subchunk1Size: 16,
		AudioFormat:   1,
		NumChannels:   1,
		SampleRate:    44100,
		ByteRate:      88200,
		BlockAlign:    2,
		BitsPerSample: 16,
		Subchunk2ID:   [4]byte{'d', 'a', 't', 'a'},
		Subchunk2Size: 0,
	}

	if err := binary.Write(file, binary.LittleEndian, &header); err != nil {
		t.Fatalf("Failed to write header: %v", err)
	}

	reader, err := reader.NewWAVReader(tempFile)
	if err != nil {
		t.Fatalf("Failed to create WAV reader: %v", err)
	}
	defer reader.Close()

	if reader.Header.SampleRate != 44100 {
		t.Errorf("Expected sample rate 44100, got %d", reader.Header.SampleRate)
	}

	if reader.Header.NumChannels != 1 {
		t.Errorf("Expected 1 channel, got %d", reader.Header.NumChannels)
	}

	if reader.Header.BitsPerSample != 16 {
		t.Errorf("Expected 16 bits per sample, got %d", reader.Header.BitsPerSample)
	}
}
