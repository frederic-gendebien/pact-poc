# Pact POC

## Introduction

Little proof-of-concept using Golang, Pact.io and PactFlow.io

## Setup

Makefile is used for all the different steps, make sure you know what it is! What? Why Makefile??? Because 
it is a unix standard and because all projects should use it. If you are new to something, it helps a lot to have it! 

### Docker

Make sure Docker is installed and Docker Swarm compliant. What? Why Docker Swarm? Because Docker Compose Runtime Sucks!

### PactFlow.io
This project uses PactFlow.io as the Pact Broker, you should create your free account and generate a token.

### Environment
The project uses environment variables. Make sure to create those variables:

```shell
export PACT_BROKER_URL=https://<account>.pactflow.io
export PACT_BROKER_TOKEN=<token>
```

## Build

### Clean
```shell
make clean # Yeah, simple
```

### Compile the Code
```shell
make compile # Happy not to dig into golang directives right?
```

### Test with Pact

#### Install Pact Binaries
You should have the pact binaries in order to execute the tests.

```shell
sudo make pact-setup # Yeah, event that
```

#### Complete Cycle
```shell
make -e test #DO NOT FORGET '-e' AND YOUR ENV VARIABLES
```

#### Test Pact Consumers
```shell
make pact-consumers-tests # Without Makefile here you would jump in the closest river
```

#### Publish pacts
```shell
make pact-publish # 
```

#### Test Pact Provider
```shell
make pact-provider-test # Same Here
```

## Docker Swarm

### Setup

1. Create a Swarm
2. Create a Network

```shell
make swarm-setup
```
### Build

```shell
make docker-build
```

### Run

#### Deploy to Swarm

```shell
make swarm-deploy
```

#### Undeploy to Swarm

```shell
make swarm-undeploy # Whaaaat, I really deploy everything with this? Yep!
```

### Update the swarm

```shell
make -C application swarm-update # Whaaaat, I really update every running services in the swarm with this command? Yep!
```

Or

```shell
make -C application/<service> swarm-update # Whaaaaat, I update a specific service running in the swarm? Yep! Got it it now?
```

### Explore

There are more makefile features, just read them!