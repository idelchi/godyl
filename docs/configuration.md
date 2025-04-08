---
layout: default
title: Configuration
---

# Configuration

Godyl can be configured in several ways (in order of priority):

1. Command-line flags
2. Environment variables
3. `.env` file(s)
4. `defaults.yml` file
5. Default embedded configuration

## Tool-specific Flags

The following flags are available for tool-related commands (`install` and `download`):

| Flag | Environment Variable | Default | Description |
|------|---------------------|---------|-------------|
| `--output`, `-o` | `GODYL_TOOL_OUTPUT` | `./bin` | Output path for the downloaded tools |
| `--tags`, `-t` | `GODYL_TOOL_TAGS` | `["!native"]` | Tags to filter tools by. Use `!` to exclude |
| `--source` | `GODYL_TOOL_SOURCE` | `github` | Source from which to install the tools |
| `--strategy` | `GODYL_TOOL_STRATEGY` | `none` | Strategy to use for updating tools |
| `--os` | `GODYL_TOOL_OS` | `""` | Operating system to use for downloading |
| `--arch` | `GODYL_TOOL_ARCH` | `""` | Architecture to use for downloading |
| `--parallel`, `-j` | `GODYL_TOOL_PARALLEL` | `0` | Number of parallel downloads (0 is unlimited) |
| `--no-verify-ssl`, `-k` | `GODYL_TOOL_NO_VERIFY_SSL` | `false` | Skip SSL verification |
| `--hint` | `GODYL_TOOL_HINT` | `[""]` | Add hint patterns with weight 1 |
| `--version`, `-v` | `GODYL_TOOL_VERSION` | `""` | Version to download (only used for `download`) |
| `--show`, `-s` | `GODYL_TOOL_SHOW` | `false` | Show the configuration and exit |

## Configuration Settings

Settings can be set in the following ways:

1. As a field in the `tools.yml` definition
   ```yaml
   output: ~/.local/bin
   ```

2. As a flag to the tool
   ```sh
   godyl --output ~/.local/bin
   ```

3. As an environment variable
   ```sh
   GODYL_TOOL_OUTPUT=~/.local/bin godyl
   ```

4. In an `.env` file
   ```
   GODYL_TOOL_OUTPUT=~/.local/bin
   ```

5. By setting the value in a `defaults.yml` file
   ```yaml
   output: ~/.local/bin
   ```

## Tools Configuration

Tools can be defined in a YAML file (typically `tools.yml`). You can use a simple form or a full form for tool definitions.

### Simple Form

```yaml
- idelchi/godyl
```

This is the simplest form to download the latest release of `godyl` from the GitHub repository `idelchi/godyl`.

If the path is a URL, it will be considered as a `source.url` type. Otherwise, it will be assumed to be a `source.github` type in the form `owner/repo`.

### Full Form

For more complex configurations, you can use the full form:

```yaml
name: godyl
description: Asset downloader
version:
  version: v0.1.0
path: https://github.com/idelchi/godyl/releases/download/v0.1.0/godyl_linux_amd64.tar.gz
output: ~/.local/bin
exe:
  name: godyl
  patterns:
    - "{{ .Exe }}{{ .EXTENSION }}$"
platform:
  os: linux
  architecture:
    type: amd64
aliases:
  - gd
source:
  type: github
  github:
    repo: godyl
    owner: idelchi
    token: ${GITHUB_TOKEN}
tags:
  - cli
  - downloader
strategy: none
```

## Templating

Many fields in the configuration support templating with variables like:

- `{{ .Name }}` - The name of the tool
- `{{ .Output }}` - The output path
- `{{ .Exe }}` - The executable name
- `{{ .OS }}` - The operating system
- `{{ .ARCH }}` - The architecture
- `{{ .EXTENSION }}` - The file extension for the platform

For example, to set a path that adapts to the current platform:

```yaml
path: https://example.com/download/{{ .Name }}_{{ .OS }}_{{ .ARCH }}.tar.gz
```