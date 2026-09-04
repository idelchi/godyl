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

| Flag                         | Environment Variable       | Default                               | Description                                          |
| :--------------------------- | :------------------------- | :------------------------------------ | :--------------------------------------------------- |
| `--log-level`, `-l`          | `GODYL_LOG_LEVEL`          | `info`                                | Log level (silent, debug, info, warn, error, always) |
| `--parallel`, `-j`           | `GODYL_PARALLEL`           | `0`                                   | Parallelism. 0 means unlimited.                      |
| `--cache-dir`                | `GODYL_CACHE_DIR`          | `~/.local/share/godyl`                | Path to cache directory                              |
| `--go`                       | `GODYL_GO`                 | ``                                    | Path to go binary for go source installs             |
| `--tmp`                      | `GODYL_TMP`                | ``                                    | Path to temporary directory                          |
| `--no-cache`                 | `GODYL_NO_CACHE`           | `false`                               | Disable cache                                        |
| `--no-verify-ssl`, `-k`      | `GODYL_NO_VERIFY_SSL`      | `false`                               | Skip SSL verification                                |
| `--no-progress`              | `GODYL_NO_PROGRESS`        | `false`                               | Disable progress bar                                 |
| `--no-verify-checksum`, `-C` | `GODYL_NO_VERIFY_CHECKSUM` | `false`                               | Skip checksum verification                           |
| `--ai`                       | `GODYL_AI`                 | `false`                               | Use AI to resolve ambiguous release asset matches    |
| `--ai-provider`              | `GODYL_AI_PROVIDER`        | `ollama`                              | AI provider (`ollama` or `openai`)                   |
| `--ai-model`                 | `GODYL_AI_MODEL`           | Ollama: `gpt-oss:20b`; OpenAI: `gpt-5-nano` | Override the AI model                          |
| `--ai-url`                   | `GODYL_AI_URL`             | Ollama: `http://localhost:11434/v1`; OpenAI: SDK default | Override the AI provider URL             |
| `--show`, `-s`               | `GODYL_SHOW`               | `0`                                   | Show the parsed configuration and exit               |
| `--config-file`, `-c`        | `GODYL_CONFIG_FILE`        | `godyl.yml`                           | Path to config file                                  |
| `--env-file`, `-e`           | `GODYL_ENV_FILE`           | `[".env"]`                            | Paths to .env files                                  |
| `--defaults`, `-d`           | `GODYL_DEFAULTS`           | `defaults.yml`                        | Path to defaults file                                |
| `--inherit`                  | `GODYL_INHERIT`            | `default`                             | Default to inherit from when unset in the tool spec  |
| `--github-token`             | `GODYL_GITHUB_TOKEN`       | See [authentication](#authentication) | GitHub token for authentication                      |
| `--gitlab-token`             | `GODYL_GITLAB_TOKEN`       | See [authentication](#authentication) | GitLab token for authentication                      |
| `--url-token`                | `GODYL_URL_TOKEN`          | See [authentication](#authentication) | URL token for authentication                         |
| `--error-file`               | `GODYL_ERROR_FILE`         | ``                                    | Path to error log file. Empty means stdout.          |
| `--keyring`                  | `GODYL_KEYRING`            | `false`                               | Enable usage of system keyring                       |
| `--verbose`, `-v`            | `GODYL_VERBOSE`            | `0`                                   | Increase verbosity (can be used multiple times)      |
| `--version`                  |                            |                                       | Show the current version and exit                    |
| `--help`, `-h`               |                            |                                       | Show help for the command and exit                   |

`--show` will display the configuration of the current command and all it's parents, and exit. Also available for all subcommands.
Can be repeated to unmask tokens and other sensitive data.

```sh
godyl -ss
```

If you get a lot of error messages for a run, use `error-file` to log them to a file for inspection.

Running with `GODYL_DEBUG=true` will enable (extremely verbose) additional debug logging.

### AI fallback

`--ai` leaves normal matching unchanged. It asks a model to select one exact
release asset only when deterministic matching cannot distinguish between
equally ranked, qualified candidates. Unmatched assets, invalid hints, source
API failures, and successful matches never invoke a model.

The implementation uses Fantasy's provider interface. Ollama is connected
through its OpenAI-compatible API and OpenAI uses Fantasy's native adapter, so
additional provider adapters can be added without changing asset matching.

Ollama is the default and uses `http://localhost:11434/v1` with `gpt-oss:20b`:

```sh
godyl --ai install tools.yml
```

OpenAI defaults to `gpt-5-nano` and reads `GODYL_AI_API_KEY`, falling back to
`OPENAI_API_KEY`:

```sh
GODYL_AI_API_KEY="${OPENAI_API_KEY}" godyl --ai --ai-provider openai install tools.yml
```

The model receives only the tied best candidates. Its answer is accepted only
when it is exactly one of those release asset names. Use `--ai-model` and
`--ai-url` to override the defaults; the same values can be supplied through the
corresponding `GODYL_` environment variables or root YAML configuration keys.
Godyl does not override the provider's input context window. Ollama therefore
uses the context length configured for its model runner, while OpenAI applies
the selected model's API limits.

The complete root YAML configuration is:

```yaml
ai: true
ai-provider: ollama
ai-model: gpt-oss:20b
ai-url: http://localhost:11434/v1
# ai-api-key: secret-value
```

The equivalent environment variables are `GODYL_AI`,
`GODYL_AI_PROVIDER`, `GODYL_AI_MODEL`, `GODYL_AI_URL`, and
`GODYL_AI_API_KEY`. `ai-api-key` is intentionally configuration/env-only and
has no CLI flag. For OpenAI, `OPENAI_API_KEY` is used when
`GODYL_AI_API_KEY` is unset. See [Configuration]({{ site.baseurl }}/configuration/)
for precedence.

### AI suggestions

`godyl install --suggest` is a non-mutating diagnostic mode for deterministic
asset matching failures. It does not require `--ai`, and the two modes cannot be
enabled together: `--ai` resolves an ambiguity for the current run, while
`--suggest` proposes a repeatable hint configuration.

For an ambiguous match, the model receives only the equally ranked best
candidates. When no candidate qualifies, it receives the complete release asset
list. Invalid hints, empty releases, source API errors, authentication failures,
and other non-matching errors are reported without consulting the model.

The request contains matching-related context such as the tool name, source,
version, target platform, installation mode, current hints, and asset names. It
does not include tokens, environment variables, commands, or local paths.

The model returns an exact asset name, a short explanation, and a complete
replacement hint list. Godyl parses those hints and reruns the deterministic
matcher. A suggestion is displayed as verified only when it uniquely resolves
to the model's chosen asset.

```sh
godyl install --suggest tools.yml
godyl --ai-model gpt-oss:20b install --suggest tools.yml
```

Suggestion mode implies dry-run behavior, processes every tool, does not update
the cache, and retains a failing exit status while any tool remains unresolved.

### Configuration management

Use `config` and `auth` to manage the tool configuration.

> **Note**: Commands that write to the `yaml` configuration file (such as `config set`, `config remove`, `auth store` and `auth remove`) will lead to loss of order and newlines.

### Authentication

Authentication tokens default to the following values (in order of precedence),
if not set anywhere else in the [configuration]({{ site.baseurl }}/configuration/index#configuration):

- `--github-token` defaults to the keyring value (see [auth]({{ site.baseurl }}/commands/auth)) (when using the keyring), or the environment variables (`GODYL_GITHUB_TOKEN`, `GITHUB_TOKEN`, `GH_TOKEN`)
- `--gitlab-token` defaults to the keyring value (see [auth]({{ site.baseurl }}/commands/auth)) (when using the keyring), or the environment variables (`GODYL_GITLAB_TOKEN`, `GITLAB_TOKEN`, `CI_JOB_TOKEN`)
- `--url-token` defaults to the keyring value (see [auth]({{ site.baseurl }}/commands/auth)) (when using the keyring), or the environment variables (`GODYL_URL_TOKEN`, `URL_TOKEN`)

If you'd like to use the keyring for authentication, it's more convenient to set the value in the `yaml` configuration file:

```yaml
keyring: true
```

or as an environment variable:

```sh
GODYL_KEYRING=true
```
