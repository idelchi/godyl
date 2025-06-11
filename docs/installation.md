---
layout: default
title: Installation
nav_order: 2
---

# Installation

## From Installation Script

For a quick installation, you can use the provided installation script:

```sh
curl -sSL https://raw.githubusercontent.com/idelchi/godyl/refs/heads/dev/install.sh | sh -s -- -d ~/.local/bin
```

For information on available options, run:

```sh
curl -sSL https://raw.githubusercontent.com/idelchi/godyl/refs/heads/dev/install.sh | sh -s -- -h
```

## Docker

You can test out the tool using Docker. The following command will run the latest version and install the default set of tools:

```sh
export GITHUB_TOKEN=<your_github_token>
docker run -it --rm --name godyl --env GITHUB_TOKEN docker.io/idelchi/godyl:dev

# Inside the container, run:
godyl dump tools -e | godyl install - --output=~/.local/bin
```

## From Source

If you have Go installed (1.24+), you can install directly from source:

```sh
go install github.com/idelchi/godyl@latest
```
