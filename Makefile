.PHONY: build run

build:
	go build -o ssssg cmd/ssssg/*.go

run: build
	./ssssg
