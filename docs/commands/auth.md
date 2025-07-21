---
layout: default
title: auth
parent: Commands
nav_order: 8
---

# Config Command

The `auth` command provides a convenient way of adding or removing tokens from either the configuration file or the system keyring.

## Syntax

```sh
godyl auth [store|remove] [flags]
```

## Subcommands

| Subcommand                           | Description                                |
| :----------------------------------- | :----------------------------------------- |
| `store [token]...`                   | Store tokens from the parsed configuration |
| `remove [token]...`, `rm [token]...` | Remove authentication tokens               |

## Examples

### Set all values from the `tokens.env` file

```sh
godyl --env-file=tokens.env auth store
```

### Set a specific token in the keyring

```sh
GODYL_GITHUB_TOKEN=token godyl --keyring auth store github-token
```

### Remove all authentication tokens

```sh
godyl auth rm
```

### Remove a specific token

```sh
godyl auth rm github-token
```
