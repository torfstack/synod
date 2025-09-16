FROM node:22-alpine AS frontend-builder

WORKDIR /opt/kayvault-frontend

COPY frontend/package*.json .

RUN ["npm", "ci"]

COPY frontend .

RUN ["npm", "run", "build"]
RUN ["npm", "prune", "--production"]

FROM golang:1.24.4 AS builder

WORKDIR /opt/kayvault

COPY go.mod go.sum ./
COPY backend backend
COPY 'cmd' 'cmd'
COPY sql sql

RUN CGO_ENABLED=0 go build -o bin/kayvault cmd/main.go

FROM alpine:edge AS runner

COPY --from=frontend-builder /opt/kayvault-frontend/dist static
COPY --from=builder /opt/kayvault/bin/kayvault kayvault
COPY sql/migrations sql/migrations

CMD ["./kayvault"]

FROM runner AS runner-debug

RUN apk add --no-cache delve

CMD ["dlv", "exec", "./kayvault", "--headless", "--listen=:4200", "--api-version=2", "--accept-multiclient", "--continue"]

