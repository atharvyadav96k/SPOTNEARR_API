FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /workspace

COPY go.work ./
RUN printf 'go 1.25.5\n\nuse (\n\t./app/cmd\n\t./database\n\t./pkg\n)\n\nreplace (\n\tgithub.com/atharvyadav96k/spotnearr/pkg v0.0.0 => ./pkg\n\tgithub.com/Developer-Aadesh/spotnearr-database v0.0.0 => ./database\n)\n' > go.work

COPY pkg/ ./pkg/
COPY database/ ./database/
COPY app/cmd/ ./app/cmd/

ENV GONOSUMDB=*
ENV GOPRIVATE=github.com/atharvyadav96k/*,github.com/Developer-Aadesh/*

RUN cd app/cmd && CGO_ENABLED=0 GOOS=linux \
    go build -ldflags="-w -s" -o /server .

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /server /server

EXPOSE 8080

CMD ["/server"]
