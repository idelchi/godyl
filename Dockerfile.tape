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
    nano \
    && rm -rf /var/lib/apt/lists/*

ARG TASK_VERSION=v3.41.0
RUN wget -qO- https://github.com/go-task/task/releases/download/${TASK_VERSION}/task_linux_${TARGETARCH}.tar.gz | tar -xz -C /usr/local/bin

RUN --mount=type=secret,id=.env \
    curl -sSL https://raw.githubusercontent.com/idelchi/godyl/refs/heads/dev/install.sh | sh -s -- -d /opt -v v0.0.12-beta && \
    # Install GoReleaser
    /opt/godyl dump tools --tags goreleaser/goreleaser,cli/cli,idelchi/go-next-tag | /opt/godyl --env-file=/run/secrets/.env i - --output=/usr/local/bin && \
    rm -rf /tmp/godyl

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

ARG TARGETOS TARGETARCH

COPY . .
ARG GODYL_VERSION="unofficial & built by unknown"
RUN --mount=type=cache,target=${GOMODCACHE},uid=${UID},gid=${UID},id=gomod-${TARGETARCH} \
    --mount=type=cache,target=${GOCACHE},uid=${UID},gid=${UID},id=gocache-${TARGETARCH} \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 go build -ldflags="-s -w -X 'main.version=${GODYL_VERSION}'" -o bin/ .

ENV PATH=$PATH:/home/${USER}/.local/bin
ENV PATH=$PATH:/root/.local/bin
ENV XDG_RUNTIME_DIR=/tmp/${UID}
ENV XDG_CONFIG_HOME=/home/${USER}/.config
ENV XDG_CACHE_HOME=/home/${USER}/.cache

RUN mkdir -p /home/${USER}/.local/bin
RUN cp bin/godyl /home/${USER}/.local/bin

WORKDIR /home/${USER}

FROM ghcr.io/charmbracelet/vhs:latest

# Create User (Debian/Ubuntu)
ARG USER=user
ARG UID=1001
RUN groupadd -r -g ${UID} ${USER} && \
    useradd -r -u ${UID} -g ${UID} -m -c "${USER} account" -d /home/${USER} -s /bin/bash ${USER}

USER ${USER}
WORKDIR /home/${USER}

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

USER ${USER}

RUN mkdir -p /home/user/gifs

ENV PATH=$PATH:/home/${USER}/.local/bin
ENV PATH=$PATH:/root/.local/bin
ENV XDG_RUNTIME_DIR=/tmp/${UID}
ENV XDG_CONFIG_HOME=/home/${USER}/.config
ENV XDG_CACHE_HOME=/home/${USER}/.cache

# Timezone
ENV TZ=Europe/Zurich

ENTRYPOINT [  ]
CMD [ "/bin/bash" ]

COPY --from=build --chown=${USER}:{USER} /home/${USER}/.local/bin/godyl /home/${USER}/.local/bin/godyl
