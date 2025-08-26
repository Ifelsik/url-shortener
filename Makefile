.PHONY: run format lint mockery test

run:
	go run ./cmd/main.go

format:
	go fmt ./...
	@echo Source code formatted!

lint:
	golangci-lint run

mock:
	mockery

test:
	gotestsum --format pkgname -- -race -coverprofile=test-cover.out ./...
	rm test-cover.out
