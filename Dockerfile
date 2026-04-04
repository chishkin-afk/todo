FROM golang:1.25.5 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o todo ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/todo .

COPY --from=builder /app/configs ./configs/
COPY --from=builder /app/docs ./docs/
COPY --from=builder /app/internal/infrastructure/persistence/postgres/migrations ./internal/infrastructure/persistence/postgres/migrations/

EXPOSE 9000

CMD ["./todo"]