---
layout: default
title: cache
parent: Commands
nav_order: 6
---

# Cache Command

The `cache` command allows interaction with `godyl`'s cache.

## Syntax

```sh
godyl cache [path|remove|clean] [flags]
```

## Description

`godyl` uses a file-based cache to keep track of the versions of the tools it has downloaded.

When running with the [sync strategy]({{site.baseurl }}/configuration/tools#strategy), `godyl` attempts to retrieve the version of the current tool by trying various flags and arguments (`--version`, `-v`, etc.).
Since this might not be so robust, it will first check the cache to see if a version is recorded there from a previous install.

## Subcommands

| Subcommand                          | Description                                                                                                    |
| :---------------------------------- | :------------------------------------------------------------------------------------------------------------- |
| `path`                              | Print the path to the cache file                                                                               |
| `remove [name]...`, `rm  [name]...` | Remove entries in the cache file                                                                               |
| `clean`                             | Compares the tools in the cache with the tools installed on the system and updates the cache file accordingly. |
