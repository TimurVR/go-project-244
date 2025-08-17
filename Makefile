.PHONY: build test lint run

build:
	@if not exist bin mkdir bin
	go build -o ./bin/gendiff.exe  ./cmd/gendiff/main.go
run: build
	./bin/gendiff.exe $(ARGS)