---
layout: default
title: Cache Command
---

# Cache Command

The `cache` command allows interaction with `godyl`'s cache.

## Syntax

```sh
godyl cache [path|remove|show|sync] [flags]
```

## Aliases

- `c`

## Description

`godyl` uses a file-based cache to keep track of the versions of the tools it has downloaded.

When running in `sync` mode, `godyl` attempts to retrieve the version of the current tool by trying various flags and arguments (`--version`, `-v`, etc.).
Since this might not be so robust, it will first check it's cache to see if a version is recorded there from a previous install.

## Subcommands

| Subcommand     | Description                                                                                                  |
| -------------- | ------------------------------------------------------------------------------------------------------------ |
| `path`         | Print the path to the cache file.                                                                            |
| `remove`, `rm` | Remove the cache file.                                                                                       |
| `show [name]`  | Show the contents of the cache file. Optionally show a specific item (by name).                              |
| `sync`         | Compare the tools in the cache with the tools installed on the system and update the cache file accordingly. |

Running without any flags will print the cache file path.

`sync` will compare the tools in the cache with the tools installed on the system and update the cache file accordingly.

## Related Topics

- [Global Flags]({{ site.baseurl }}/commands/commands#global-flags)
- [Dump Cache]({{ site.baseurl }}/commands/dump#dump-cache)
