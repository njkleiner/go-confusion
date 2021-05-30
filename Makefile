.PHONY: all format lint test

all: format test

format:
	go fmt ./...

lint:
	go vet ./...

	gofmt -d ./ | diff -u /dev/null -

test: lint
	go test -cover ./...

