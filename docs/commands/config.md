---
layout: default
title: Config
parent: Commands
nav_order: 7
---

# Config Command

The `config` command allows interaction with `godyl`'s config file.

## Syntax

```sh
godyl config [path|set] [flags]
```

## Aliases

- `cfg`

## Subcommands

| Subcommand          | Description                       |
| :------------------ | :-------------------------------- |
| `path`              | Print the path to the config file |
| `set <key> <value>` | Set a value in the config file    |

> **Note**: The `set` command will lead to loss of order and comments in the config file.

## Examples

### Store tokens in the config file

```sh
godyl config set github-token secret-token
```
