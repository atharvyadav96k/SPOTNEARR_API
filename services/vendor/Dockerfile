FROM golang:1.25-alpine AS builder

WORKDIR /workspace

# Copy workspace files first for better layer caching
COPY go.work go.work.sum ./

# Copy module manifests so `go work download` can pull external deps
COPY pkg/go.mod pkg/go.sum ./pkg/
COPY services/vendor/go.mod services/vendor/go.sum ./services/vendor/

# Download all external dependencies declared across the workspace
RUN go work download

# Copy source
COPY pkg/ ./pkg/
COPY services/vendor/ ./services/vendor/

# Build
WORKDIR /workspace/services/vendor
RUN go build -o server ./cmd/

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /workspace/services/vendor/server .

EXPOSE 8081

CMD ["./server"]
