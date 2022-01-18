include ../../config.mk

SERVICE=reverseproxy

docker-build:
	docker build -t $(GROUP)/$(SERVICE):latest build

docker-deploy:
	docker run -d --rm --name $(SERVICE) --publish "8081:80" $(GROUP)/$(SERVICE)

docker-undeploy:
	docker rm --force $(SERVICE)

docker-shell:
	docker exec -ti $(DOCKER_ID) sh

docker-logs:
	docker logs -f $(DOCKER_ID)

docker-kill:
	docker rm --force $(DOCKER_ID)
