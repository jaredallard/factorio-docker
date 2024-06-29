# syntax=docker/dockerfile:1
# GO_VERSION is the version of Go to use for building all dependencies.
# This is set via the --build-arg flag when running `docker build`.
ARG GO_VERSION
ARG DEBIAN_VERSION=bookworm
FROM golang:${GO_VERSION}-${DEBIAN_VERSION} AS builder
WORKDIR /src

# Used for build caching.
ENV GOCACHE=/tmp/.cache/go-build
ENV GOMODCACHE=/tmp/.cache/go-mod

# See --from reasoning near the end.
ENV \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

RUN mkdir -p /artifacts

# Copy the go module and download the dependencies, only rebuild if the
# go.mod or go.sum files change.
# TODO(jaredallard): Buildkit cache mount.
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/tmp/.cache \
  go mod download

# Copy over the rest
COPY . .

RUN --mount=type=cache,target=/tmp/.cache \
  go build -trimpath -v -ldflags "-s -w" -o /artifacts ./cmd/...

FROM golang:${GO_VERSION}-${DEBIAN_VERSION} AS factocord
ENV \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

# Why: We're building Factocord-3.0...
# hadolint ignore=DL3003
RUN --mount=type=cache,target=/tmp/.cache \
  git clone https://github.com/maxsupermanhd/FactoCord-3.0 && \
  cd "FactoCord-3.0" && \
  git checkout v3.2.19 && \
  go mod download && \
  go build -o /artifacts/ .

# We hardcode linux/amd64 here because Factorio can only be ran on
# amd64.
# hadolint ignore=DL3029
FROM --platform=linux/amd64 debian:${DEBIAN_VERSION}-slim
ENTRYPOINT [ "/usr/local/bin/wrapper" ]
ARG VERSION=stable
ARG SHA256SUM=""
VOLUME /data /opt/factorio
WORKDIR /data
EXPOSE 34197/udp 27015/tcp
SHELL ["/bin/bash", "-euox", "pipefail", "-c"]

ARG USER=factorio
ARG GROUP=factorio
ARG PUID=845
ARG PGID=845

# Ensure CA certificates are up to date.
#
# Why: ca-certificates don't need to be pinned.
# hadolint ignore=DL3008
RUN apt-get update -y \
  && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends ca-certificates \
  && rm -rf /var/lib/apt/lists/*

# Create a user to run the Factorio server under.
RUN addgroup --system --gid "${PGID}" "${GROUP}" \
  && adduser --system --uid "${PUID}" --gid "${PGID}" --no-create-home --disabled-password --shell /bin/sh "${USER}" \
  && mkdir -p /data /opt/factorio \
  && chown -R "${USER}:${GROUP}" /data /opt/factorio \
  && chown -R "${USER}:${GROUP}" /data

# Copy over the built binaries from earlier.
COPY --from=factocord /artifacts/ /usr/local/bin/
COPY --from=builder /artifacts/ /usr/local/bin/

USER $USER