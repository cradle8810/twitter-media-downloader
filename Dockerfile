FROM golang:trixie AS builder

ENV DEBIAN_FRONTEND=noninteractive

WORKDIR /src
COPY * /src
RUN make

#----------------------------------------
FROM ubuntu:24.04
LABEL org.opencontainers.image.source=https://github.com/cradle8810/twitter-media-downloader

RUN apt update && \
    apt install -y \
      ca-certificates \
      openssl \
    && \
    rm -rf /var/lib/apt/*

WORKDIR /opt/bin
COPY --from=builder --chmod=755 /src/twmd /opt/bin/twmd

WORKDIR /work
ENTRYPOINT ["/opt/bin/twmd"]
CMD ["-h"]
