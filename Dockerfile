ARG GO_VERSION=1.22.5@sha256:fcae9e0e7313c6467a7c6632ebb5e5fab99bd39bd5eb6ee34a211353e647827a
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
