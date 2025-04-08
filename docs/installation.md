---
layout: default
title: Installation
---

# Installation

There are several ways to install `godyl`:

## From Source

If you have Go installed (1.20+), you can install directly from source:

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

## Without Installation

If you prefer not to install `godyl`, you can use the convenience scripts provided:

```sh
# Install all tools defined in the embedded tools file
curl -sSL https://raw.githubusercontent.com/idelchi/godyl/refs/heads/dev/scripts/tools.sh | sh -s -- -o ~/.local/bin

# Install Kubernetes-related tools
curl -sSL https://raw.githubusercontent.com/idelchi/godyl/refs/heads/dev/scripts/k8s.sh | sh -s -- -o ~/.local/bin

# Extract specific tools
curl -sSL https://raw.githubusercontent.com/idelchi/godyl/refs/heads/dev/scripts/extract.sh | sh -s -- -o ~/.local/bin idelchi/gogen idelchi/tcisd
```

## Platform Support

Godyl has been tested on:

- **Linux**: `amd64`, `arm64`
- **Windows**: `amd64`
- **MacOS**: `arm64`

For the tools listed in the default `tools.yml` file.

> **Note**: To avoid GitHub API rate limiting when using `github` as a source type, set up a GitHub API token by either using the `--github-token` flag or setting the `GODYL_GITHUB_TOKEN` environment variable.
