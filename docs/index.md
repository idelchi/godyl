---
layout: default
---

# Godyl

Asset downloader for GitHub releases, URLs, and Go projects.

## Documentation

### Getting Started
- [Installation](installation.md)
- [Usage](usage.md)

### Configuration
- [Configuration Basics](configuration.md)
- [Tools Format](tools-format.md)
- [Default Configuration](defaults.md)
- [Advanced Features](advanced-features.md)

### External Links
- [GitHub Repository](https://github.com/idelchi/godyl)
- [Go Package Documentation](https://pkg.go.dev/github.com/idelchi/godyl)

## Features

- Download from GitHub releases
- Download from URLs
- Download Go projects
- Execute custom commands for installation
- Infer platform and architecture from system
- Filter tools by tags
- Various installation strategies (none, upgrade, force)
- Create aliases for downloaded tools
- Cache downloaded artifacts
- Works on Linux, Windows, and MacOS

## What is Godyl?

Godyl is a tool that helps with batch-downloading and installing statically compiled binaries from GitHub releases, URLs, and Go projects. It uses simple heuristics to select the correct binary to download based on your platform and architecture.

Most properties can be overridden, with `hints` and `skip` options to help Godyl make the correct decisions when downloading tools.