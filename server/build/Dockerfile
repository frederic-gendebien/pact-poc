FROM golang:alpine AS builder

RUN apk add make

WORKDIR /build
COPY . .

RUN make -C server clean app

FROM alpine
WORKDIR /server
COPY --from=builder /build/server/bin/app .

CMD "./app"
