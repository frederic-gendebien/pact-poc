FROM golang:alpine

RUN apk update; \
    apk add make; \
    apk add curl

RUN cd /opt/; \
    curl -fsSL https://raw.githubusercontent.com/pact-foundation/pact-ruby-standalone/master/install.sh | sh

ENV PACT_HOME="/opt/pact"
ENV PACT_BIN="/opt/pact/bin"
ENV PATH="$PATH:$PACT_BIN"

