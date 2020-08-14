.PHONY: bindata build cover install run test

bindata:
	go-bindata -pkg assets -o assets/data.go -ignore '.*\.go' -ignore '.DS_Store' -prefix assets assets/*

build:
	go build -o ssssg cmd/ssssg/*.go

cover:
	go tool cover -html=cover.out

install: build
	cp -a ssssg ~/bin

run: build
	./ssssg build
	prettier --write "docs/**/*.html"

test:
	go test -coverprofile=cover.out -mod=readonly ./...
