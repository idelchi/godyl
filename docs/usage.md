---
layout: default
title: Usage
render_with_liquid: false
---

# Usage

The basic syntax for using `godyl` is:

```sh
godyl [flags] [command] [flags]
```

## Available Commands

- `install` - Install tools from YAML files
- `download` - Download and unpack individual tools
- `dump` - Display configuration information
- `update` - Update the godyl application
- `cache` - Manage the godyl cache

## Global Flags

The following global flags are available for all commands:

| Flag                | Environment Variable | Default                 | Description                                  |
| ------------------- | -------------------- | ----------------------- | -------------------------------------------- |
| `--help`, `-h`      | `GODYL_HELP`         | `false`                 | Show help message and exit                   |
| `--version`         | `GODYL_VERSION`      | `false`                 | Show version information and exit            |
| `--dry`             | `GODYL_DRY`          | `false`                 | Run without making any changes (dry run)     |
| `--log`             | `GODYL_LOG`          | `info`                  | Log level (debug, info, warn, error, silent) |
| `--env-file`, `-e`  | `GODYL_ENV_FILE`     | `[".env"]`              | Path to `.env` file(s)                       |
| `--defaults`, `-d`  | `GODYL_DEFAULTS`     | `defaults.yml`          | Path to defaults file                        |
| `--github-token`    | `GODYL_GITHUB_TOKEN` | `${GODYL_GITHUB_TOKEN}` | GitHub token for authentication              |
| `--gitlab-token`    | `GODYL_GITLAB_TOKEN` | `${GODYL_GITLAB_TOKEN}` | GitLab token for authentication              |
| `--url-token`       | `GODYL_URL_TOKEN`    | `${GODYL_URL_TOKEN}`    | URL token for authentication                 |
| `--cache-dir`, `-c` | `GODYL_CACHE_DIR`    | `${XDG_CACHE_HOME}`     | Path to cache directory                      |

## Install Command

The `install` command allows you to install tools defined in YAML configuration files:

```sh
godyl install [[tools.yml]...|STDIN] --output ./bin
```

If no file is specified, `godyl` defaults to using `tools.yml` in the current directory.

If the argument is set to `-`, `godyl` will read from `stdin`.

## Download Command

The `download` command allows you to download a single tool without a configuration file:

```sh
godyl download idelchi/godyl --output ./bin
```

You can also download multiple tools:

```sh
godyl download idelchi/tcisd idelchi/gogen idelchi/wslint
```

When using the `download` command, the tool will be unarchived directly into the output directory.

Override `os` and `arch` to download a specific binary:

```sh
godyl download idelchi/godyl --os linux --arch amd64 --output ./bin
```

You can also download tools from direct URLs:

```sh
godyl download "https://github.com/idelchi/go-next-tag/releases/download/v0.0.1/go-next-tag_{{ .OS }}_{{ .ARCH }}.tar.gz" --output ./bin
```

## Dump Command

Display various configuration settings and information:

```sh
godyl dump [config|defaults|env|platform|tools]
```

Subcommands:

- `config` - Display the current configuration settings
- `defaults` - Display the default configuration settings
- `env` - Display environment variables that affect the application
- `platform` - Display information about the current platform
- `tools` - Display information about available tools
- `cache` - Display information about the cache

For example, install all tools that were embedded when the application was built:

```sh
godyl dump tools | godyl install - --output ./bin
```

## Update Command

Update the godyl application to the latest version:

```sh
godyl update [flags]
```

> **Note**: On Windows, this will launch a background process to clean up the old version.

## Cache Command

Manage the godyl cache:

```sh
godyl cache [flags]
```

Use the `--delete` flag to delete the cache:

```sh
godyl cache --delete
```
