GO=go

build-simulator:
	$(GO) build cmd/simulator/main.go -o bin/simulator

vet:
	$(GO) vet -v

test:
	$(GO) test -race --coverprofile=coverage.coverprofile --covermode=atomic ./...

build: vet test build-simulator