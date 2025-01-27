# SDR Parser

Software Defined Radio (SDR) signal processing tool with support for various modulation types and signal processing capabilities.

## 🚀 Features

- **Signal Generation & Demodulation:**
  - AM (Amplitude Modulation)
  - FM (Frequency Modulation)
  - USB (Upper Sideband)
  - LSB (Lower Sideband)
- **Signal Processing:**
  - Automatic Gain Control (AGC)
  - Noise Reduction Filters (Moving Average, Median, Butterworth)
  - WAV file support with various sample rates

## 📦 Installation

```bash
go get github.com/Vivirinter/sdr-parser
```

## 🎯 Quick Start

### Generate Signals

```bash
# Generate AM signal (AM radio station at 1000 kHz)
sdrparser generate -o am_signal.wav -f 1000000 -d 10.0 -m am

# Generate FM signal (FM radio station at 100.5 MHz)
sdrparser generate -o fm_signal.wav -f 100500000 -d 10.0 -m fm

# Generate USB/LSB signals (Amateur radio)
sdrparser generate -o usb_signal.wav -f 14200000 -d 10.0 -m usb
sdrparser generate -o lsb_signal.wav -f 7100000 -d 10.0 -m lsb
```

Common frequencies:
- AM Radio: 530 kHz - 1.7 MHz
- FM Radio: 88 MHz - 108 MHz
- Amateur Radio: 7 MHz, 14 MHz, 21 MHz
- Weather Radio: 162 MHz
- ADS-B (Aircraft): 1090 MHz

### Apply Filters

```bash
# Apply Moving Average filter
sdrparser filter -i input.wav -o filtered.wav -t moving_average -w 5

# Apply Median filter
sdrparser filter -i input.wav -o filtered.wav -t median -w 5

# Apply Butterworth filter
sdrparser filter -i input.wav -o filtered.wav -t butterworth -c 1000 -n 4
```

Filter parameters:
- `-i/--input`: Input WAV file
- `-o/--output`: Output WAV file
- `-t/--type`: Filter type (moving_average, median, butterworth)
- `-w/--window`: Window size for MA/median filters
- `-c/--cutoff`: Cutoff frequency for Butterworth
- `-n/--order`: Filter order for Butterworth
- `--amplitude`: Signal amplitude scaling
- `--snr`: Signal-to-noise ratio (dB)
- `--normalize`: Normalize output signal

## 🧪 Testing

```bash
go test ./...
```

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.
