GROUP=pocpact
DOCKER_ID=$$(docker ps | grep $(SERVICE) | awk '{ print $$1 }')