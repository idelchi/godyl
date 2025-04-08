---
layout: default
title: Update Command
parent: Commands
---

# Update Command

The `update` command updates the Godyl application itself to the latest version or a specified version.

## Syntax

```sh
godyl update [flags]
```

## Description

The `update` command allows you to keep your Godyl installation up to date by downloading and installing the latest version from GitHub. This command handles the update process automatically, ensuring that you always have access to the latest features and bug fixes.

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

## How the Update Process Works

1. The `update` command first checks the GitHub repository for the latest available version.
2. If a newer version is available (or a specific version is requested), it downloads the appropriate release asset for your platform.
3. The downloaded asset is extracted to a temporary location.
4. The Godyl binary is replaced with the new version.
5. On Windows, a background process is launched to clean up the old version.

## Platform-Specific Considerations

### Windows

On Windows, the update process launches a background process to clean up the old version, as the running binary cannot be directly replaced. This process will wait for the current Godyl process to exit before completing the update.

### Linux and macOS

On Linux and macOS, the update process directly replaces the Godyl binary. If the binary is in a location that requires elevated privileges, you may need to run the update command with `sudo`:

```sh
sudo godyl update
```

## Troubleshooting

If you encounter any issues during the update process, you can try the following:

1. Use the `--no-verify-ssl` flag if you're experiencing SSL certificate validation issues:

   ```sh
   godyl update --no-verify-ssl
   ```

2. If you're unable to update to the latest version, try specifying a specific version:

   ```sh
   godyl update --version v0.1.0
   ```

3. If the update command fails, you can always manually download and install Godyl from the GitHub releases page.

## Related Topics

- [Installation Guide](../installation.html)
- [Platform Detection](../advanced-features.html#platform-inference)
