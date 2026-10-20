ARG GO_VERSION=1.27.3@sha256:957861185dbf29d93668ace41e011d091c6e779953d47f0316ef2cc58e28fe2b
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
