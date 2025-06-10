# Build
FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.22-alpine AS builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

ARG Version
ARG GitCommit

ENV CGO_ENABLED=0
ENV GO111MODULE=on
WORKDIR /go/src/github.com/rverchere/ovh-mks-exporter

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY internal internal
COPY cmd cmd

RUN CGO_ENABLED=${CGO_ENABLED} GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
  go build -v \
  ./cmd/ovh-mks-exporter

# Release
FROM --platform=${TARGETPLATFORM:-linux/amd64} gcr.io/distroless/static:nonroot

LABEL maintainer="Rémi Verchère <remi@verchere.fr>"

ENV OVH_ENDPOINT="ovh-eu"
ENV OVH_APPLICATION_KEY=""
ENV OVH_APPLICATION_SECRET=""
ENV OVH_CONSUMER_KEY=""
ENV OVH_CLOUD_PROJECT_SERVICE=""
ENV OVH_CLOUD_PROJECT_KUBEID=""

WORKDIR /

EXPOSE 9101

COPY --from=builder /go/src/github.com/rverchere/ovh-mks-exporter/ovh-mks-exporter /ovh-mks-exporter

USER nonroot:nonroot

CMD [ "/ovh-mks-exporter" ]
