FROM haproxy:lts-alpine

USER root
RUN apk add --no-cache socat

USER haproxy

COPY haproxy.cfg /usr/local/etc/haproxy/haproxy.cfg