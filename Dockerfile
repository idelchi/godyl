#[=======================================================================[
# Description : Docker image containing the godyl binary
#]=======================================================================]

ARG GO_VERSION=1.26.0
ARG DISTRO=trixie

#### ---- Build ---- ####
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-${DISTRO} AS build

LABEL maintainer=arash.idelchi

ARG TARGETARCH
ARG TARGETOS

USER root

RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    ca-certificates \
    git \
    jq \
    yq \
    nano \
    && rm -rf /var/lib/apt/lists/*

RUN apt-get update && apt-get install -y --no-install-recommends \
    gnome-keyring \
    dbus-x11 \
    && rm -rf /var/lib/apt/lists/*

ARG TASK_VERSION=v3.41.0
RUN wget -qO- https://github.com/go-task/task/releases/download/${TASK_VERSION}/task_linux_${TARGETARCH}.tar.gz | tar -xz -C /usr/local/bin

WORKDIR /work

# Create User (Debian/Ubuntu)
ARG USER=user
ARG UID=1001
RUN groupadd -r -g ${UID} ${USER} && \
    useradd -r -u ${UID} -g ${UID} -m -c "${USER} account" -d /home/${USER} -s /bin/bash ${USER}

USER ${USER}
WORKDIR /tmp/go

ENV GOMODCACHE=/home/${USER}/.cache/.go-mod
ENV GOCACHE=/home/${USER}/.cache/.go

COPY go.mod go.sum ./
RUN --mount=type=cache,target=${GOMODCACHE},uid=1001,gid=1001 \
    --mount=type=cache,target=${GOCACHE},uid=1001,gid=1001 \
    go mod download

RUN go mod download

ARG TARGETOS TARGETARCH

COPY . .
ARG GODYL_VERSION="unofficial & built by unknown"
RUN --mount=type=cache,target=${GOMODCACHE},uid=${UID},gid=${UID},id=gomod-${TARGETARCH} \
    --mount=type=cache,target=${GOCACHE},uid=${UID},gid=${UID},id=gocache-${TARGETARCH} \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 go build -ldflags="-s -w -X 'main.version=${GODYL_VERSION}'" -o bin/ .

RUN go run . || true

ENV PATH=$PATH:/home/${USER}/.local/bin
ENV PATH=$PATH:/root/.local/bin
ENV XDG_RUNTIME_DIR=/tmp/${UID}
ENV XDG_CONFIG_HOME=/home/${USER}/.config
ENV XDG_CACHE_HOME=/home/${USER}/.cache

RUN mkdir -p /home/${USER}/.local/bin
RUN cp bin/godyl /home/${USER}/.local/bin
RUN mkdir -p $XDG_RUNTIME_DIR && chmod 700 $XDG_RUNTIME_DIR

WORKDIR /tmp
USER root
RUN rm -rf /tmp/go

USER ${USER}
WORKDIR /home/${USER}

# Timezone
ENV TZ=Europe/Zurich

FROM debian:trixie-slim AS final

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Create User (Debian/Ubuntu)
ARG USER=user
ARG UID=1001
RUN groupadd -r -g ${UID} ${USER} && \
    useradd -r -u ${UID} -g ${UID} -m -c "${USER} account" -d /home/${USER} -s /bin/bash ${USER}

USER ${USER}
WORKDIR /home/${USER}

COPY --from=build --chown=${USER}:{USER} /home/${USER}/.local/bin/godyl /home/${USER}/.local/bin/godyl

ENV PATH=$PATH:/home/${USER}/.local/bin
ENV PATH=$PATH:/root/.local/bin
ENV XDG_RUNTIME_DIR=/tmp/${UID}
ENV XDG_CONFIG_HOME=/home/${USER}/.config
ENV XDG_CACHE_HOME=/home/${USER}/.cache
