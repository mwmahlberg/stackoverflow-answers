FROM golang:1.24-alpine3.21 AS build
WORKDIR /usr/local/src/redis-client
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o redis-client .

FROM alpine:3.21
WORKDIR /usr/local/bin
COPY --from=build /usr/local/src/redis-client/redis-client .
ENTRYPOINT ["/usr/local/bin/redis-client"]