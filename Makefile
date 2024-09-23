all: db-gen build

db-gen:
	sqlc generate

build:
	go build -o ./bin/api ./cmd/api/

.PHONY: all db-gen build