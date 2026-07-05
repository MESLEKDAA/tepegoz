

FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /app




COPY go.mod go.sum ./
RUN go mod download


COPY . .


RUN go build -o tepegoz cmd/tepegoz/main.go





FROM alpine:latest

WORKDIR /root/


RUN apk add --no-cache tzdata ca-certificates






COPY --from=builder /app/tepegoz .
COPY --from=builder /app/configs ./configs


RUN mkdir -p reports


VOLUME ["/root/logs", "/root/reports"]
CMD ["./tepegoz"]