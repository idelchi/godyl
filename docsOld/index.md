---
layout: default
---

# Godyl

Asset downloader for GitHub releases, GitLab release, URLs, and Go projects.

## What is Godyl?

`godyl` helps with batch-downloading and installing statically compiled binaries from:

- GitHub releases
- GitLab releases
- URLs
- Go projects

As an alternative to above, custom commands can be used as well.

`godyl` will infer the platform and architecture from the system it is running on, and will attempt to download the appropriate binary.

This uses simple heuristics to select the correct binary to download, and will not work for all projects.

However, most properties can be overridden, with `hints` and `skip` used to help the tool make the correct decision.

Godyl has been tested on:

- **Linux**: `amd64`, `arm64`
- **Windows**: `amd64`
- **MacOS**: `arm64`

for the tools listed in the default [tools.yml](https://github.com/idelchi/godyl/blob/main/tools.yml) file.

> **Note**: To avoid GitHub API rate limiting when using `github` as a source type, set up a GitHub API token by either using the `--github-token` flag or setting the `GODYL_GITHUB_TOKEN` environment variable.

Tool is inspired by [task](https://github.com/go-task/task), [dra](https://github.com/devmatteini/dra) and [ansible](https://github.com/ansible/ansible)

## Getting Started

### Quick Installation

```sh
curl -sSL https://raw.githubusercontent.com/idelchi/godyl/refs/heads/dev/install.sh | sh -s -- -d ~/.local/bin
```

### Basic Usage

Download a single tool:

```sh
godyl download idelchi/godyl
```

Create a `tools.yml` file to define multiple tools:

```yaml
- name: syncthing/syncthing
  tags:
    - sync

- name: helm/helm
  path: https://get.helm.sh/helm-{{ .Version }}-{{ .OS }}-{{ .ARCH }}.tar.gz
  tags:
    - kubernetes
```

Then install them all at once:

```sh
godyl install tools.yml
```

## Documentation

### Getting Started

- [Installation](installation#content-start)
- [Commands](commands/index#content-start)

### Configuration

- [Configuration Basics](configuration/configuration#content-start)
- [Default Configuration](configuration/defaults#content-start)
- [Tools Format](configuration/tools#content-start)
- [Templates](configuration/templates#content-start)

### Command Reference

- [Commands Overview](commands/index#content-start)
- [Install Command](commands/install#content-start)
- [Download Command](commands/download#content-start)
- [Dump Command](commands/dump#content-start)
- [Update Command](commands/update#content-start)
- [Cache Command](commands/cache#content-start)

## Use Cases

### Setting Up a Development Environment

Create a `dev-tools.yml` file with all the tools you need for development and install them all with a single command:

```sh
godyl install dev-tools.yml --output ~/.local/bin
```

### Creating Project-specific Toolchains

Include a `tools.yml` file in your project repository to ensure everyone uses the same tool versions:

```yaml
- name: google/go-jsonnet
  version: v0.18.0
  tags:
    - json

- name: golangci/golangci-lint
  version: v1.52.2
  tags:
    - go
    - linter
```

## External Links

- [GitHub Repository](https://github.com/idelchi/godyl)
- [Go Package Documentation](https://pkg.go.dev/github.com/idelchi/godyl)
