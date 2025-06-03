---
layout: default
title: Dump
parent: Commands
nav_order: 4
---

# Dump Command

The `dump` command displays various configuration settings and information about `godyl`.

## Syntax

```sh
godyl [flags] dump [defaults|env|platform|tools] [flags]
```

## Description

The `dump` command provides a way to inspect `godyl`'s configuration, available tools, and system information. This can be helpful for debugging, understanding the current setup, or creating custom tool configurations.

## Subcommands

| Subcommand           | Description                                               |
| :------------------- | :-------------------------------------------------------- |
| `defaults [name...]` | Display the default configuration settings                |
| `env`                | Display environment variables that affect the application |
| `platform`           | Display information about the current platform            |
| `tools`              | Display information about available tools                 |
| `cache [name]`       | Display information about the cache                       |
| `config [name]`      | Display information about the config                      |

## Flags for `dump tools`

| Flag              | Environment Variable        | Default | Description                |
| :---------------- | :-------------------------- | :------ | :------------------------- |
| `--full`, `-f`    | `GODYL_DUMP_TOOLS_FULL`     | `false` | Show full tool information |
| `--tags`, `-f`    | `GODYL_DUMP_TOOLS_TAGS`     | `false` | Filter by tags             |
| `--embedded, `-e` | `GODYL_DUMP_TOOLS_EMBEDDED` | `true`  | Show only embedded tools   |

## Flags for `dump config`

| Flag           | Environment Variable    | Default | Description                |
| :------------- | :---------------------- | :------ | :------------------------- |
| `--full`, `-f` | `GODYL_DUMP_TOOLS_FULL` | `false` | Show full tool information |

## Examples

### Display the default configuration embedded in the binary

Output the whole default map:

```sh
godyl dump defaults
```

Select specific keys:

```sh
godyl dump defaults linux default
```

### Display environment variables

```sh
godyl dump env
```

Output will list all environment variables that affect `godyl`'s behavior.

### Display platform information

```sh
godyl dump platform
```

Output will show details about your current platform, including OS, architecture, and other system information.

### Display available tools

```sh
godyl dump tools
```

Output will list the tools that are embedded in the `godyl` binary. This can be used as a starting point for creating custom tool configurations.

### Display full tool information

```sh
godyl dump tools --full
```

Output will show detailed information about each tool, including all available configuration options.

### Display cache information

```sh
godyl dump cache
```

### Display cache information for a specific item

```sh
godyl dump cache idelchi/godyl
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

or only specific tags:

```sh
godyl dump tools --tags docker | godyl install - --output ~/.local/bin
```
