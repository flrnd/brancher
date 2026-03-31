BINARY=brancher

build:
	go build -o bin/$(BINARY) ./cmd/brancher

run:
	go run ./cmd/brancher $(ARGS)

test:
	go test ./...

clean:
	rm -rf bin/

lint:
	golangci-lint run

fmt:
	go fmt ./...

.PHONY: build run test clean lint fmt
