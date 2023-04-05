# syntax=docker/dockerfile:experimental
FROM golang AS builder

LABEL maintainer="https://github.com/AkronimBlack"

WORKDIR /app

COPY . .

RUN GO111MODULE=on go mod download -x \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o process-manager main.go


FROM alpine

RUN apk update && \
    apk add --no-cache bash ca-certificates dumb-init gettext tzdata && \
    rm -rf /var/cache/apk/* && \
    # update certificates
    update-ca-certificates

COPY start.json .
COPY --from=builder /app/process-manager .

EXPOSE 8080

ENTRYPOINT ["./process-manager"]
