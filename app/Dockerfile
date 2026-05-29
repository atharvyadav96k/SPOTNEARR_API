FROM golang:1.25.5-alpine

WORKDIR /app

COPY . .

WORKDIR /app/cmd

RUN go mod download
RUN go build -o server .

EXPOSE 8080

CMD ["./server"]