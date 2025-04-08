---
layout: default
title: Configuration
nav_order: 4
---

# Configuration

`godyl` can be configured in multiple ways (in order of priority):

1. Command-line flags
2. Environment variables
3. `.env` file(s)
4. `defaults.yml` file

## Global Flags

The following global flags are available for all commands:

| Flag                 | Environment Variable     | Default                          | Description                                  |
| -------------------- | ------------------------ | -------------------------------- | -------------------------------------------- |
| `--help`, `-h`       | `GODYL_HELP`             | `false`                          | Show help message and exit                   |
| `--version`          | `GODYL_VERSION`          | `false`                          | Show version information and exit            |
| `--dry`              | `GODYL_DRY`              | `false`                          | Run without making any changes (dry run)     |
| `--log`              | `GODYL_LOG`              | `info`                           | Log level (debug, info, warn, error, silent) |
| `--env-file`, `-e`   | `GODYL_ENV_FILE`         | `[".env"]`                       | Path to `.env` file(s).                      |
| `--defaults`, `-d`   | `GODYL_DEFAULTS`         | `defaults.yml`                   | Path to defaults file.                       |
| `--github-token`     | `GODYL_GITHUB_TOKEN`     | `${GODYL_GITHUB_TOKEN}`          | GitHub token for authentication              |
| `--gitlab-token`     | `GODYL_GITLAB_TOKEN`     | `${GODYL_GITLAB_TOKEN}`          | GitLab token for authentication              |
| `--url-token`        | `GODYL_URL_TOKEN`        | `${GODYL_URL_TOKEN}`             | URL token for authentication                 |
| `--url-token-header` | `GODYL_URL_TOKEN_HEADER` | `Authorization`                  | URL header for authentication                |
| `--cache-dir`, `-c`  | `GODYL_CACHE_DIR`        | `${XDG_CACHE_HOME}` (or similar) | Path to cache directory                      |
| `--cache-type`       | `GODYL_CACHE_TYPE`       | `file`                           | Type of cache (file, sqlite)                 |

For `--env-file` and `--defaults`, the defaults are used only if no issue is encountered while loading them.

## Environment Variables

In addition to the environment variables that correspond to flags, the following environment variables will be read directly from the environment (not from `.env` files):

- `--github-token` from `GITHUB_TOKEN` or `GH_TOKEN`
- `--gitlab-token` from `GITLAB_TOKEN`
- `--url-token` from `URL_TOKEN`

## Tool-specific Flags

The following flags are available for tool-related commands (`install` and `download`):

| Flag                    | Environment Variable       | Default       | Description                                    |
| ----------------------- | -------------------------- | ------------- | ---------------------------------------------- |
| `--output`, `-o`        | `GODYL_TOOL_OUTPUT`        | `./bin`       | Output path for the downloaded tools           |
| `--tags`, `-t`          | `GODYL_TOOL_TAGS`          | `["!native"]` | Tags to filter tools by. Use `!` to exclude    |
| `--source`              | `GODYL_TOOL_SOURCE`        | `github`      | Source from which to install the tools         |
| `--strategy`            | `GODYL_TOOL_STRATEGY`      | `none`        | Strategy to use for updating tools             |
| `--os`                  | `GODYL_TOOL_OS`            | `""`          | Operating system to use for downloading        |
| `--arch`                | `GODYL_TOOL_ARCH`          | `""`          | Architecture to use for downloading            |
| `--parallel`, `-j`      | `GODYL_TOOL_PARALLEL`      | `0`           | Number of parallel downloads (0 is unlimited)  |
| `--no-verify-ssl`, `-k` | `GODYL_TOOL_NO_VERIFY_SSL` | `false`       | Skip SSL verification                          |
| `--hint`                | `GODYL_TOOL_HINT`          | `[""]`        | Add hint patterns with weight 1                |
| `--version`, `-v`       | `GODYL_TOOL_VERSION`       | `""`          | Version to download (only used for `download`) |
| `--show`, `-s`          | `GODYL_TOOL_SHOW`          | `false`       | Show the configuration and exit                |

## Update Command Flags

The following flags are available for the `update` command:

| Flag                    | Environment Variable         | Default | Description           |
| ----------------------- | ---------------------------- | ------- | --------------------- |
| `--no-verify-ssl`, `-k` | `GODYL_UPDATE_NO_VERIFY_SSL` | `false` | Skip SSL verification |
| `--version`, `-v`       | `GODYL_UPDATE_VERSION`       | `""`    | Version to download   |
| `--pre`                 | `GODYL_UPDATE_PRE`           | `false` | Include pre-releases  |

## Cache Command Flags

The following flags are available for the `cache` command:

| Flag             | Environment Variable | Default | Description      |
| ---------------- | -------------------- | ------- | ---------------- |
| `--delete`, `-d` | `GODYL_CACHE_DELETE` | `false` | Delete the cache |

## Dump Tools Command Flags

The following flags are available for the `dump tools` command:

| Flag          | Environment Variable    | Default | Description                |
| ------------- | ----------------------- | ------- | -------------------------- |
| `--full`, `f` | `GODYL_DUMP_TOOLS_FULL` | `false` | Show full tool information |

## Priority Order and Defaults

Settings can be specified in multiple ways, with the following order of priority:

1. Fields in the `tools.yml` definition
2. Command-line flags
3. Environment variables
4. Values in `.env` file(s)
5. Values in `defaults.yml`
6. Default configuration embedded in the application

## Examples

### Using Environment Variables

```sh
GODYL_TOOL_OUTPUT=~/.local/bin GODYL_GITHUB_TOKEN=gh_token godyl install
```

### Using an .env File

Create a `.env` file:

```
GODYL_TOOL_OUTPUT=~/.local/bin
GODYL_GITHUB_TOKEN=gh_token
```

Then run:

```sh
godyl install
```

### Using Command-line Flags

```sh
godyl --output ~/.local/bin --github-token gh_token install
```

### Using a defaults.yml File

Create a `defaults.yml` file:

```yaml
output: ~/.local/bin
source:
  github:
    token: gh_token
```

Then run:

```sh
godyl install
```
