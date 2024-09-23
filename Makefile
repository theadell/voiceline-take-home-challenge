BIN_DIR = bin

all: db-gen build

db-gen:
	sqlc generate

build:
	go build -o $(BIN_DIR)/api ./cmd/api/

run: build 
	./$(BIN_DIR)/api -migrate-db

.PHONY: all db-gen build