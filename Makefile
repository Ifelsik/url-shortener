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
	gotestsum --format pkgname -- -race -coverprofile=test-cover.temp.out ./...
	touch test-cover.out
	grep -v "mock" test-cover.temp.out >> test-cover.out
	go tool cover -func test-cover.out
	rm test-cover.temp.out test-cover.out

