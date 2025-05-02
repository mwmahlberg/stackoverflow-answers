FROM alpine:latest
ARG DOCKERIZE_VERSION="0.9.3"
RUN apk --no-cache upgrade && apk add --no-cache \
    ca-certificates \
    curl \
    tar \
    gzip \
    && update-ca-certificates
RUN wget https://github.com/jwilder/dockerize/releases/download/v${DOCKERIZE_VERSION}/dockerize-alpine-linux-amd64-v${DOCKERIZE_VERSION}.tar.gz -O - | \
    tar -xvz -C /usr/local/bin dockerize && \
    chown root:root /usr/local/bin/dockerize
# Install your application here
ENTRYPOINT [ "/usr/local/bin/dockerize" ]