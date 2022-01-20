include config.mk

info:
	@echo "group: $(GROUP)"

clean:
	go clean -cache -testcache
	$(MAKE) -C infrastructure clean
	$(MAKE) -C application clean

configure: pact deps

pact:
	@echo "--- Installing Pact CLI dependencies"
	cd /opt/; \
	curl -fsSL https://raw.githubusercontent.com/pact-foundation/pact-ruby-standalone/master/install.sh | bash

deps:
	go mod download

build:
	go build -v ./...

test: publish-pacts
	go test -v -cover ./...

/application/tests/pact/pacts:
	go test -v github.com/frederic-gendebien/pact-poc/application/server/pkg/interfaces/client

publish-pacts: /application/tests/pact/pacts
	@pact-broker publish application/tests/pact/pacts \
		--broker-base-url=$(PACT_BROKER_URL) \
		--broker-token=$(PACT_BROKER_TOKEN) \
		--consumer-app-version=$(VERSION) \
		--tag=main

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
