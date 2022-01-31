include ../../config.mk

.PHONY: info clean build run docker-build docker-deploy docker-shell docker-undeploy

SERVICE=server
STACK=$(SERVICE)

info:
	@echo "group: $(GROUP)"
	@echo "service: $(SERVICE)"

clean:
	-rm bin/app

build:
	go build -v ./...

bin/app:
	-mkdir bin
	go build -v -o bin/app cmd/main.go

app: bin/app

run:
	go run cmd/main.go

docker-build:
	docker build  -t $(GROUP)/$(SERVICE):latest -f build/Dockerfile ../..

docker-deploy:
	docker run -d --rm \
		--name $(SERVICE) \
		--publish "8080:8080" \
		$(GROUP)/$(SERVICE):latest

docker-undeploy:
	docker rm --force $(SERVICE)

docker-redeploy: docker-undeploy docker-deploy

docker-shell:
	docker exec -ti $(DOCKER_ID) sh

docker-logs:
	docker logs -f $(DOCKER_ID)

docker-kill:
	docker rm --force $(DOCKER_ID)

swarm-deploy:
	docker stack deploy -c build/docker-compose.yml $(STACK)

swarm-undeploy:
	docker stack rm $(STACK)

swarm-redeploy: swarm-undeploy swarm-deploy

swarm-update: docker-build docker-kill
