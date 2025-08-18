.PHONY: build test lint run

build:
	@if not exist bin mkdir bin
	go build -o ./bin/gendiff.exe  ./cmd/gendiff/main.go
	go build -o ./bin/gendiff  ./cmd/gendiff/main.go
lint:
	golangci-lint run
run: build
	./bin/gendiff.exe $(ARGS)