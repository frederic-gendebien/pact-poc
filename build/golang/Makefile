REGISTRY=fredericgendebien
IMAGE=golang
TAG=pact-bullseye

login:
	docker login

compile:
	docker build -t $(IMAGE):$(TAG) .

push: login build
	docker tag $(IMAGE):$(TAG) $(REGISTRY)/$(IMAGE):$(TAG)
	docker push $(REGISTRY)/$(IMAGE):$(TAG)

shell: build
	docker run -ti --rm $(IMAGE):$(TAG) /bin/bash