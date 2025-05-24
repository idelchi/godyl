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

---

## What is Godyl?

`godyl` aims to help with batch-downloading from:

- GitHub releases
- GitLab releases
- URLs
- Go projects

Furthermore, custom commands can be used.

`godyl` uses simple heuristics to select the correct binary to download, matching the current platform and architecture.

Most properties can be overridden, with `hints` and `skip` used to help the tool make the correct decision.

`godyl` has been tested on:

- **Linux**: `amd64`, `arm64`
- **Windows**: `amd64`
- **MacOS**: `arm64`

for the tools listed in the default [tools.yml](https://github.com/idelchi/godyl/blob/dev/tools.yml) file.

> **Note**: You'll have a very short journey with this tool without a GitHub API token. To avoid rate limiting when using `github` as a source type, set up an API token and use it with the `--github-token` flag or the `GODYL_GITHUB_TOKEN` environment variable.

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

To periodically update/sync, run above command with the `--strategy` flag:

```sh
godyl install tools.yml --strategy=sync
```

to bring down the latest version, if the current one is out of date.

{% endraw %}

For a sample, see [tools.yml](https://github.com/idelchi/godyl/blob/main/tools.yml) or run `godyl dump tools > tools.yml` to inspect the default configuration.
