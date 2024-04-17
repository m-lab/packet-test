FROM golang:1.22-alpine AS build
RUN apk add gcc git linux-headers musl-dev
WORKDIR /packet-test

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./ ./
RUN ./build.sh

FROM alpine:3.19.1
WORKDIR /packet-test
COPY --from=build /packet-test/server /packet-test/

ENTRYPOINT ["./server"]