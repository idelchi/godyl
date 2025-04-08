---
layout: default
title: Cache Command
parent: Commands
---

# Cache Command

The `cache` command allows you to manage Godyl's cache for downloaded artifacts.

## Syntax

```sh
godyl cache [flags]
```

## Description

Godyl uses a cache to store downloaded artifacts, which helps to speed up future installations and reduce bandwidth usage. The `cache` command allows you to view information about the cache and manage its contents.

## Flags

| Flag             | Environment Variable | Default | Description      |
| ---------------- | -------------------- | ------- | ---------------- |
| `--delete`, `-d` | `GODYL_CACHE_DELETE` | `false` | Delete the cache |

## Examples

### View cache information

```sh
godyl cache
```

This will display information about the cache, including its location, size, and contents.

### Delete the cache

```sh
godyl cache --delete
```

This will remove all cached artifacts, freeing up disk space. Note that this will cause future installations to re-download artifacts.

## Cache Location

By default, Godyl stores its cache in the standard cache directory for your operating system:

- **Linux**: `~/.cache/godyl` or the location specified by `$XDG_CACHE_HOME/godyl`
- **macOS**: `~/Library/Caches/godyl`
- **Windows**: `%LOCALAPPDATA%\godyl\Cache`

You can override this location by setting the `--cache-dir` global flag or the `GODYL_CACHE_DIR` environment variable.

## How Caching Works

When you download a tool using Godyl, the downloaded artifact is stored in the cache. If you later attempt to download the same tool with the same version, Godyl will use the cached version instead of downloading it again.

The cache is organized by source type, tool name, and version. Each cached artifact includes metadata about when it was downloaded and the URL or source it came from.

## Cache Types

Godyl supports different cache backends:

- **File**: Stores artifacts as files on disk (default)
- **SQLite**: Stores artifacts in a SQLite database

You can select the cache type using the `--cache-type` global flag or the `GODYL_CACHE_TYPE` environment variable.

## Benefits of Caching

Using the cache provides several benefits:

1. **Reduced download time**: Previously downloaded artifacts can be reused without re-downloading.
2. **Offline installation**: Once an artifact is cached, you can install it even without an internet connection.
3. **Reduced bandwidth usage**: Useful for environments with limited or metered connections.
4. **Improved reliability**: If a source becomes temporarily unavailable, cached artifacts can still be used.

## When to Delete the Cache

You might want to delete the cache in the following situations:

- To free up disk space
- If you suspect cached artifacts are corrupted
- After upgrading to a new version of Godyl with incompatible cache format changes
- To force re-downloading of all artifacts

## Related Topics

- [Global Flags](../commands/index.html#global-flags)
- [Installation Strategies](../advanced-features.html#installation-strategies)
