FROM golang:1.25-alpine AS builder

WORKDIR /workspace

# Copy workspace files, then trim to only the modules this service needs.
# go.work lists services/search and services/vendor which aren't needed here.
COPY go.work go.work.sum ./
RUN printf 'go 1.25.5\n\nuse (\n\t./app/cmd\n\t./database\n\t./pkg\n)\n' > go.work

# Local dependencies resolved via workspace
COPY pkg/ ./pkg/
COPY database/ ./database/

# go.mod/go.sum first so dependency download is cached separately from source
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
