FROM golang:1.15-alpine as builder
WORKDIR /data
ADD . .
RUN go build -o harbor-img-clear

FROM alpine:3.12.0
WORKDIR /data
COPY --from=builder /data/harbor-img-clear .
