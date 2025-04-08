---
layout: default
title: Command Reference
---

# Command Reference

This page provides detailed information about all Godyl commands and their options.

## Global Flags

These flags are available for all commands:

| Flag                 | Environment Variable     | Default                 | Description                                  |
| -------------------- | ------------------------ | ----------------------- | -------------------------------------------- |
| `--help`, `-h`       | `GODYL_HELP`             | `false`                 | Show help message and exit                   |
| `--version`          | `GODYL_VERSION`          | `false`                 | Show version information and exit            |
| `--dry`              | `GODYL_DRY`              | `false`                 | Run without making any changes (dry run)     |
| `--log`              | `GODYL_LOG`              | `info`                  | Log level (debug, info, warn, error, silent) |
| `--env-file`, `-e`   | `GODYL_ENV_FILE`         | `[".env"]`              | Path to `.env` file(s)                       |
| `--defaults`, `-d`   | `GODYL_DEFAULTS`         | `defaults.yml`          | Path to defaults file                        |
| `--github-token`     | `GODYL_GITHUB_TOKEN`     | `${GODYL_GITHUB_TOKEN}` | GitHub token for authentication              |
| `--gitlab-token`     | `GODYL_GITLAB_TOKEN`     | `${GODYL_GITLAB_TOKEN}` | GitLab token for authentication              |
| `--url-token`        | `GODYL_URL_TOKEN`        | `${GODYL_URL_TOKEN}`    | URL token for authentication                 |
| `--url-token-header` | `GODYL_URL_TOKEN_HEADER` | `Authorization`         | URL header for authentication                |
| `--cache-dir`, `-c`  | `GODYL_CACHE_DIR`        | `${XDG_CACHE_HOME}`     | Path to cache directory                      |
| `--cache-type`       | `GODYL_CACHE_TYPE`       | `file`                  | Type of cache (file, sqlite)                 |

## Install Command

Install tools defined in YAML configuration files.

```sh
godyl install [[tools.yml]...|STDIN] [flags]
```

### Flags

| Flag                    | Environment Variable       | Default       | Description                                   |
| ----------------------- | -------------------------- | ------------- | --------------------------------------------- |
| `--output`, `-o`        | `GODYL_TOOL_OUTPUT`        | `./bin`       | Output path for the downloaded tools          |
| `--tags`, `-t`          | `GODYL_TOOL_TAGS`          | `["!native"]` | Tags to filter tools by. Use `!` to exclude   |
| `--source`              | `GODYL_TOOL_SOURCE`        | `github`      | Source from which to install the tools        |
| `--strategy`            | `GODYL_TOOL_STRATEGY`      | `none`        | Strategy to use for updating tools            |
| `--os`                  | `GODYL_TOOL_OS`            | `""`          | Operating system to use for downloading       |
| `--arch`                | `GODYL_TOOL_ARCH`          | `""`          | Architecture to use for downloading           |
| `--parallel`, `-j`      | `GODYL_TOOL_PARALLEL`      | `0`           | Number of parallel downloads (0 is unlimited) |
| `--no-verify-ssl`, `-k` | `GODYL_TOOL_NO_VERIFY_SSL` | `false`       | Skip SSL verification                         |
| `--hint`                | `GODYL_TOOL_HINT`          | `[""]`        | Add hint patterns with weight 1               |
| `--show`, `-s`          | `GODYL_TOOL_SHOW`          | `false`       | Show the configuration and exit               |

### Examples

```sh
# Install tools from tools.yml in the current directory
godyl install

# Install tools from a specific file
godyl install my-tools.yml

# Install tools from multiple files
godyl install tools1.yml tools2.yml

# Install tools from stdin
cat tools.yml | godyl install -

# Install tools with custom output directory
godyl install tools.yml --output ~/.local/bin

# Filter tools by tag
godyl install tools.yml --tags cli,downloader

# Exclude tools with a specific tag
godyl install tools.yml --tags '!experimental'

# Force reinstallation of all tools
godyl install tools.yml --strategy force

# Install tools for a different platform
godyl install tools.yml --os linux --arch arm64

# Limit parallel downloads
godyl install tools.yml --parallel 4
```

## Download Command

Download and install individual tools without a configuration file.

```sh
godyl download [tool|URL]... [flags]
```

### Flags

| Flag                    | Environment Variable       | Default  | Description                             |
| ----------------------- | -------------------------- | -------- | --------------------------------------- |
| `--output`, `-o`        | `GODYL_TOOL_OUTPUT`        | `./bin`  | Output path for the downloaded tools    |
| `--source`              | `GODYL_TOOL_SOURCE`        | `github` | Source from which to install the tools  |
| `--os`                  | `GODYL_TOOL_OS`            | `""`     | Operating system to use for downloading |
| `--arch`                | `GODYL_TOOL_ARCH`          | `""`     | Architecture to use for downloading     |
| `--no-verify-ssl`, `-k` | `GODYL_TOOL_NO_VERIFY_SSL` | `false`  | Skip SSL verification                   |
| `--hint`                | `GODYL_TOOL_HINT`          | `[""]`   | Add hint patterns with weight 1         |
| `--version`, `-v`       | `GODYL_TOOL_VERSION`       | `""`     | Version to download                     |

### Examples

```sh
# Download a tool from GitHub
godyl download idelchi/godyl

# Download multiple tools
godyl download idelchi/godyl idelchi/tcisd

# Download a specific version
godyl download idelchi/godyl --version v0.1.0

# Download for a different platform
godyl download idelchi/godyl --os linux --arch arm64

# Download from a direct URL
godyl download "https://example.com/download/tool_{{ .OS }}_{{ .ARCH }}.tar.gz"

# Specify an output directory
godyl download idelchi/godyl --output ~/.local/bin
```

## Dump Command

Display various configuration settings and information.

```sh
godyl dump [config|defaults|env|platform|tools|cache] [flags]
```

### Subcommands

- `config` - Display the current configuration settings
- `defaults` - Display the default configuration settings
- `env` - Display environment variables that affect the application
- `platform` - Display information about the current platform
- `tools` - Display information about available tools
- `cache` - Display information about the cache

### Flags for `dump tools`

| Flag           | Environment Variable    | Default | Description                |
| -------------- | ----------------------- | ------- | -------------------------- |
| `--full`, `-f` | `GODYL_DUMP_TOOLS_FULL` | `false` | Show full tool information |

### Examples

```sh
# Display the current configuration
godyl dump config

# Display the default configuration
godyl dump defaults

# Display environment variables
godyl dump env

# Display platform information
godyl dump platform

# Display information about available tools
godyl dump tools

# Display full information about available tools
godyl dump tools --full

# Display cache information
godyl dump cache

# Dump tools and pipe to install
godyl dump tools | godyl install -
```

## Update Command

Update the godyl application to the latest version.

```sh
godyl update [flags]
```

### Flags

| Flag                    | Environment Variable         | Default | Description           |
| ----------------------- | ---------------------------- | ------- | --------------------- |
| `--no-verify-ssl`, `-k` | `GODYL_UPDATE_NO_VERIFY_SSL` | `false` | Skip SSL verification |
| `--version`, `-v`       | `GODYL_UPDATE_VERSION`       | `""`    | Version to download   |
| `--pre`                 | `GODYL_UPDATE_PRE`           | `false` | Include pre-releases  |

### Examples

```sh
# Update to the latest version
godyl update

# Update to a specific version
godyl update --version v0.1.0

# Include pre-releases when updating
godyl update --pre
```

## Cache Command

Manage the godyl cache.

```sh
godyl cache [flags]
```

### Flags

| Flag             | Environment Variable | Default | Description      |
| ---------------- | -------------------- | ------- | ---------------- |
| `--delete`, `-d` | `GODYL_CACHE_DELETE` | `false` | Delete the cache |

### Examples

```sh
# Display cache information
godyl cache

# Delete the cache
godyl cache --delete
```
