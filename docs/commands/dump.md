---
layout: default
title: Dump Command
parent: Commands
---

# Dump Command

The `dump` command displays various configuration settings and information about Godyl.

## Syntax

```sh
godyl dump [config|defaults|env|platform|tools|cache] [flags]
```

## Description

The `dump` command provides a way to inspect Godyl's configuration, available tools, and system information. This can be helpful for debugging, understanding the current setup, or creating custom tool configurations.

## Subcommands

| Subcommand | Description                                               |
| ---------- | --------------------------------------------------------- |
| `config`   | Display the current configuration settings                |
| `defaults` | Display the default configuration settings                |
| `env`      | Display environment variables that affect the application |
| `platform` | Display information about the current platform            |
| `tools`    | Display information about available tools                 |
| `cache`    | Display information about the cache                       |

## Flags for `dump tools`

| Flag           | Environment Variable    | Default | Description                |
| -------------- | ----------------------- | ------- | -------------------------- |
| `--full`, `-f` | `GODYL_DUMP_TOOLS_FULL` | `false` | Show full tool information |

## Examples

### Display the current configuration

```sh
godyl dump config
```

Output will contain the active configuration settings derived from all sources (command-line flags, environment variables, defaults, etc.).

### Display the default configuration

```sh
godyl dump defaults
```

Output will show the default configuration that's embedded in the Godyl binary.

### Display environment variables

```sh
godyl dump env
```

Output will list all environment variables that affect Godyl's behavior.

### Display platform information

```sh
godyl dump platform
```

Output will show details about your current platform, including OS, architecture, and other system information.

### Display available tools

```sh
godyl dump tools
```

Output will list the tools that are embedded in the Godyl binary. This can be used as a starting point for creating custom tool configurations.

### Display full tool information

```sh
godyl dump tools --full
```

Output will show detailed information about each tool, including all available configuration options.

### Display cache information

```sh
godyl dump cache
```

Output will show details about the cache, including its location, size, and contents.

## Practical Uses

### Creating a Custom Tools Configuration

You can use the `dump tools` command to create a starting point for your own tools configuration:

```sh
godyl dump tools > my-tools.yml
```

This creates a YAML file containing all the embedded tools, which you can then modify according to your needs.

### Installing Embedded Tools

You can use the `dump tools` command in combination with the `install` command to install all embedded tools:

```sh
godyl dump tools | godyl install - --output ~/.local/bin
```

### Filtering Tools with YQ

You can combine the `dump tools` command with tools like `yq` to filter tools:

```sh
godyl dump tools | yq --yaml-output '[.[] | try (select(.tags != null and (.tags[] == "docker")))]' > docker-tools.yml
```

This creates a new YAML file containing only tools tagged with "docker".

## Related Topics

- [Global Flags](index#global-flags)
