# syntax=docker.io/docker/dockerfile:1

# Build the Go server
FROM golang:1.25.0-alpine AS go_builder
WORKDIR /app
COPY go.* ./
ENV GOMODCACHE=/.cache/go-mod
RUN --mount=type=cache,target="/.cache/go-mod" go mod download

WORKDIR /app
COPY . .
ENV GOCACHE=/.cache/go-build
RUN --mount=type=cache,target="/.cache/go-mod" \
    --mount=type=cache,target="/.cache/go-build" \
    go build ./cmd/polaris

# Production image
FROM alpine AS runner
WORKDIR /app
COPY --from=go_builder /app/polaris ./polaris

# Start the Go server
CMD ["./polaris"]

