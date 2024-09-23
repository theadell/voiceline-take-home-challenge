BIN_DIR = bin

mock:
	mockgen -source=internal/db/querier.go -destination=internal/db/mock/querier_mock.go -package=mock

db-gen:
	sqlc generate

build:
	go build -o $(BIN_DIR)/api ./cmd/api/

test:
	go test ./... -v -cover

run: build 
	./$(BIN_DIR)/api -migrate-db


.PHONY: mock db-gen build test run 