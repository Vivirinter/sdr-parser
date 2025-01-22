BUILD_DIR = build
BINARY_NAME = sdrparser

.PHONY: all build clean test

all: clean build test

build:
	go build -o $(BINARY_NAME) ./cmd/sdrparser

clean:
	rm -f $(BINARY_NAME)
	rm -f *.wav

test:
	go test -v ./...

generate-test-signals: build
	./$(BINARY_NAME) generate -o am_test.wav -f 440 -d 3.0 -m am
	./$(BINARY_NAME) generate -o fm_test.wav -f 440 -d 3.0 -m fm

run-demod-tests: generate-test-signals
	./$(BINARY_NAME) demod -i am_test.wav -o demod_audio.wav -t am
	./$(BINARY_NAME) demod -i fm_test.wav -o demod_fm.wav -t fm
