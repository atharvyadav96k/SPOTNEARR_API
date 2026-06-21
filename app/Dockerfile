FROM golang:1.25-alpine AS builder

WORKDIR /workspace

# Workspace files — must come first so Go can resolve local modules
COPY go.work go.work.sum ./

# Local module dependencies (resolved via go.work)
COPY pkg/ ./pkg/
COPY database/ ./database/

# App go.mod/go.sum first for layer caching
COPY app/cmd/go.mod app/cmd/go.sum ./app/cmd/
RUN cd app/cmd && go mod download

# Full app source
COPY app/cmd/ ./app/cmd/

RUN cd app/cmd && CGO_ENABLED=0 GOOS=linux \
    go build -ldflags="-w -s" -o /server .

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /server /server

EXPOSE 8080

CMD ["/server"]
