FROM oven/bun:1.3.9-alpine AS frontend-builder

WORKDIR /opt/synod-frontend

COPY frontend/package.json frontend/bun.lock ./
RUN --mount=type=cache,target=/root/.bun/install/cache \
    bun install --frozen-lockfile

COPY frontend/ .

RUN bun run build

FROM golang:1.27.1-alpine AS builder-base

WORKDIR /opt/synod

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

FROM builder-base AS builder-prod
COPY backend/ backend/
COPY sql/ sql/
COPY cmd/ cmd/

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -ldflags="-s -w" -o bin/synod cmd/main.go

FROM builder-base AS builder-debug
COPY backend/ backend/
COPY sql/ sql/
COPY cmd/ cmd/

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -gcflags="all=-N -l" -o bin/synod-debug cmd/main.go

FROM cgr.dev/chainguard/wolfi-base AS runner-base

WORKDIR /app

COPY --from=frontend-builder /opt/synod-frontend/dist ./static
COPY sql/migrations/ ./sql/migrations/

USER 65532:65532

FROM runner-base AS runner
COPY --from=builder-prod /opt/synod/bin/synod ./synod

EXPOSE 8080
CMD ["./synod"]

FROM runner-base AS runner-debug

USER root
RUN apk add --no-cache delve

USER 65532:65532

COPY --from=builder-debug /opt/synod/bin/synod-debug ./synod
EXPOSE 8080 4200
CMD ["dlv", "exec", "./synod", "--headless", "--listen=:4200", "--api-version=2", "--accept-multiclient", "--continue"]

