---
layout: default
title: Cache Command
---

# Cache Command

The `cache` command allows interaction with Godyl's cache.

## Syntax

```sh
godyl cache [flags]
```

## Description

Godyl uses a file based cache to keep track of the versions of the tools it has downloaded.

When running in `upgrade` mode, Godyl attempts to retrieve the version of the current tool by trying various flags and arguments (`--version`, `-v`, etc.).
Since this might not be so robust, it will first check it's cache to see if a version is recorded there from a previous install.

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

## Cache Types

Godyl supports different cache backends:

- **File**: Stores artifacts as files on disk (default)
- **SQLite**: Stores artifacts in a SQLite database

## Related Topics

- [Global Flags]({{ site.baseurl }}/commands/commands#global-flags)
- [Dump Cache]({{ site.baseurl }}/commands/dump#dump-cache)
