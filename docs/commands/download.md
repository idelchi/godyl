---
layout: default
title: Download Command
---

# Download Command

The `download` command allows you to download and unarchive individual tools without requiring a configuration file.

## Syntax

```sh
godyl [flags] download [tool|URL]... [flags]
```

## Aliases

- `dl`
- `unpack`
- `extract`
- `x`

## Description

The `download` command provides a quick way to download tools. You can specify tools in the format `owner/repo` for GitHub/GitLab repositories, or provide a full URL.

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

### Download a specific version

```sh
godyl download idelchi/godyl --version v0.1.0
```

### Download multiple tools

```sh
godyl download idelchi/godyl idelchi/tcisd
```

### Download from a direct URL

{% raw  %}

```sh
godyl download "https://github.com/idelchi/go-next-tag/releases/download/v0.0.1/go-next-tag_{{ .OS }}_{{ .ARCH }}.tar.gz"
```

{% endraw %}

## Related Topics

- [Global Flags]({{ site.baseurl }}/commands/commands#global-flags)
