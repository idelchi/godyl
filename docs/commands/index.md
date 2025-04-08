---
layout: default
title: Commands
has_children: true
---

# Command Reference

Godyl provides several commands to help you manage your tools. This section provides detailed information about each command and its options.

## Available Commands

| Command                              | Description                           |
| ------------------------------------ | ------------------------------------- |
| [`install`](install#content-start)   | Install tools from YAML files         |
| [`download`](download#content-start) | Download and install individual tools |
| [`dump`](dump#content-start)         | Display configuration information     |
| [`update`](update#content-start)     | Update the godyl application          |
| [`cache`](cache#content-start)       | Manage the godyl cache                |

## Global Flags

These flags are available for all commands:

| Flag                 | Environment Variable     | Default                 | Description                                  |
| -------------------- | ------------------------ | ----------------------- | -------------------------------------------- |
| `--help`, `-h`       | `GODYL_HELP`             | `false`                 | Show help message and exit                   |
| `--version`          | `GODYL_VERSION`          | `false`                 | Show version information and exit            |
| `--dry`              | `GODYL_DRY`              | `false`                 | Run without making any changes (dry run)     |
| `--log`              | `GODYL_LOG`              | `info`                  | Log level (debug, info, warn, error, silent) |
| `--env-file`, `-e`   | `GODYL_ENV_FILE`         | `[".env"]`              | Path to `.env` file(s)                       |
| `--defaults`, `-d`   | `GODYL_DEFAULTS`         | `defaults.yml`          | Path to defaults file                        |
| `--github-token`     | `GODYL_GITHUB_TOKEN`     | `${GODYL_GITHUB_TOKEN}` | GitHub token for authentication              |
| `--gitlab-token`     | `GODYL_GITLAB_TOKEN`     | `${GODYL_GITLAB_TOKEN}` | GitLab token for authentication              |
| `--url-token`        | `GODYL_URL_TOKEN`        | `${GODYL_URL_TOKEN}`    | URL token for authentication                 |
| `--url-token-header` | `GODYL_URL_TOKEN_HEADER` | `Authorization`         | URL header for authentication                |
| `--cache-dir`, `-c`  | `GODYL_CACHE_DIR`        | `${XDG_CACHE_HOME}`     | Path to cache directory                      |
| `--cache-type`       | `GODYL_CACHE_TYPE`       | `file`                  | Type of cache (file, sqlite)                 |

## Command Structure

All Godyl commands follow a consistent structure:

```sh
godyl [global flags] [command] [command flags] [arguments]
```

For detailed information about each command, please refer to the individual command pages linked above.
