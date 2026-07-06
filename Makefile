.PHONY: build run clean test

build:
	go build -o foundry-plans

run: build
	./foundry-plans

clean:
	rm -f foundry-plans

test:
	go test ./...
