---
layout: default
title: config
parent: Commands
nav_order: 7
---

# Config Command

The `config` command allows interaction with `godyl`'s config file.

## Syntax

```sh
godyl [flags] config [path|set|remove] [flags]
```

## Aliases

- `cfg`

## Subcommands

| Subcommand                       | Description                       |
| :------------------------------- | :-------------------------------- |
| `path`                           | Print the path to the config file |
| `set <key> <value>`              | Set a value in the config file    |
| `remove [key]...`, `rm [key]...` | Remove entries in the config file |

> **Note**: The `set` command will lead to loss of order and newlines in the config file.

## Examples

### Set a key in the config file

```sh
godyl config set dump.tools.embedded true
```

Root settings use their names directly. For example, this persistently enables
the Ollama-backed AI fallback and selects its model:

```sh
godyl config set ai true
godyl config set ai-provider ollama
godyl config set ai-model gemma3:4b
```

An OpenAI key can be stored as `ai-api-key`, although an environment variable
is preferable for secrets:

```sh
GODYL_AI_API_KEY="${OPENAI_API_KEY}" godyl --ai --ai-provider openai install tools.yml
```

### Remove a key from the config file

```sh
godyl config remove dump
```

### Remove all entries in the config file

```sh
godyl config remove
```
