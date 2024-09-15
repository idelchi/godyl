---
layout: default
title: Cache
parent: Commands
nav_order: 6
---

# Cache Command

The `cache` command allows interaction with `godyl`'s cache.

## Syntax

```sh
godyl cache [path|remove|dump|sync] [flags]
```

## Aliases

- `c`

## Description

`godyl` uses a file-based cache to keep track of the versions of the tools it has downloaded.

When running with the [sync strategy]({{site.baseurl }}/configuration/tools#strategy), `godyl` attempts to retrieve the version of the current tool by trying various flags and arguments (`--version`, `-v`, etc.).
Since this might not be so robust, it will first check the cache to see if a version is recorded there from a previous install.

## Subcommands

| Subcommand                | Description                                                                                                                                                     |
| :------------------------ | :-------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `path`                    | Print the path to the cache file                                                                                                                                |
| `remove`, `rm`            | Remove the cache file                                                                                                                                           |
| `dump [name]` `ls [name]` | Show the contents of the cache file. Optionally show a specific item (by name)                                                                                  |
| `clean`                   | Compares the tools in the cache with the tools installed on the system and updates the cache file accordingly. Currently removes entries or updates the version |
