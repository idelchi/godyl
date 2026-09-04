---
layout: default
title: config
parent: Commands
nav_order: 7
---

# Config Command

The `config` command allows interaction with `godyl`'s config file.
Use `godyl config path` to see whether the active file is local or global. See
[Configuration File Location]({{ site.baseurl }}/configuration/index#configuration-file-location)
for the lookup order and platform-specific global directories.

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

### Display the active config file

```sh
godyl config path
```

### Set a key in the config file

```sh
godyl config set dump.tools.embedded true
```

Keys follow the command hierarchy: `dump.tools.embedded` corresponds to
`godyl dump tools --embedded` and `GODYL_DUMP_TOOLS_EMBEDDED`.

Root settings use their names directly. For example, this persistently enables
the Ollama-backed AI fallback and selects its model:

```sh
godyl config set ai true
godyl config set ai-provider ollama
godyl config set ai-model gpt-oss:20b
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
