FROM        golang:alpine AS builder
WORKDIR     /app
COPY        . .
COPY        .env .env
RUN         apk add --no-cache ca-certificates openssl
RUN         go mod vendor
RUN         CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-w -s -extldflags "-static"' -o /go/bin/multi-iaas -mod=vendor main.go
EXPOSE      9999
ENTRYPOINT  ["/go/bin/multi-iaas"]
