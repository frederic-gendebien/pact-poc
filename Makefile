include config.mk

info:
	@echo "group: $(GROUP)"

clean:
	$(MAKE) -C infrastructure clean
	$(MAKE) -C server clean

configure: pact deps

pact:
	@echo "--- Installing Pact CLI dependencies"
	cd /opt/; \
	curl -fsSL https://raw.githubusercontent.com/pact-foundation/pact-ruby-standalone/master/install.sh | bash

deps:
	go mod download

build:
	go build -v ./...

test:
	go test -v -cover ./...

publish-pacts:
	pact-broker publish tests/pact/pacts \
		--broker-base-url=$(PACT_BROKER_URL) \
		--consumer-app-version=$(VERSION)

docker-build:
	$(MAKE) -C infrastructure docker-build
	$(MAKE) -C server docker-build

docker-kill:
	$(MAKE) -C infrastructure docker-kill
	$(MAKE) -C server docker-kill

swarm-setup: swarm-init swarm-network

swarm-init:
	-docker swarm init

swarm-network: swarm-init
	-docker network create --attachable -d overlay napoleongames

swarm-deploy:
	$(MAKE) -C infrastructure swarm-deploy
	$(MAKE) -C server swarm-deploy

swarm-undeploy:
	$(MAKE) -C infrastructure swarm-undeploy
	$(MAKE) -C server swarm-undeploy

swarm-redeploy: swarm-undeploy swarm-deploy
