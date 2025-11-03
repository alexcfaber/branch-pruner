.PHONY: test lint build fmt ci-check

test:
	@echo "running unit tests..."
	go test ./...

lint:
	@echo "running golangci-lint..."
	golangci-lint run

build:
	@echo "building..."
	go build ./...

fmt:
	@echo "formatting..."
	gofmt -w .

ci-check: fmt lint test


