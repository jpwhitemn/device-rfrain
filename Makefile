.PHONY: build run clean 

GO = CGO_ENABLED=0 GO111MODULE=on go

MICROSERVICES=cmd/device-rfrain

.PHONY: $(MICROSERVICES)

VERSION=$(shell cat ./VERSION 2>/dev/null || echo 0.0.0)
GIT_SHA=$(shell git rev-parse HEAD)
GOFLAGS=-ldflags "-X github.com/edgexfoundry/device-rfrain.Version=$(VERSION)"

build: $(MICROSERVICES)

cmd/device-rfrain:
	$(GO) build $(GOFLAGS) -o $@ ./cmd

run:
	cd cmd && ./device-rfrain

clean:
	rm -f $(MICROSERVICES)

