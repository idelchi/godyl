---
layout: default
title: Home
nav_order: 1
description: "Asset downloader for GitHub releases, GitLab releases, URLs, and Go projects."
permalink: /
---

{: .text-center }
![Godyl Logo]({{ site.baseurl }}/assets/images/godyl.png){: style="height: 320px; width: auto;"}

# Godyl

Asset downloader for GitHub releases, GitLab releases, URLs, and Go projects.
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

`godyl` uses deterministic heuristics to select the correct binary to download, matching the current platform and architecture.

Most properties can be overridden, with `hints` and `skip` used to help the tool make the correct decision.

For releases with equally ranked deterministic matches, an opt-in AI fallback
can select from the tied candidates. It is not called for successful or
unmatched results, and its answer must match one supplied candidate exactly.
See [AI fallback]({{ site.baseurl }}/commands/index#ai-fallback)
for Ollama and OpenAI configuration.

`godyl` has been tested on:

- **Linux**: `amd64`, `arm64`
- **Windows**: `amd64`
- **MacOS**: `arm64`

for the tools listed in the default [tools.yml](https://github.com/idelchi/godyl/blob/main/tools.yml) file.

> **Note**: You'll have a very short journey with this tool without a GitHub API token. To avoid rate limiting when using `github` as a source type, set up an API token and use it with the `--github-token` flag or the `GODYL_GITHUB_TOKEN` environment variable. See [Authentication]({{ site.baseurl }}/commands/index#authentication) for more details. By not using a token, `godyl` will attempt to use the unauthenticated web API firstly, which might lead to rate limiting / blocking if you make too many requests in a short time. As such, the parallelism is set to `1` by default when no token is provided.

Tool is inspired by [task](https://github.com/go-task/task), [dra](https://github.com/devmatteini/dra) and [ansible](https://github.com/ansible/ansible)

## Getting Started

### Quick Installation

```sh
curl -sSL https://raw.githubusercontent.com/idelchi/godyl/refs/heads/main/install.sh | sh -s -- -d ~/.local/bin
```

### Basic Usage

{% raw %}

Export the `GODYL_GITHUB_TOKEN` environment variable with your GitHub API token:

```sh
export GODYL_GITHUB_TOKEN=<your_github_token>
```

Download (and extract) a single tool:

```sh
godyl download idelchi/envprof
```

Create a `tools.yml` file to define multiple tools:

```yaml
- name: syncthing/syncthing
  tags:
    - sync

- name: helm/helm
  url: https://get.helm.sh/helm-{{ .Version }}-{{ .OS }}-{{ .ARCH }}.tar.gz
  checksum:
    type: file
    value: "{{ .URL }}.sha256sum"
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

To resolve ambiguous deterministic asset matches, enable the optional AI
fallback for the invocation:

```sh
godyl --ai install tools.yml
```

To inspect asset matching failures and receive verified hint suggestions without
installing anything:

```sh
godyl install --suggest tools.yml
```

{% endraw %}

For a sample, see [tools.yml](https://github.com/idelchi/godyl/blob/main/tools.yml) or run `godyl dump tools -e > tools.yml` to inspect the default configuration.
