---
layout: default
title: Commands
nav_order: 4
has_children: true
---

# Command Reference

`godyl` provides several commands to help you manage your tools. This section provides detailed information about each command and its options.

## Available Commands

| Command                                            | Description                           |
| :------------------------------------------------- | :------------------------------------ |
| [`install`]({{ site.baseurl }}/commands/install)   | Install tools from YAML files         |
| [`download`]({{ site.baseurl }}/commands/download) | Download and install individual tools |
| [`status`]({{ site.baseurl }}/commands/status)     | Check the status of installed tools   |
| [`dump`]({{ site.baseurl }}/commands/dump)         | Display configuration information     |
| [`update`]({{ site.baseurl }}/commands/update)     | Update the godyl application          |
| [`cache`]({{ site.baseurl }}/commands/cache)       | Manage the godyl cache                |
| [`validate`]({{ site.baseurl }}/commands/validate) | Validate the configuration            |
| [`version`]({{ site.baseurl }}/commands/version)   | Display the current version of godyl  |

## Global Flags

These flags are available for all commands:

| Flag                    | Environment Variable  | Default                                      | Description                                          |
| :---------------------- | :-------------------- | :------------------------------------------- | :--------------------------------------------------- |
| `--log-level`, `-l`     | `GODYL_LOG_LEVEL`     | `info`                                       | Log level (silent, debug, info, warn, error, always) |
| `--parallel`, `-j`      | `GODYL_PARALLEL`      | `runtime.NumCPU()`                           | Parallelism. 0 means unlimited.                      |
| `--cache-dir`           | `GODYL_CACHE_DIR`     | `${XDG_CACHE_HOME}/godyl`                    | Path to cache directory                              |
| `--no-cache`            | `GODYL_NO_CACHE`      | `false`                                      | Disable cache                                        |
| `--no-verify-ssl`, `-k` | `GODYL_NO_VERIFY_SSL` | `false`                                      | Skip SSL verification                                |
| `--no-progress`         | `GODYL_NO_PROGRESS`   | `false`                                      | Disable progress bar                                 |
| `--show`, `-s`          | `GODYL_SHOW`          | `false`                                      | Show the parsed flags and exit                       |
| `--config-file`, `-c`   | `GODYL_CONFIG_FILE`   | `${XDG_CONFIG_HOME}/godyl/godyl.yml`         | Path to config file                                  |
| `--env-file`, `-e`      | `GODYL_ENV_FILE`      | `[".env"]`                                   | Paths to .env files                                  |
| `--defaults`, `-d`      | `GODYL_DEFAULTS`      | `defaults.yml`                               | Path to defaults file                                |
| `--inherit`             | `GODYL_INHERIT`       | `default`                                    | Default to inherit from when unset in the tool spec  |
| `--github-token`        | `GODYL_GITHUB_TOKEN`  | `${GITHUB_TOKEN}` or `${GH_TOKEN}`           | GitHub token for authentication                      |
| `--gitlab-token`        | `GODYL_GITLAB_TOKEN`  | `${GODYL_GITLAB_TOKEN}` or `${GITLAB_TOKEN}` | GitLab token for authentication                      |
| `--url-token`           | `GODYL_URL_TOKEN`     | `${GODYL_URL_TOKEN}` or `${URL_TOKEN}`       | URL token for authentication                         |
| `--error-file`          | `GODYL_ERROR_FILE`    | ``                                           | Path to error log file. Empty means stdout.          |

`--show` will display the configuration of the current command and all it's parents, and exit. Also available for all subcommands.
Can be repeated to unmask tokens and other sensitive data.

```sh
godyl -ss
```

If you get a lot of error messages for a run, use `error-file` to log them to a file for inspection.

`godyl` searched for a `tools.yml` file in the current directory if not given as an argument. Additionally, it will respect the `GODYL_INSTALL_TOOLS` environment variable, as well as the `install.tools` key in the config file.
