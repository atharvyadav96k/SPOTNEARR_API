FROM golang:1.25-alpine AS builder

WORKDIR /workspace

# Copy workspace descriptor
COPY go.work go.work.sum ./

# Copy go.mod + go.sum for every module listed in go.work.
# go mod download fails if any module referenced in go.work is missing its go.mod.
COPY pkg/go.mod pkg/go.sum ./pkg/
COPY app/cmd/go.mod app/cmd/go.sum ./app/cmd/
COPY services/vendor/go.mod services/vendor/go.sum ./services/vendor/
COPY services/search/go.mod services/search/go.sum ./services/search/

# Download external deps (pkg is a local workspace replace, not downloaded)
RUN cd /workspace/pkg && go mod download && \
    cd /workspace/services/vendor && go mod download

# Copy only the source needed for this binary
COPY pkg/ ./pkg/
COPY services/vendor/ ./services/vendor/

WORKDIR /workspace/services/vendor
RUN go build -o server ./cmd/

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /workspace/services/vendor/server .

EXPOSE 8081

CMD ["./server"]
