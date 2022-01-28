REGISTRY=fredericgendebien
IMAGE=golang
TAG=pact-alpine

login:
	docker login

build:
	docker build -t $(IMAGE):$(TAG) .

push: login build
	docker tag $(IMAGE):$(TAG) $(REGISTRY)/$(IMAGE):$(TAG)
	docker push $(REGISTRY)/$(IMAGE):$(TAG)

shell: build
	docker run -ti $(IMAGE):$(TAG) /bin/sh