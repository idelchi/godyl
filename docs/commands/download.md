---
layout: default
title: Download Command
parent: Commands
---

# Download Command

The `download` command allows you to download and install individual tools without requiring a configuration file.

## Syntax

```sh
godyl download [tool|URL]... [flags]
```

## Description

The `download` command provides a quick way to download tools directly from GitHub repositories or URLs. You can specify tools in the format `owner/repo` for GitHub repositories, or provide a full URL.

When using the `download` command, the tool will be downloaded and extracted directly into the output directory.

## Flags

| Flag                    | Environment Variable       | Default  | Description                             |
| ----------------------- | -------------------------- | -------- | --------------------------------------- |
| `--output`, `-o`        | `GODYL_TOOL_OUTPUT`        | `./bin`  | Output path for the downloaded tools    |
| `--source`              | `GODYL_TOOL_SOURCE`        | `github` | Source from which to install the tools  |
| `--os`                  | `GODYL_TOOL_OS`            | `""`     | Operating system to use for downloading |
| `--arch`                | `GODYL_TOOL_ARCH`          | `""`     | Architecture to use for downloading     |
| `--no-verify-ssl`, `-k` | `GODYL_TOOL_NO_VERIFY_SSL` | `false`  | Skip SSL verification                   |
| `--hint`                | `GODYL_TOOL_HINT`          | `[""]`   | Add hint patterns with weight 1         |
| `--version`, `-v`       | `GODYL_TOOL_VERSION`       | `""`     | Version to download                     |

## Examples

### Download a tool from GitHub

```sh
godyl download idelchi/godyl
```

### Download multiple tools

```sh
godyl download idelchi/godyl idelchi/tcisd
```

### Download a specific version

```sh
godyl download idelchi/godyl --version v0.1.0
```

### Download for a different platform

```sh
godyl download idelchi/godyl --os linux --arch arm64
```

### Download from a direct URL

```sh
godyl download "https://example.com/download/tool_{{ .OS }}_{{ .ARCH }}.tar.gz"
```

### Specify an output directory

```sh
godyl download idelchi/godyl --output ~/.local/bin
```

## How It Works

1. The `download` command first determines the source of the tool (GitHub repository or URL).
2. It then fetches information about available releases (for GitHub sources).
3. Using platform information and hints, it selects the appropriate release asset to download.
4. The asset is downloaded and extracted to the specified output directory.
5. Any executable files are identified and made available in the output directory.

## Differences from the Install Command

The `download` command differs from the `install` command in several ways:

- It operates on individual tools rather than a collection of tools defined in a YAML file.
- It always uses the `extract` mode, extracting files directly to the output directory.
- It doesn't support features like tagging, customized aliases, or pre/post commands.

For more complex installation scenarios, consider using the [`install`](install.html) command with a YAML configuration file.

## Related Topics

- [Install Command](install.html)
- [URL Templates](../templates.html#url-templates)
- [Platform Detection](../advanced-features.html#platform-inference)
