# syntax=docker/dockerfile:1
# GO_VERSION is the version of Go to use for building all dependencies.
# This is set via the --build-arg flag when running `docker build`.
ARG GO_VERSION
ARG DEBIAN_VERSION=bookworm
FROM golang:${GO_VERSION}-${DEBIAN_VERSION} AS builder
WORKDIR /src

RUN mkdir -p /artifacts

# Copy the go module and download the dependencies, only rebuild if the
# go.mod or go.sum files change.
# TODO(jaredallard): Buildkit cache mount.
COPY go.mod go.sum ./
RUN go mod download

# Copy over the rest
COPY . .

RUN go build -trimpath -v -ldflags "-s -w" -o /artifacts ./cmd/...

FROM golang:${GO_VERSION}-${DEBIAN_VERSION} AS factocord
ENV GOBIN=/artifacts
RUN go install github.com/maxsupermanhd/FactoCord-3.0/v3@v3.2.19

FROM debian:${DEBIAN_VERSION}-slim
ARG VERSION=latest
ARG SHA256SUM=""
VOLUME /factorio
WORKDIR /factorio
EXPOSE 34197/udp 27015/tcp
SHELL ["/bin/bash", "-euox", "pipefail", "-c"]

ARG USER=factorio
ARG GROUP=factorio
ARG PUID=845
ARG PGID=845

ENV SAVES=/factorio/saves \
  CONFIG=/factorio/config \
  MODS=/factorio/mods \
  SCENARIOS=/factorio/scenarios

# Ensure CA certificates are up to date.
#
# Why: ca-certificates don't need to be pinned.
# hadolint ignore=DL3008
RUN apt-get update -y\
  && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends ca-certificates \
  && rm -rf /var/lib/apt/lists/*

# Create a user to run the Factorio server under.
RUN addgroup --system --gid "${PGID}" "${GROUP}" \
  && adduser --system --uid "${PUID}" --gid "${PGID}" --no-create-home --disabled-password --shell /bin/sh "${USER}"

# Copy over the built binaries from earlier.
COPY --from=factocord /artifacts/ /usr/local/bin/
COPY --from=builder /artifacts/ /usr/local/bin/

# Download Factorio
RUN downloader --version="${VERSION}" --sha256sum="${SHA256SUM}" "/opt/factorio" \
  && chmod ugo=rwx "/opt/factorio"

# Default config
COPY docker/default-factorio-config.ini /opt/factorio/config/config.ini
COPY docker/default-factocord.json /etc/factorio-docker/

USER $USER