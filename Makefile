BIN_DIR = bin

db-gen:
	sqlc generate

build:
	go build -o $(BIN_DIR)/api ./cmd/api/

test:
	go test ./... -v -cover

run: build 
	./$(BIN_DIR)/api -migrate-db


.PHONY: db-gen build test run 