FROM golang:alpine AS builder

RUN apk add make

WORKDIR /build
COPY . .

RUN make -C application/server clean app

FROM alpine
WORKDIR /server
COPY --from=builder /build/application/server/bin/app .

CMD "./app"
