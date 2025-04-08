---
layout: default
title: Cache Command
parent: Commands
---

# Cache Command

The `cache` command allows interaction with Godyl's cache.

## Syntax

```sh
godyl cache [flags]
```

## Description

Godyl uses a file based cache to keep track of the versions of the tools it has downloaded.

The `cache` command allows you to view information about the cache and manage its contents.

## Flags

| Flag             | Environment Variable | Default | Description      |
| ---------------- | -------------------- | ------- | ---------------- |
| `--delete`, `-d` | `GODYL_CACHE_DELETE` | `false` | Delete the cache |

## Examples

### Output cache file path

```sh
godyl cache
```

### Delete the cache file

```sh
godyl cache --delete
```

## Cache Location

By default, Godyl stores its cache in the standard cache directory for your operating system:

- **Linux**: `~/.cache/godyl` or the location specified by `$XDG_CACHE_HOME/godyl`
- **macOS**: `~/Library/Caches/godyl`
- **Windows**: `%LOCALAPPDATA%\godyl\Cache`

You can override this location by setting the `--cache-dir` global flag or the `GODYL_CACHE_DIR` environment variable.

## How Caching Works

When running in `upgrade` mode, Godyl attempts to retrieve the version of the current tool by trying various flags and arguments (`--version`, `-v`, etc.).
Since this might not be so robust, it will first check it's cache to see if a version is recorded there from a previous install.

## Cache Types

Godyl supports different cache backends:

- **File**: Stores artifacts as files on disk (default)
- **SQLite**: Stores artifacts in a SQLite database

You can select the cache type using the `--cache-type` global flag or the `GODYL_CACHE_TYPE` environment variable.

## Related Topics

- [Global Flags](../commands/index#global-flags)
- [Dump Cache](dump#dump-cache)
