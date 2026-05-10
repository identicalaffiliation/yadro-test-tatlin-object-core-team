.PHONY: run build clean unit_test integration_test lint cover

BIN_NAME=app
PATH_TO_MAIN=./cmd/app/main.go
CONFIG_FILE=config.yml

lint:
	golangci-lint run ./...

unit_test:
	go test ./internal/...

integration_test:
	go test -race ./tests/integration

cover:
	go test -coverprofile=coverage.out -coverpkg=./internal/...,./pkg/... ./...
	go tool cover -func=coverage.out

build:
	go build -o $(BIN_NAME) $(PATH_TO_MAIN)

run:
	./$(BIN_NAME) -c=$(CONFIG_FILE)

clean:
	rm -rf $(BIN_NAME)
	rm -rf coverage.out


