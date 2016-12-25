# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=main
ENTRYPOINT="cmd/main.go"
BINARY_UNIX=$(BINARY_NAME)_unix

all: build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(ENTRYPOINT)

test:
	$(GOTEST) -v ./...

clean:
	rm -rf atom_data/*
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

cleanall:
	rm -f glide.lock
	rm -f glide.yaml

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

deps:
	rm -r glide.yaml
	rm -f glide.lock
	glide init
	glide install

rockdb:
	$(MAKE) -C $rockdb;

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
docker-build:
	docker run --rm -it -v "$(GOPATH)":/go -w /go/src/bitbucket.org/rsohlich/makepost golang:latest go build -o "$(BINARY_UNIX)" -v    
