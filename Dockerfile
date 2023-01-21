ARG GO_VERSION=1.19

FROM golang:${GO_VERSION} as build

WORKDIR /go/src/app
COPY . .

RUN CGO_ENABLED=0 make build

FROM gcr.io/distroless/static-debian11

COPY --from=build /go/src/bin/simulator /
CMD ["simulator"]