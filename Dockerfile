FROM golang:1.13 AS builder
WORKDIR $GOPATH/src/github.com/clovergrp/vault-copy
COPY go.mod go.sum ./
ENV GO111MODULE=on \
    GOPROXY=https://proxy.golang.org
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -ldflags="-s -w" -v ./...

FROM scratch
LABEL maintainer="Igor Diakonov <aidos.tanatos@gmail.com>"
EXPOSE 8080
COPY --from=builder /go/bin/vault-copy /vault-copy
ENTRYPOINT ["/vault-copy"]
