GO=go

tidy:
	$(GO) mod tidy

build-simulator:
	$(GO) build -o bin/simulator cmd/simulator/main.go

deps:
	$(GO) get ./...

vet:
	$(GO) vet ./...

test:
	$(GO) test -race --coverprofile=coverage.coverprofile --covermode=atomic ./...

build: deps tidy vet test build-simulator