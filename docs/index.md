---
layout: default
title: Home
nav_order: 1
description: "Asset downloader for GitHub releases, GitLab release, URLs, and Go projects."
permalink: /
---

{: .text-center }
![Godyl Logo]({{ site.baseurl }}/assets/images/godyl.png){: style="height: 320px; width: auto;"}

# Godyl

Asset downloader for GitHub releases, GitLab release, URLs, and Go projects.
{: .fs-6 .fw-300 }

[Get started now](#getting-started){: .btn .btn-primary .fs-5 .mb-4 .mb-md-0 .me-2 } [View it on GitHub](https://github.com/idelchi/godyl){: .btn .fs-5 .mb-4 .mb-md-0 }

![Godyl in Action (Install)]({{ site.baseurl }}/assets/gifs/install.gif)
![Godyl in Action (Download)]({{ site.baseurl }}/assets/gifs/download.gif)

---

## What is Godyl?

`godyl` aims to help with batch-downloading and "installing" statically compiled binaries from:

- GitHub releases
- GitLab releases
- URLs
- Go projects

`godyl` uses simple heuristics to select the correct binary to download, matching the current platform and architecture.

Most properties can be overridden, with `hints` and `skip` used to help the tool make the correct decision.

`godyl` has been tested on:

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

{% raw %}

Download (and extract) a single tool:

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

{% endraw %}

## Use Cases

`godyl` can be used to set up the same set of tools on machines, and periodically running it to keep them up to date.

For a sample, see [tools.yml](https://github.com/idelchi/godyl/blob/main/tools.yml) or run `godyl dump tools > tools.yml` to inspect the default configuration.

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
