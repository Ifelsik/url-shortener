run:
	go run ./cmd/main.go

format:
	go fmt ./...
	@echo Source code formatted!

lint:
	golangci-lint run
