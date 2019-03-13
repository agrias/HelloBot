FROM jrottenberg/ffmpeg:alpine

ADD YmirBot /app/YmirBot
ADD ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["/app/YmirBot"]

EXPOSE 9095

