FROM golang:1.25-alpine AS builder

WORKDIR /workspace

COPY cmd/go.mod cmd/go.sum ./
RUN go mod download

COPY cmd/ .

RUN CGO_ENABLED=0 GOOS=linux \
    go build -ldflags="-w -s" -o /server .

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /server /server

EXPOSE 8080

CMD ["/server"]
