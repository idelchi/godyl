---
layout: default
title: Install Command
---

# Install Command

The `install` command allows you to install tools defined in YAML configuration files.

## Syntax

```sh
godyl [flags] install [tools.yml...|STDIN] [flags]
```

## Aliases

- `get`
- `i`

## Description

The `install` command reads tool definitions from one or more YAML files (or from standard input) and installs them according to the specified configuration. If no file is specified, it defaults to using `tools.yml` in the current directory.

## Flags

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

## Examples

### Install tools from tools.yml in the current directory

```sh
godyl install
```

### Install tools from multiple files

```sh
godyl install tools1.yml tools2.yml
```

### Install tools from stdin

```sh
cat tools.yml | godyl install -
```

### Filter tools by tag

```sh
godyl install tools.yml --tags cli,downloader
```

### Exclude tools with a specific tag

```sh
godyl install tools.yml --tags '!experimental'
```

### Force reinstallation of all tools

```sh
godyl install tools.yml --strategy force
```

### Install tools for a different platform

```sh
godyl install tools.yml --os linux --arch arm64
```

## Related Topics

- [Global Flags]({{ site.baseurl }}/commands/commands#global-flags)
- [Tool Configuration Format]({{ site.baseurl }}/configuration/tools#content-start)
