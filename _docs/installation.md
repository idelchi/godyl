---
layout: default
title: Installation
nav_order: 2
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

### Install all tools defined in the embedded tools file

```sh
curl -sSL https://raw.githubusercontent.com/idelchi/godyl/refs/heads/dev/scripts/tools.sh | sh -s -- -o ~/.local/bin
```

### Install Kubernetes-related tools

```sh
curl -sSL https://raw.githubusercontent.com/idelchi/godyl/refs/heads/dev/scripts/k8s.sh | sh -s -- -o ~/.local/bin
```

### Extract specific tools

```sh
curl -sSL https://raw.githubusercontent.com/idelchi/godyl/refs/heads/dev/scripts/extract.sh | sh -s -- -o ~/.local/bin idelchi/gogen idelchi/tcisd
```

## Verifying Installation

After installation, verify that `godyl` is working correctly:

```sh
godyl --version
```

This should display the installed version of `godyl`.
