---
layout: default
title: Commands
nav_order: 4
has_children: true
---

# Command Reference

`godyl` provides several commands to help you manage your tools. This section provides detailed information about each command and its options.

## Available Commands

### Main commands

| Command                                            | Description                           |
| :------------------------------------------------- | :------------------------------------ |
| [`install`]({{ site.baseurl }}/commands/install)   | Install tools from YAML files         |
| [`download`]({{ site.baseurl }}/commands/download) | Download and install individual tools |
| [`update`]({{ site.baseurl }}/commands/update)     | Update the godyl application          |

### Auxiliary commands

| Command                                            | Description                         |
| :------------------------------------------------- | :---------------------------------- |
| [`status`]({{ site.baseurl }}/commands/status)     | Check the status of installed tools |
| [`dump`]({{ site.baseurl }}/commands/dump)         | Display configuration information   |
| [`cache`]({{ site.baseurl }}/commands/cache)       | Manage the cache                    |
| [`config`]({{ site.baseurl }}/commands/config)     | Manage the configuration            |
| [`auth`]({{ site.baseurl }}/commands/auth)         | Manage the authentication tokens    |
| [`validate`]({{ site.baseurl }}/commands/validate) | Validate the configuration          |
| [`paths`]({{ site.baseurl }}/commands/paths)       | Show active filesystem paths        |
| [`version`]({{ site.baseurl }}/commands/version)   | Display the current version         |

## Global Flags

The following global flags are available:

| Flag                    | Environment Variable  | Default                               | Description                                          |
| :---------------------- | :-------------------- | :------------------------------------ | :--------------------------------------------------- |
| `--log-level`, `-l`     | `GODYL_LOG_LEVEL`     | `info`                                | Log level (silent, debug, info, warn, error, always) |
| `--parallel`, `-j`      | `GODYL_PARALLEL`      | `runtime.NumCPU()`                    | Parallelism. 0 means unlimited.                      |
| `--cache-dir`           | `GODYL_CACHE_DIR`     | `${XDG_CACHE_HOME}/godyl`             | Path to cache directory                              |
| `--no-cache`            | `GODYL_NO_CACHE`      | `false`                               | Disable cache                                        |
| `--no-verify-ssl`, `-k` | `GODYL_NO_VERIFY_SSL` | `false`                               | Skip SSL verification                                |
| `--no-progress`         | `GODYL_NO_PROGRESS`   | `false`                               | Disable progress bar                                 |
| `--show`, `-s`          | `GODYL_SHOW`          | `false`                               | Show the parsed configuration and exit               |
| `--config-file`, `-c`   | `GODYL_CONFIG_FILE`   | `${XDG_CONFIG_HOME}/godyl/godyl.yml`  | Path to config file                                  |
| `--env-file`, `-e`      | `GODYL_ENV_FILE`      | `[".env"]`                            | Paths to .env files                                  |
| `--defaults`, `-d`      | `GODYL_DEFAULTS`      | `defaults.yml`                        | Path to defaults file                                |
| `--inherit`             | `GODYL_INHERIT`       | `default`                             | Default to inherit from when unset in the tool spec  |
| `--github-token`        | `GODYL_GITHUB_TOKEN`  | See [authentication](#authentication) | GitHub token for authentication                      |
| `--gitlab-token`        | `GODYL_GITLAB_TOKEN`  | See [authentication](#authentication) | GitLab token for authentication                      |
| `--url-token`           | `GODYL_URL_TOKEN`     | See [authentication](#authentication) | URL token for authentication                         |
| `--error-file`          | `GODYL_ERROR_FILE`    | ``                                    | Path to error log file. Empty means stdout.          |
| `--keyring`             | `GODYL_KEYRING`       | `false`                               | Enable usage of system keyring                       |
| `--verbose`, `-v`       | `GODYL_VERBOSE`       | `false`                               | Increase verbosity (can be used multiple times)      |
| `--version`             | `GODYL_VERSION`       | `false`                               | Show the current version and exit                    |
| `--help`, `-h`          | `GODYL_HELP`          | `false`                               | Show help for the command and exit                   |

`--show` will display the configuration of the current command and all it's parents, and exit. Also available for all subcommands.
Can be repeated to unmask tokens and other sensitive data.

```sh
godyl -ss
```

If you get a lot of error messages for a run, use `error-file` to log them to a file for inspection.

### Configuration management

Use `config` and `auth` to manage the tool configuration.

> **Note**: Commands that write to the `yaml` configuration file (such as `config set`, `config remove`, `auth store` and `auth remove`) will lead to loss of order and newlines.

### Authentication

Authentication tokens default to the following values (in order of precedence),
if not set anywhere else in the [configuration]({{ site.baseurl }}/configuration/index#configuration):

- `--github-token` defaults to the keyring value (see [auth]({{ site.baseurl }}/commands/auth)) (when using the keyring), or the environment variables (`GITHUB_TOKEN`, `GH_TOKEN`)
- `--gitlab-token` defaults to the keyring value (see [auth]({{ site.baseurl }}/commands/auth)) (when using the keyring), or the environment variables (`GITLAB_TOKEN`, `CI_JOB_TOKEN`)
- `--url-token` defaults to the keyring value (see [auth]({{ site.baseurl }}/commands/auth)) (when using the keyring), or the environment variable `URL_TOKEN`

If you'd like to use the keyring for authentication, it's more convenient to set the value in the `yaml` configuration file:

```yaml
keyring: true
```

or as an environment variable:

```sh
GODYL_KEYRING=true
```
