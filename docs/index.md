---
layout: default
---

# Godyl

Asset downloader for GitHub releases, GitLab releases, URLs, and Go projects.

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

## External Links

- [GitHub Repository](https://github.com/idelchi/godyl)
- [Go Package Documentation](https://pkg.go.dev/github.com/idelchi/godyl)
