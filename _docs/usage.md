---
layout: default
title: Usage
nav_order: 3
---

# Usage

The basic syntax for using `godyl` is:

```sh
godyl [flags] [command] [flags]
```

## Available Commands

- `install` - Install tools from YAML files
- `download` - Download and unpack individual tools
- `dump` - Display configuration information
- `update` - Update the godyl application
- `cache` - Manage the godyl cache

## Install Command

The `install` command allows you to install tools defined in YAML configuration files:

```sh
godyl install [[tools.yml]...|STDIN] --output ./bin
```

If no file is specified, `godyl` defaults to using `tools.yml` in the current directory.

If the argument is set to `-`, `godyl` will read from `stdin`.

### Examples

Install tools from the default `tools.yml` file:

```sh
godyl install --output ~/.local/bin
```

Install tools from a specific YAML file:

```sh
godyl install custom-tools.yml --output ~/.local/bin
```

Install tools from multiple YAML files:

```sh
godyl install tools1.yml tools2.yml --output ~/.local/bin
```

## Download Command

The `download` command allows you to download a single tool without a configuration file:

```sh
godyl download idelchi/godyl --output ./bin
```

You can also download multiple tools:

```sh
godyl download idelchi/tcisd idelchi/gogen idelchi/wslint
```

When using the `download` command, the tool will be unarchived directly into the output directory.

### Override Platform Settings

You can override the OS and architecture detection:

```sh
godyl download idelchi/godyl --os linux --arch amd64 --output ./bin
```

### Download from Direct URL

You can download tools from direct URLs using templates:

```sh
godyl download "https://github.com/idelchi/go-next-tag/releases/download/v0.0.1/go-next-tag_{{ .OS }}_{{ .ARCH }}.tar.gz" --output ./bin
```

## Dump Command

The `dump` command displays various configuration settings and information:

```sh
godyl dump [config|defaults|env|platform|tools]
```

### Subcommands

- `config` - Display the current configuration settings
- `defaults` - Display the default configuration settings
- `env` - Display environment variables that affect the application
- `platform` - Display information about the current platform
- `tools` - Display information about available tools
- `cache` - Display information about the cache

### Examples

Display information about the current platform:

```sh
godyl dump platform
```

Install all tools that were embedded when the application was built:

```sh
godyl dump tools | godyl install - --output ./bin
```

## Update Command

The `update` command updates the godyl application to the latest version:

```sh
godyl update [flags]
```

{: .note }

> On Windows, this will launch a background process to clean up the old version.

### Examples

Update to the latest version:

```sh
godyl update
```

Update to a specific version:

```sh
godyl update --version v0.1.0
```

## Cache Command

The `cache` command allows you to manage the godyl cache:

```sh
godyl cache [flags]
```

### Examples

Display information about the cache:

```sh
godyl cache
```

Delete the cache:

```sh
godyl cache --delete
```
