ARG GO_VERSION=1.24.0@sha256:2b1cbf278ce05a2a310a3d695ebb176420117a8cfcfcc4e5e68a1bef5f6354da
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
