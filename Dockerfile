#[=======================================================================[
# Description : Docker image containing the godyl binary
#]=======================================================================]

ARG GO_VERSION=1.24.2
ARG DISTRO=bookworm

#### ---- Build ---- ####
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-${DISTRO} AS build

LABEL maintainer=arash.idelchi

ARG TARGETARCH
ARG TARGETOS

USER root

RUN apt-get update && apt-get install -y \
    curl \
    ca-certificates \
    git \
    jq \
    yq \
    && rm -rf /var/lib/apt/lists/*

ARG TASK_VERSION=v3.41.0
RUN wget -qO- https://github.com/go-task/task/releases/download/${TASK_VERSION}/task_linux_${TARGETARCH}.tar.gz | tar -xz -C /usr/local/bin

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

ARG TARGETOS TARGETARCH

COPY . .
ARG GODYL_VERSION="unofficial & built by unknown"
RUN --mount=type=cache,target=${GOMODCACHE},uid=1001,gid=1001 \
    --mount=type=cache,target=${GOCACHE},uid=1001,gid=1001 \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 go build -ldflags="-s -w -X 'main.version=${GODYL_VERSION}'" -o bin/ .

RUN go run . || true

ENV PATH=$PATH:/home/${USER}/.local/bin
ENV PATH=$PATH:/root/.local/bin

RUN mkdir -p /home/${USER}/.local/bin
RUN cp bin/godyl /home/${USER}/.local/bin

WORKDIR /home/${USER}

USER root
RUN rm -rf /tmp/go
USER ${USER}

# Timezone
ENV TZ=Europe/Zurich

FROM debian:bookworm-slim AS final

RUN apt-get update && apt-get install -y \
    curl \
    ca-certificates \
    git \
    jq \
    && rm -rf /var/lib/apt/lists/*

# Create User (Debian/Ubuntu)
ARG USER=user
RUN groupadd -r -g 1001 ${USER} && \
    useradd -r -u 1001 -g 1001 -m -c "${USER} account" -d /home/${USER} -s /bin/bash ${USER}

USER ${USER}
WORKDIR /home/${USER}

COPY --from=build --chown=${USER}:{USER} /home/${USER}/.local/bin/godyl /home/${USER}/.local/bin/godyl

ENV PATH=$PATH:/home/${USER}/.local/bin
ENV PATH=$PATH:/root/.local/bin

# Timezone
ENV TZ=Europe/Zurich
