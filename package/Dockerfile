FROM alpine:3.5
RUN apk add --no-cache ca-certificates
COPY portainer-agent /usr/bin/
CMD ["portainer-agent", "--config-file", "/config/config.json"]
