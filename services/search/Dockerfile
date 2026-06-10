FROM golang:1.25-alpine AS builder

WORKDIR /workspace

# Copy workspace descriptor
COPY go.work go.work.sum ./

# Copy go.mod (+ go.sum where present) for every module listed in go.work.
# database has no go.sum — go.work.sum at the root covers its deps.
COPY database/go.mod ./database/
COPY pkg/go.mod pkg/go.sum ./pkg/
COPY app/cmd/go.mod app/cmd/go.sum ./app/cmd/
COPY services/vendor/go.mod services/vendor/go.sum ./services/vendor/
COPY services/search/go.mod services/search/go.sum ./services/search/

# Download external deps for modules that reach the network.
RUN cd /workspace/pkg && go mod download && \
    cd /workspace/services/search && go mod download

# Copy source needed for this binary.
COPY database/ ./database/
COPY pkg/ ./pkg/
COPY services/search/ ./services/search/

WORKDIR /workspace/services/search
RUN go build -o server ./cmd/

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /workspace/services/search/server .

EXPOSE 8082

CMD ["./server"]
