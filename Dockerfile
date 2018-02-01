ARG GO_VERSION=1.9.3
ARG ALPINE_VERSION=3.7

FROM golang:${GO_VERSION} AS BUILD
WORKDIR /go/src/github.com/frankgreco/demogods/
COPY Gopkg.toml Gopkg.lock Makefile /go/src/github.com/frankgreco/demogods/
RUN make install
COPY ./ /go/src/github.com/frankgreco/demogods/
RUN CGO_ENABLED=0 make binary

FROM alpine:${ALPINE_VERSION}
COPY --from=BUILD /go/src/github.com/frankgreco/demogods/demogods /
ENTRYPOINT ["/demogods"]
