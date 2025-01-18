FROM golang:1.23.5 AS builder

WORKDIR /opt/kayvault-backend

COPY go.mod go.sum ./
COPY internal internal
COPY cmd cmd
COPY sql sql

RUN CGO_ENABLED=0 go build -o bin/kayvault-backend cmd/main.go

FROM alpine:3.21.2 AS runner

COPY --from=builder /opt/kayvault-backend/bin/kayvault-backend kayvault-backend
COPY sql/migrations sql/migrations

CMD ["./kayvault-backend"]
