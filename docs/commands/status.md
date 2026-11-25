---
layout: default
title: status
parent: Commands
nav_order: 3
---

# Status Command

The `status` command allows you to check whether there are tools missing or syncs to be made.

## Syntax

```sh
godyl [flags] status [tools.yml|-]...
```

## Aliases

- `diff`
- `s`

## Description

The `status` command checks the status of the tools defined in the provided YAML file(s) or from standard input (STDIN). It compares the installed versions of the tools with the versions specified in the YAML file(s) and reports any discrepancies.

It reads the installed-version cache when available, but does not create or update it.

## Exit Codes

| Code | Meaning |
| :--- | :------ |
| `0` | Nothing to install under the configured strategies, including when no tools match the tags. |
| `1` | The check failed, even if other tools need installation. |
| `2` | At least one selected tool needs installation and the check completed without errors. |

## Flags

| Flag             | Environment Variable  | Default     | Description                                 |
| :--------------- | :-------------------- | :---------- | :------------------------------------------ |
| `--output`, `-o` | `GODYL_STATUS_OUTPUT` | `./bin`     | Output path for the downloaded tools        |
| `--tags`, `-t`   | `GODYL_STATUS_TAGS`   | `[!native]` | Tags to filter tools by. Use `!` to exclude |

`tags` may use wildcards `*` which matches any sequence of characters.

## Examples

### Check the status of all tools

```sh
godyl status tools.yml
```

### Check the status of a specific tool

```sh
godyl status tools.yml --tags idelchi/godyl
```
