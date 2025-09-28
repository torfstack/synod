FROM node:22-alpine AS frontend-builder

WORKDIR /opt/synod-frontend

COPY frontend/package*.json .

RUN ["npm", "ci"]

COPY frontend .

RUN ["npm", "run", "build"]
RUN ["npm", "prune", "--production"]

FROM golang:1.25.1 AS builder

WORKDIR /opt/synod

COPY go.mod go.sum ./
COPY backend backend
COPY 'cmd' 'cmd'
COPY sql sql

RUN CGO_ENABLED=0 go build -o bin/synod cmd/main.go

FROM alpine:edge AS runner

COPY --from=frontend-builder /opt/synod-frontend/dist static
COPY --from=builder /opt/synod/bin/synod synod
COPY sql/migrations sql/migrations

CMD ["./synod"]

FROM runner AS runner-debug

RUN apk add --no-cache delve

CMD ["dlv", "exec", "./synod", "--headless", "--listen=:4200", "--api-version=2", "--accept-multiclient", "--continue"]

