FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY . .

WORKDIR /app/cmd

RUN go mod download
RUN go build -o server .

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/cmd/server .

EXPOSE 8080

CMD ["./server"]