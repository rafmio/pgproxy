FROM golang:alpine AS builder

WORKDIR /app

COPY . .

RUN go mod init pgproxy && \
    go get github.com/lib/pq && \
    go get . && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app

FROM alpine:latest

COPY --from=builder /app/app /usr/local/bin/app
# COPY migrations /app/migrations

CMD ["/usr/local/bin/app"]