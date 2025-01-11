FROM golang:1.23.4 AS builder

WORKDIR /opt/kayvault

COPY go.mod go.sum ./
COPY internal internal
COPY cmd cmd
COPY sql sql

RUN CGO_ENABLED=0 go build -o bin/kayvault-backend cmd/main.go

FROM alpine:3.20.0 AS runner

WORKDIR /opt/kayvault

COPY --from=builder /opt/kayvault/bin/kayvault-backend kayvault-backend
# COPY sql/migrations sql/migrations # TODO add migrations

CMD ["./kayvault-backend"]
