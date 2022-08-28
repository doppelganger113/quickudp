SHELL = /bin/bash
PROJECTNAME=quickudp

GOBASE=$(shell pwd)
GOBIN=bin

default:
	go build -i -v -o $(GOBIN)/$(PROJECTNAME) ./cmd/$(PROJECTNAME)/main.go || exit

install:
	go mod download

benchmark:
	CGO_ENABLED=1 go build -race -o $(GOBIN)/benchmark ./cmd/benchmark/main.go

example:
	CGO_ENABLED=1 go build -race -o $(GOBIN)/example ./cmd/example/main.go

start:
	go build -i -v -o $(GOBIN)/$(PROJECTNAME) ./cmd/$(PROJECTNAME)/main.go|| exit
	$(GOBIN)/$(PROJECTNAME) || exit

verify:
	go test ./... -race

cover:
	go test ./... -race -coverprofile="c.out" && go tool cover -func=c.out

cover-html:
	go test ./... -race -coverprofile="c.out" && go tool cover -html=c.out

vet:
	go vet ./...

tidy:
	go mod tidy
