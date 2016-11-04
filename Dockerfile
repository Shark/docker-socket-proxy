FROM scratch
MAINTAINER Felix Seidel <felix@seidel.me>
ADD docker-socket-proxy /
VOLUME "/docker"
ENTRYPOINT ["/docker-socket-proxy"]
