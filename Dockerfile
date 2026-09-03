ARG GO_VERSION=1.27.1@sha256:512690a5660563b57d37ecc31129e7f136e831db2aed24a1dbeb8ad7380dc0fa
FROM golang:${GO_VERSION} AS build

ARG TARGETOS
ARG TARGETARCH

ADD . /go/src/github.com/nvalembois/gh-pr-auto-approver
WORKDIR /go/src/github.com/nvalembois/gh-pr-auto-approver

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build \
    -ldflags="-s -w" \
    -o gh-pr-auto-approver cmd/gh-pr-auto-approver.go

FROM scratch
COPY --from=build /go/src/github.com/nvalembois/gh-pr-auto-approver/gh-pr-auto-approver /gh-pr-auto-approver
COPY --from=build /etc/ssl/certs /etc/ssl/certs
ENTRYPOINT ["/gh-pr-auto-approver"]
