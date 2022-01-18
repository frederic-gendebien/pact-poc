include ../config.mk

STACK=infra

docker-build:
	$(MAKE) -C reverseproxy docker-build

swarm-deploy:
	docker stack deploy -c build/docker-compose.yml $(STACK)

swarm-undeploy:
	docker stack rm $(STACK)

swarm-redeploy: swarm-undeploy swarm-deploy