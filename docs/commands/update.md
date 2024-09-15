---
layout: default
title: Update Command
parent: Commands
nav_order: 5
---

# Update Command

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

| Flag                    | Environment Variable         | Default | Description                                                                                                    |
| :---------------------- | :--------------------------- | :------ | :------------------------------------------------------------------------------------------------------------- |
| `--no-verify-ssl`, `-k` | `GODYL_UPDATE_NO_VERIFY_SSL` | `false` | Skip SSL verification                                                                                          |
| `--version`, `-v`       | `GODYL_UPDATE_VERSION`       | `""`    | Version to download (empty means latest)                                                                       |
| `--pre`                 | `GODYL_UPDATE_PRE`           | `false` | Include pre-releases                                                                                           |
| `--check`               | `GODYL_UPDATE_CHECK`         | `false` | Check for updates                                                                                              |
| `--cleanup`             | `GODYL_UPDATE_CLEANUP`       | `false` | Cleanup old versions (Windows only; see [Platform-Specific Considerations](#platform-specific-considerations)) |
| `--force`               | `GODYL_UPDATE_FORCE`         | `false` | Force update even if the current version is the latest                                                         |

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

On Windows, the update process launches a background process to clean up the old version, as the running binary cannot be directly replaced. This process will wait for the current `godyl` process to exit before completing the update. Can be enabled with the `--cleanup` flag.
