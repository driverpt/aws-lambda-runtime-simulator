GO=go

ifdef LD_FLAGS
LD_FLAGS_ARG=-ldflags "${LD_FLAGS}"
endif

tidy:
	$(GO) mod tidy

build-simulator:
	$(GO) build -o bin/simulator ${LD_FLAGS_ARG} cmd/simulator/main.go

deps:
	$(GO) get ./...

vet:
	$(GO) vet ./...

test: deps
	$(GO) test -race --coverprofile=coverage.coverprofile --covermode=atomic ./...

build: deps vet build-simulator