.PHONY: build test clean

build:
	go build -o llmem .

test:
	go test ./... -v

clean:
	rm -f llmem
