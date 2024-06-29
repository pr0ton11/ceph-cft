FROM docker.io/golang:1.22.4

LABEL org.opencontainers.image.title=ceph-cft
LABEL org.opencontainers.image.description="Configure your ceph cluster with environment variables"
LABEL org.opencontainers.image.version=0.1.0
LABEL org.opencontainers.image.licenses=GPL-3.0
LABEL org.opencontainers.image.url=https://github.com/pr0ton11/ceph-cft
LABEL org.opencontainers.image.source=https://github.com/pr0ton11/ceph-cft
LABEL org.opencontainers.image.authors=pr0ton11

ENV CGO_ENABLED=0
ENV GOOS=linux

COPY . /app

WORKDIR /app

RUN go mod download && \
    go build -ldflags="-s -w" -o /app/ceph-cft && \
    chmod +x /app/ceph-cft && \
    cp /app/ceph-cft /usr/local/bin/ceph-cft

ENTRYPOINT [ "sh", "-c", "ceph-cft" ]
