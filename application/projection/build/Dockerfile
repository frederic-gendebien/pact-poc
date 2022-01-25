FROM golang:alpine AS builder

RUN apk add make

WORKDIR /build
COPY . .

RUN make -C application/projection clean app

FROM alpine
WORKDIR /projection
COPY --from=builder /build/application/projection/bin/app .

CMD "./app"
