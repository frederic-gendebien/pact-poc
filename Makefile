include config.mk

.PHONY=clean,build,configure,test

PACT_FOLDERS=application/tests/pact/pacts

info:
	@echo "group: $(GROUP)"

clean:
	go clean -cache -testcache
	-rm $(PACT_FOLDERS)
	$(MAKE) -C infrastructure clean
	$(MAKE) -C application clean

configure: pact-setup deps

pact-setup:
	@echo "--- Installing Pact CLI dependencies"
	cd /opt/; \
	curl -fsSL https://raw.githubusercontent.com/pact-foundation/pact-ruby-standalone/master/install.sh | bash

deps:
	go mod download

compile:
	go build -v ./...

test: pact-publish pact-provider-test

pact-consumers-tests: $(PACT_FOLDERS)

$(PACT_FOLDERS): client-pact-test projection-pact-test

client-pact-test:
	go test -v github.com/frederic-gendebien/pact-poc/application/server/pkg/interfaces/client

projection-pact-test:
	go test -v github.com/frederic-gendebien/pact-poc/application/projection/internal/interfaces/eventbus

server-pact-test:
	go test -v github.com/frederic-gendebien/pact-poc/application/server/internal/interfaces/http
	go test -v github.com/frederic-gendebien/pact-poc/application/server/internal/usecase

pact-publish: $(PACT_FOLDERS)
	@pact-broker publish application/tests/pact/pacts \
		--broker-base-url=$(PACT_BROKER_URL) \
		--broker-token=$(PACT_BROKER_TOKEN) \
		--consumer-app-version=$(VERSION) \
		--tag=main

pact-provider-test: server-pact-test

docker-build:
	$(MAKE) -C infrastructure docker-build
	$(MAKE) -C application docker-build

docker-kill:
	$(MAKE) -C infrastructure docker-kill
	$(MAKE) -C application docker-kill

swarm-setup: swarm-init swarm-network

swarm-init:
	-docker swarm init

swarm-network: swarm-init
	-docker network create --attachable -d overlay napoleongames

swarm-deploy:
	$(MAKE) -C infrastructure swarm-deploy
	$(MAKE) -C application swarm-deploy

swarm-undeploy:
	$(MAKE) -C infrastructure swarm-undeploy
	$(MAKE) -C application swarm-undeploy

swarm-redeploy: swarm-undeploy swarm-deploy
