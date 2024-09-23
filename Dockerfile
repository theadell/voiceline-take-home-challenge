FROM golang:1.23-alpine as builder

WORKDIR /app

RUN apk --no-cache add ca-certificates gcc musl-dev sqlite-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o api ./cmd/api/

FROM alpine:3.19.1

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/openapi3.sepc.yaml /openapi3.sepc.yaml
COPY --from=builder /app/migrations /migrations

COPY --from=builder /app/api /api

ENV HOST=0.0.0.0
ENV PORT=8080
ENV MIGRATE_DB=true

ENTRYPOINT ["/api"]
CMD []

EXPOSE 8080