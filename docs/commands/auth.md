---
layout: default
title: auth
parent: Commands
nav_order: 8
---

# Auth Command

The `auth` command provides a convenient way of adding or removing tokens from either the configuration file or the system keyring.

## Syntax

```sh
godyl [flags] auth [store|remove|status] [flags]
```

## Subcommands

| Subcommand                           | Description                                |
| :----------------------------------- | :----------------------------------------- |
| `store [token]...`                   | Store tokens from the parsed configuration |
| `remove [token]...`, `rm [token]...` | Remove authentication tokens               |
| `status`                             | Show the status of authentication tokens   |

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

### Show the status of authentication tokens

```sh
godyl auth status
```
