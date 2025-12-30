FROM golang:1.24-alpine

RUN apk add --no-cache git bash

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 安装 air
RUN go install github.com/cosmtrek/air@latest

CMD ["air", "-c", ".air.toml"]
