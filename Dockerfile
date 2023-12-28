FROM golang:alpine AS builder
# add CA for HTTPS(SSL)
# RUN apk update && apk add git && apk add ca-certificates

ENV GO111MODULE=auto CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main

FROM debian:10.4-slim
# FROM scratch

RUN apt update && apt install tzdata -y
ENV TZ="Asia/Seoul"

# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# COPY --from=builder /etc/passwd /etc/passwd

COPY --from=builder ["/go/src/app/main", "/"]

EXPOSE 3000
CMD ["/main"]