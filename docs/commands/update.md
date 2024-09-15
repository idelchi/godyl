---
layout: default
title: update
parent: Commands
nav_order: 5
---

# Update Command

![Godyl in Action (Update)]({{ site.baseurl }}/assets/gifs/update.gif)

The `update` command updates the `godyl` application itself to the latest version or a specified version.

## Syntax

```sh
godyl [flags] update [flags]
```

## Aliases

- `upgrade`
- `u`

## Description

The `update` command allows you to keep your `godyl` installation up to date by downloading and installing the latest version from GitHub.

## Flags

| Flag              | Environment Variable   | Default | Description                                                                                                    |
| :---------------- | :--------------------- | :------ | :------------------------------------------------------------------------------------------------------------- |
| `--version`, `-v` | `GODYL_UPDATE_VERSION` | `""`    | Version to download (empty means latest)                                                                       |
| `--pre`           | `GODYL_UPDATE_PRE`     | `false` | Include pre-releases                                                                                           |
| `--check`         | `GODYL_UPDATE_CHECK`   | `false` | Check for updates                                                                                              |
| `--cleanup`       | `GODYL_UPDATE_CLEANUP` | `false` | Cleanup old versions (Windows only; see [Platform-Specific Considerations](#platform-specific-considerations)) |
| `--force`         | `GODYL_UPDATE_FORCE`   | `false` | Force update even if the current version is the latest                                                         |

## Examples

### Update to the latest stable version

```sh
godyl update
```

This will download and install the latest stable version of `godyl`.

### Update to a specific version

```sh
godyl update --version v0.1.0
```

This will download and install version 0.1.0 of `godyl`.

### Include pre-releases when updating

```sh
godyl update --pre
```

This will include pre-release versions when determining the latest version to install.

## Platform-Specific Considerations

### Windows

On Windows, the running binary cannot be directly replaced. Use the `--cleanup` flag/option to launch a background process to remove up the old version.
