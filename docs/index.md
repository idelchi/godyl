---
layout: default
title: Home
nav_order: 1
---

<p align="center">
  <img alt="godyl logo" src="assets/images/go.png" height="150" />
  <h1 align="center">godyl</h1>
  <p align="center">Asset downloader</p>
</p>

# Godyl

`godyl` helps with batch-downloading and installing statically compiled binaries from:

- GitHub releases
- URLs
- Go projects

As an alternative to above, custom commands can be used as well.

`godyl` will infer the platform and architecture from the system it is running on, and will attempt to download the appropriate binary.

This uses simple heuristics to select the correct binary to download, and will not work for all projects.

However, most properties can be overridden, with `hints` and `skip` used to help `godyl` make the correct decision.

{: .note }

> **Note:** Set up a GitHub API token to avoid rate limiting when using `github` as a source type.
> See [configuration](configuration.html) for more information, or simply `export GODYL_GITHUB_TOKEN=<token>`

{: .note }

> **Note:** Tested on:
>
> **Linux**: `amd64`, `arm64`
>
> **Windows**: `amd64`
>
> **MacOS**: `arm64`
>
> for tools listed in [tools.yml](https://github.com/idelchi/godyl/blob/main/tools.yml)

Tool is inspired by [task](https://github.com/go-task/task), [dra](https://github.com/devmatteini/dra) and [ansible](https://github.com/ansible/ansible)

## Quick Links

- [Installation](installation.html) - Install godyl
- [Usage](usage.html) - Learn how to use godyl
- [Configuration](configuration.html) - Configure godyl
- [Tools Definition](tools.html) - Define tools to download
- [Advanced](advanced.html) - Advanced features and concepts

## Features

- **Multiple Sources**: Download from GitHub releases, direct URLs, or Go packages
- **Platform Detection**: Automatically detects OS and architecture
- **Smart Selection**: Uses heuristics to select the appropriate asset
- **Customizable**: Override defaults with hints, patterns, and custom commands
- **Caching**: Caches downloads to avoid unnecessary network traffic
- **Batch Processing**: Install multiple tools at once from YAML configuration
