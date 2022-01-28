FROM golang:bullseye

RUN apt update; \
    apt install -y build-essential

RUN cd /opt/; \
    curl -fsSL https://raw.githubusercontent.com/pact-foundation/pact-ruby-standalone/master/install.sh | bash

ENV PACT_HOME="/opt/pact"
ENV PACT_BIN="/opt/pact/bin"
ENV PATH="$PATH:$PACT_BIN"

