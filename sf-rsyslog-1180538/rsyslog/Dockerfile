FROM library/alpine:3.21

RUN apk add --no-cache rsyslog
COPY rsyslog.conf /etc/

EXPOSE 514
CMD ["rsyslogd", "-n"]