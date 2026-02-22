FROM oven/bun:1.3.9-alpine AS frontend-builder

WORKDIR /opt/synod-frontend

COPY frontend/package.json frontend/bun.lock ./

RUN bun install --frozen-lockfile

COPY frontend .

RUN bun run build

FROM golang:1.26.0-alpine AS builder

WORKDIR /opt/synod

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY 'cmd' 'cmd'
COPY backend backend
COPY sql sql

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/synod cmd/main.go

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 go build -gcflags="all=-N -l" -o bin/synod-debug cmd/main.go

FROM alpine:3.23 AS runner

WORKDIR /app

COPY --from=builder /opt/synod/bin/synod ./synod
COPY --from=frontend-builder /opt/synod-frontend/dist ./static
COPY sql/migrations ./sql/migrations

CMD ["./synod"]

FROM alpine:3.23 AS runner-debug

RUN apk add --no-cache delve

WORKDIR /app

COPY --from=builder /opt/synod/bin/synod-debug ./synod
COPY --from=frontend-builder /opt/synod-frontend/dist ./static
COPY sql/migrations ./sql/migrations

CMD ["dlv", "exec", "./synod", "--headless", "--listen=:4200", "--api-version=2", "--accept-multiclient", "--continue"]

