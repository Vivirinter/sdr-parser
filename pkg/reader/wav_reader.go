package reader

import (
	"encoding/binary"
	"io"
	"os"
)

type WAVHeader struct {
	ChunkID       [4]byte
	ChunkSize     uint32
	Format        [4]byte
	Subchunk1ID   [4]byte
	Subchunk1Size uint32
	AudioFormat   uint16
	NumChannels   uint16
	SampleRate    uint32
	ByteRate      uint32
	BlockAlign    uint16
	BitsPerSample uint16
	Subchunk2ID   [4]byte
	Subchunk2Size uint32
}

type WAVReader struct {
	file   *os.File
	Header WAVHeader
}

func NewWAVReader(filename string) (*WAVReader, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	reader := &WAVReader{file: file}
	if err := reader.readHeader(); err != nil {
		file.Close()
		return nil, err
	}

	return reader, nil
}

func (r *WAVReader) readHeader() error {
	return binary.Read(r.file, binary.LittleEndian, &r.Header)
}

func (r *WAVReader) ReadSamples() ([]float64, error) {
	dataSize := r.Header.Subchunk2Size
	numSamples := int(dataSize) / int(r.Header.BlockAlign)
	samples := make([]float64, numSamples)

	for i := 0; i < numSamples; i++ {
		var sample int16
		err := binary.Read(r.file, binary.LittleEndian, &sample)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		samples[i] = float64(sample) / 32767.0
	}

	return samples, nil
}

func (r *WAVReader) Close() error {
	return r.file.Close()
}
