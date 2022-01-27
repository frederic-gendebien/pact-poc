include ../config.mk

info:
	@echo "group: $(GROUP)"

clean:
	-rm -rf tests/pact
	$(MAKE) -C projection clean
	$(MAKE) -C server clean

docker-build:
	$(MAKE) -C projection docker-build
	$(MAKE) -C server docker-build

docker-kill:
	$(MAKE) -C projection docker-kill
	$(MAKE) -C server docker-kill

swarm-deploy:
	$(MAKE) -C projection swarm-deploy
	$(MAKE) -C server swarm-deploy

swarm-undeploy:
	$(MAKE) -C projection swarm-undeploy
	$(MAKE) -C server swarm-undeploy

swarm-redeploy: swarm-undeploy swarm-deploy

swarm-update: docker-build docker-kill