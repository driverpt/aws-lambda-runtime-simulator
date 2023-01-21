ARG GO_VERSION=1.19

FROM golang:${GO_VERSION} as build

WORKDIR /go/src/app
COPY . .

RUN make test
RUN CGO_ENABLED=0 LD_FLAGS="-s -w" make build

FROM gcr.io/distroless/static-debian11

COPY --from=build /go/src/app/bin/simulator /usr/bin/simulator
CMD ["simulator"]