.PHONY: build test install clean

build:
	go build -o bin/veo ./cmd/veo

test:
	go test -v ./...

install:
	go install ./cmd/veo

clean:
	rm -rf bin/

run:
	go run ./cmd/veo
