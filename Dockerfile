ARG GO_VERSION=1.22.4@sha256:4950c1cce16bb754e23ee70d25a67e906149d0faedc0aaeed49c075b696fa889
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
