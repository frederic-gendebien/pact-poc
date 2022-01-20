include ../config.mk

info:
	@echo "group: $(GROUP)"

clean:
	-rm -rf tests/pact
	$(MAKE) -C server clean

docker-build:
	$(MAKE) -C server docker-build

docker-kill:
	$(MAKE) -C server docker-kill

swarm-deploy:
	$(MAKE) -C server swarm-deploy

swarm-undeploy:
	$(MAKE) -C server swarm-undeploy

swarm-redeploy: swarm-undeploy swarm-deploy
