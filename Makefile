GO=go

ifdef LD_FLAGS
LD_FLAGS_ARG=-ldflags "${LD_FLAGS}"
endif

ifeq ($(RACE), true)
RACE_ARG=-race
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
	$(GO) test ${RACE_ARG} --coverprofile=coverage.coverprofile --covermode=atomic ./...

build: deps vet build-simulator