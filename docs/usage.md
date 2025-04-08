---
layout: default
title: Usage
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

## Download Command

The `download` command allows you to download a single tool without a configuration file:

```sh
godyl download idelchi/godyl --output ./bin
```

You can also download multiple tools:

```sh
godyl download idelchi/tcisd idelchi/gogen idelchi/wslint
```