FROM node:22-alpine AS frontend-builder

WORKDIR /opt/kayvault-frontend

COPY frontend/package*.json .

RUN ["npm", "ci"]

COPY frontend .

RUN ["npm", "run", "build"]
RUN ["npm", "prune", "--production"]

FROM golang:1.24.0 AS builder

WORKDIR /opt/kayvault

COPY go.mod go.sum ./
COPY backend backend
COPY cmd cmd
COPY sql sql
COPY --from=frontend-builder /opt/kayvault-frontend/build static

RUN CGO_ENABLED=0 go build -o bin/kayvault cmd/main.go

FROM alpine:3.21.2 AS runner

COPY --from=builder /opt/kayvault/bin/kayvault kayvault
COPY sql/migrations sql/migrations

CMD ["./kayvault"]
