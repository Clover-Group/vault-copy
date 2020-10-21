FROM golang:1.15 AS builder
WORKDIR $GOPATH/src/github.com/clovergrp/vault-copy
RUN groupadd -r appuser && useradd -r -g appuser appuser
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates
COPY go.mod go.sum ./
ENV GO111MODULE=on \
    GOPROXY=https://proxy.golang.org
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -ldflags="-s -w" -v ./...

FROM scratch
LABEL maintainer="Igor Diakonov <aidos.tanatos@gmail.com>"
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
USER appuser
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/vault-copy /vault-copy
ENTRYPOINT ["/vault-copy"]
