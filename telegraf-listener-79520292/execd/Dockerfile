FROM telegraf:alpine

WORKDIR /usr/src/app
USER root
RUN apk add --no-cache python3 py3-pip
ADD --chown=telegraf execd.py /usr/local/bin/myplugin
USER telegraf