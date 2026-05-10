.PHONY: run build clean unit_test integration_test lint

BIN_NAME=app
PATH_TO_MAIN=./cmd/app/main.go
CONFIG_FILE=config.yml

lint:
	golangci-lint run ./...

unit_test:
	go test ./...

build:
	go build -o $(BIN_NAME) $(PATH_TO_MAIN)

run:
	./$(BIN_NAME) -c=$(CONFIG_FILE)

clean:
	rm -rf $(BIN_NAME)
