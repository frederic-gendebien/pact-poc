GROUP=pactpoc
VERSION=0.0.1
DOCKER_ID=$$(docker ps | grep $(SERVICE) | awk '{ print $$1 }')
PACT_BROKER_URL=
PACT_BROKER_TOKEN=