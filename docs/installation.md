---
layout: default
title: Installation
---

# Installation

There are several ways to install `godyl`:

## From Source

If you have Go installed (1.24+), you can install directly from source:

```sh
go install github.com/idelchi/godyl@latest
```

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

You can also run `godyl` using Docker. The following command will run the latest version of `godyl`:

```sh
docker run --rm -v $(pwd):~/.local/bin docker.io/idelchi/godyl:latest install --output=/data
```
