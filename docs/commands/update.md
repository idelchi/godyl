---
layout: default
title: Update Command
---

# Update Command

The `update` command updates the Godyl application itself to the latest version or a specified version.

## Syntax

```sh
godyl [flags] update [flags]
```

## Description

The `update` command allows you to keep your Godyl installation up to date by downloading and installing the latest version from GitHub.

## Flags

| Flag                    | Environment Variable         | Default | Description           |
| ----------------------- | ---------------------------- | ------- | --------------------- |
| `--no-verify-ssl`, `-k` | `GODYL_UPDATE_NO_VERIFY_SSL` | `false` | Skip SSL verification |
| `--version`, `-v`       | `GODYL_UPDATE_VERSION`       | `""`    | Version to download   |
| `--pre`                 | `GODYL_UPDATE_PRE`           | `false` | Include pre-releases  |

## Examples

### Update to the latest stable version

```sh
godyl update
```

This will download and install the latest stable version of Godyl.

### Update to a specific version

```sh
godyl update --version v0.1.0
```

This will download and install version 0.1.0 of Godyl.

### Include pre-releases when updating

```sh
godyl update --pre
```

This will include pre-release versions when determining the latest version to install.

## Platform-Specific Considerations

### Windows

On Windows, the update process launches a background process to clean up the old version, as the running binary cannot be directly replaced. This process will wait for the current Godyl process to exit before completing the update.

## Related Topics

- [Global Flags]({{ site.baseurl }}/commands/commands#global-flags)
