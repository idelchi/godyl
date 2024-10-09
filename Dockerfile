#[=======================================================================[
# Description : Docker image containing the godyl binary
#]=======================================================================]

ARG GO_VERSION=1.23.1
ARG DISTRO
#### ---- Build ---- ####
FROM golang:${GO_VERSION}-${DISTRO} AS build

LABEL maintainer=arash.idelchi

USER root

RUN apt-get update && apt-get install -y \
    python3 \
    python3-pip \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /work

# Create User (Debian/Ubuntu)
ARG USER=user
RUN groupadd -r -g 1001 ${USER} && \
    useradd -r -u 1001 -g 1001 -m -c "${USER} account" -d /home/${USER} -s /bin/bash ${USER}

USER ${USER}
WORKDIR /tmp/go

ENV GOMODCACHE=/home/${USER}/.cache/.go-mod
ENV GOCACHE=/home/${USER}/.cache/.go

COPY go.mod go.sum ./
RUN --mount=type=cache,target=${GOMODCACHE},uid=1001,gid=1001 \
    --mount=type=cache,target=${GOCACHE},uid=1001,gid=1001 \
    go mod download

RUN go mod download

ENV PATH=$PATH:/home/${USER}/.local/bin
ENV PATH=$PATH:/root/.local/bin

COPY . .
ARG GODYL_VERSION="unofficial & built by unknown"
RUN --mount=type=cache,target=${GOMODCACHE},uid=1001,gid=1001 \
    --mount=type=cache,target=${GOCACHE},uid=1001,gid=1001 \
    CGO_ENABLED=0 go install -ldflags="-s -w -X 'main.version=${GODYL_VERSION}'" ./cmd/...

# Timezone
ENV TZ=Europe/Zurich

COPY .bashrc /home/${USER}/.bashrc


FROM build AS final


USER root

RUN rm -rf /usr/local/go

USER user
