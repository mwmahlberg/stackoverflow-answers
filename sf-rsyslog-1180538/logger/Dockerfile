FROM golang:1.24.2-alpine3.21 AS builder
WORKDIR /app
COPY go.mod go.sum main.go ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o log-generator .
FROM alpine:3.21
COPY --from=builder /app/log-generator /usr/local/bin/log-generator
ENTRYPOINT [ "/usr/local/bin/log-generator" ]