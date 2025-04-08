---
layout: default
---

# Godyl

Asset downloader for GitHub releases, URLs, and Go projects.

## What is Godyl?

Godyl is a powerful tool designed to simplify the process of downloading and installing statically compiled binaries from various sources. Whether you need a single tool or want to set up an entire development environment, Godyl has you covered.

### Key Features

- **Multiple Sources**: Download assets from GitHub releases, GitLab releases, URLs, or Go projects
- **Intelligent Platform Detection**: Automatically selects the right binary for your platform
- **Batch Installation**: Install multiple tools from a single YAML configuration file
- **Caching**: Reduces bandwidth usage and speeds up repeated installations
- **Templating**: Customize paths and behaviors using Go templates
- **Cross-Platform**: Works on Linux, Windows, and macOS

## Getting Started

### Quick Installation

Install Godyl in seconds:

```sh
curl -sSL https://raw.githubusercontent.com/idelchi/godyl/refs/heads/dev/install.sh | sh -s -- -d ~/.local/bin
```

### Basic Usage

Download a single tool:

```sh
godyl download idelchi/godyl --output ~/.local/bin
```

Create a `tools.yml` file to define multiple tools:

```yaml
- name: kubectl
  source:
    type: url
    path: https://dl.k8s.io/release/v1.27.3/bin/{{ .OS }}/{{ .ARCH }}/kubectl{{ .EXTENSION }}
  mode: extract

- name: helm
  source:
    type: github
    github:
      owner: helm
      repo: helm
```

Then install them all at once:

```sh
godyl install tools.yml --output ~/.local/bin
```

## Documentation

Explore the comprehensive documentation to learn more about Godyl's capabilities:

### Getting Started

- [Installation](installation)
- [Usage](usage)

### Configuration

- [Configuration Basics](configuration)
- [Tools Format](tools-format)
- [Default Configuration](defaults)
- [Advanced Features](advanced-features)

### Command Reference

- [Commands Overview](commands/)
- [Install Command](commands/install)
- [Download Command](commands/download)
- [Dump Command](commands/dump)
- [Update Command](commands/update)
- [Cache Command](commands/cache)

### Examples and Templates

- [Simple Examples](examples)
- [Advanced Examples](advanced-examples)
- [Template Reference](templates)

## Why Use Godyl?

Godyl solves common challenges in managing development tools:

- **Consistency**: Ensure everyone on your team uses the same tool versions
- **Convenience**: Set up a complete development environment with a single command
- **Flexibility**: Works with a wide variety of tools and source types
- **Platform Independence**: The same configuration works across different operating systems

## Use Cases

### Setting Up a Development Environment

Create a `dev-tools.yml` file with all the tools you need for development:

```yaml
- idelchi/godyl
- helm/helm
- derailed/k9s
- kubernetes/kubectl
- hashicorp/terraform
```

Install them all with a single command:

```sh
godyl install dev-tools.yml --output ~/.local/bin
```

### Creating a Project-specific Toolchain

Include a `tools.yml` file in your project repository to ensure everyone uses the same tool versions:

```yaml
- name: gopls
  version: v0.11.0
  tags: [go, lsp]

- name: golangci-lint
  version: v1.52.2
  tags: [go, linter]
```

### CI/CD Setup

Use Godyl in your CI/CD pipelines to ensure consistent tool availability:

```yaml
# In your GitHub Actions workflow
steps:
  - name: Install Tools
    run: |
      curl -sSL https://raw.githubusercontent.com/idelchi/godyl/refs/heads/dev/install.sh | sh -s -- -d /usr/local/bin
      godyl install tools.yml --output /usr/local/bin
```

## External Links

- [GitHub Repository](https://github.com/idelchi/godyl)
- [Go Package Documentation](https://pkg.go.dev/github.com/idelchi/godyl)
