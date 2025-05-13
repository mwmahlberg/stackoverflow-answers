FROM golang:1.24-alpine3.21 AS build
WORKDIR /usr/local/src/mqtt-client

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o mqtt-client


FROM alpine:3.21
COPY --from=build --chown=nobody:nogroup --chmod=0550 /usr/local/src/mqtt-client/mqtt-client /usr/local/bin/mqtt-client
USER nobody
ENTRYPOINT ["/usr/local/bin/mqtt-client"]