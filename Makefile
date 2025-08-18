.PHONY: build test lint run

build:
	@if not exist bin mkdir bin
	go build -o ./bin/gendiff.exe  ./cmd/gendiff/main.go
	go build -o ./bin/gendiff  ./cmd/gendiff/main.go

lint:
	golangci-lint run

test:
	go test -v  ./test/main_test.go

test-with-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

run: build
	./bin/gendiff.exe $(ARGS)