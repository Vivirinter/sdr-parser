# SDR Parser

Software Defined Radio (SDR) signal processing tool with support for various modulation types and signal processing capabilities.

## 🚀 Features

- **Signal Generation:**
  - AM (Amplitude Modulation)
  - FM (Frequency Modulation)
  - USB (Upper Sideband)
  - LSB (Lower Sideband)
- **Signal Processing:**
  - Automatic Gain Control (AGC)
  - Various filter types (Low-pass, High-pass, Band-pass)
  - Sample rate conversion
- **File Format Support:**
  - WAV file reading and writing
  - Support for various sample rates and durations

## 📦 Installation

```bash
go get github.com/Vivirinter/sdr-parser
```

## 🎯 Quick Start

### Generate Signals

```bash
# Generate AM signal
sdrparser generate -o am_signal.wav -f 440 -d 3.0 -m am

# Generate FM signal
sdrparser generate -o fm_signal.wav -f 440 -d 3.0 -m fm

# Generate USB signal
sdrparser generate -o usb_signal.wav -f 440 -d 3.0 -m usb

# Generate LSB signal
sdrparser generate -o lsb_signal.wav -f 440 -d 3.0 -m lsb
```

### Demodulate Signals

```bash
# Demodulate AM signal
sdrparser demod -i am_signal.wav -o demod_am.wav -t am

# Demodulate FM signal
sdrparser demod -i fm_signal.wav -o demod_fm.wav -t fm

# Demodulate SSB signals
sdrparser demod -i usb_signal.wav -o demod_usb.wav -t usb
sdrparser demod -i lsb_signal.wav -o demod_lsb.wav -t lsb
```

## 🛠️ Command Line Options

### Generate Command
```bash
sdrparser generate [flags]

Flags:
  -o, --output string    Output WAV file (default "signal.wav")
  -f, --freq float      Signal frequency in Hz (default 440.0)
  -d, --duration float  Signal duration in seconds (default 5.0)
  -m, --mod string      Modulation type (am, fm, usb, lsb) (default "am")
```

### Demodulate Command
```bash
sdrparser demod [flags]

Flags:
  -i, --input string    Input WAV file
  -o, --output string   Output WAV file (default "audio.wav")
  -t, --type string     Demodulation type (am, fm, usb, lsb) (default "am")
```

## 🔧 Development

### Project Structure
```
.
├── cmd/sdrparser/     # Main application entry point
├── internal/          # Internal packages
│   └── processing/    # Signal processing implementations
├── pkg/              # Public packages
│   ├── demod/        # Demodulation algorithms
│   └── reader/       # WAV file handling
└── test/             # Test files
```

### Signal Processing Details

#### AM (Amplitude Modulation)
```
s(t) = A * (1 + m(t)) * cos(2π * fc * t)
```
where:
- A: carrier amplitude
- m(t): modulating signal
- fc: carrier frequency

#### FM (Frequency Modulation)
```
s(t) = A * cos(2π * fc * t + β * ∫m(t)dt)
```
where:
- β: frequency deviation
- m(t): modulating signal

#### SSB (Single Sideband)
```
USB: s(t) = A * cos(2π * fc * t) + j * sin(2π * fc * t)
LSB: s(t) = A * cos(2π * fc * t) - j * sin(2π * fc * t)
```

### Build and Test

```bash
# Build the project
make build

# Run tests
make test

# Generate test signals
make generate-test-signals

# Run demodulation tests
make run-demod-tests
```

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
