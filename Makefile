.PHONY: bindata build clean cover install run serve test

bindata:
	go-bindata -pkg assets -o assets/data.go -ignore '.*\.go' -ignore '.DS_Store' -prefix assets assets/*

build:
	go build -o ssssg cmd/ssssg/*.go

clean:
	rm -rf docs/*

cover:
	go tool cover -html=cover.out

install: build
	cp -a ssssg ~/bin

run: build clean
	./ssssg build
	prettier --write "docs/**/*.html"

serve: build
	./ssssg serve

test:
	go test -coverprofile=cover.out -mod=readonly ./...
