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

## Flags

| Flag           | Environment Variable | Default | Description                                 |
| :------------- | :------------------- | :------ | :------------------------------------------ |
| `--tags`, `-t` | `GODYL_STATUS_TAGS`  | `[]`    | Tags to filter tools by. Use `!` to exclude |

## Examples

### Check the status of all tools

```sh
godyl status tools.yml
```

### Check the status of a specific tool

```sh
godyl status tools.yml --tags idelchi/godyl
```
