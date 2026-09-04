---
layout: default
title: Configuration
nav_order: 3
has_children: true
---

# Configuration

`godyl` can be configured in several ways.

Settings are applied (in order of highest to lowest priority) by:

1. Definition in `tools.yml`
2. Command-line flags
3. Environment variables
4. `.env` file(s)
5. `.yml` file
6. `defaults.yml` file (embedded or passed as configuration)
7. Command-line flag default values

In addition, tokens can be set in the keyring, or from a few other commonly used environment variables. See [Authentication]({{ site.baseurl }}/commands/index#authentication) for more details.

All of above will be merged for run-time settings and to form the complete definition for each tool.

## Configuration File Location

Unless `--config-file` or `GODYL_CONFIG_FILE` selects a specific file, `godyl`
uses the first existing configuration file in this order:

1. `./godyl.yaml`
2. `./godyl.yml`
3. The global `godyl.yaml`
4. The global `godyl.yml`

If none exists, the active path is the global `godyl.yml`. The global
configuration directory is platform-specific:

| Platform                     | Directory                                                                        |
| :--------------------------- | :------------------------------------------------------------------------------- |
| Linux and other Unix systems | `${XDG_CONFIG_HOME}/godyl`, or `~/.config/godyl` when `XDG_CONFIG_HOME` is unset |
| macOS                        | `~/Library/Application Support/godyl`                                            |
| Windows                      | `%AppData%\godyl`                                                                |

Print the exact active path on the current system with:

```sh
godyl config path
```

A configuration file in the current directory intentionally takes precedence
over the global file. Consequently, `godyl config set`, `godyl config remove`,
and authentication commands that are not using the keyring operate on the
active file shown by `godyl config path`.

## Commandline flags

See [Command Reference]({{ site.baseurl }}/commands) and sub-commands for details on available flags and their default values.

## Environment variables

{% raw %}

Environment variables are available for all flag arguments and are prefixed with `GODYL_` and further with `<SUBCOMMAND>_` for each subcommand.

The `.env` files follow the same format.

Examples:

```sh
# Set the output directory for the `install` command
GODYL_INSTALL_OUTPUT=~/.local/bin

# Set the full flag for the `dump tools` command
GODYL_DUMP_TOOLS_FULL=true
```

The AI asset-matching fallback is configured through the same mechanism:

```sh
GODYL_AI=true
GODYL_AI_PROVIDER=ollama
GODYL_AI_MODEL=gpt-oss:20b
GODYL_AI_URL=http://localhost:11434/v1
```

For OpenAI, set `GODYL_AI_PROVIDER=openai` and provide
`GODYL_AI_API_KEY`. If that variable is unset, `OPENAI_API_KEY` is used.
The API key deliberately has no command-line flag, avoiding exposure in shell
history and process arguments.

All environment variables are also loaded into the run-time environment, regardless of whether they came directly from the environment or from a `.env` file.

They can be accessed with `{{ .Env.<ENV_VAR> }}`.

{% endraw %}

## YAML Configuration

A YAML file can set the same values as long command-line flags and environment
variables. The names are derived mechanically:

| Scope             | Command-line flag         | Environment variable    | YAML key          |
| :---------------- | :------------------------ | :---------------------- | :---------------- |
| Root              | `--log-level`             | `GODYL_LOG_LEVEL`       | `log-level`       |
| Subcommand        | `godyl install --output`  | `GODYL_INSTALL_OUTPUT`  | `install.output`  |
| Nested subcommand | `godyl dump tools --full` | `GODYL_DUMP_TOOLS_FULL` | `dump.tools.full` |

YAML preserves hyphens in long flag names. Environment variables replace
hyphens with underscores and use uppercase letters. Short flag aliases do not
affect either name.

In YAML, the dot-separated keys above are represented by nested mappings. A
representative global configuration file looks like this:

```yaml
# Root flags
log-level: info
parallel: 4
no-progress: true
keyring: true
env-file:
  - .env

# Optional AI tie-breaker for ambiguous release-asset matching
ai: true
ai-provider: ollama
ai-model: gpt-oss:20b
ai-url: http://localhost:11434/v1
# ai-api-key: secret-value

# `godyl install` flags
install:
  output: ~/.local/bin
  strategy: sync
  tags:
    - "!native"
  suggest: false

# `godyl download` flags
download:
  output: ./downloads
  source: github

# `godyl dump tools` flags
dump:
  tools:
    embedded: true
    full: true
```

Only include settings that should differ from the built-in defaults. See the
[command reference]({{ site.baseurl }}/commands/) for available flags, or run
`godyl validate` to validate the complete active configuration.

The AI keys are root-command settings. Their defaults are:

| Key           | Environment variable | Default                                                               |
| :------------ | :------------------- | :-------------------------------------------------------------------- |
| `ai`          | `GODYL_AI`           | `false`                                                               |
| `ai-provider` | `GODYL_AI_PROVIDER`  | `ollama`                                                              |
| `ai-model`    | `GODYL_AI_MODEL`     | `gpt-oss:20b` for Ollama; `gpt-5-nano` for OpenAI                     |
| `ai-url`      | `GODYL_AI_URL`       | `http://localhost:11434/v1` for Ollama; OpenAI SDK default for OpenAI |
| `ai-api-key`  | `GODYL_AI_API_KEY`   | `ollama` for Ollama; `OPENAI_API_KEY` fallback for OpenAI             |

CLI flags, environment variables, `.env`, and YAML retain the precedence shown
at the top of this page. For example, `--ai-model` overrides `GODYL_AI_MODEL`,
which overrides `ai-model` in `godyl.yml`.

Suggestion mode is an install-command setting. Its flag, environment variable,
and YAML key are `--suggest`, `GODYL_INSTALL_SUGGEST`, and `install.suggest`.
It should normally be enabled for an individual diagnostic run rather than
stored permanently.

## Defaults Configuration

The `defaults.yml` file is used to set default values for all tools. It supports the same fields as the `tools.yml` file.

An example of sane default values are provided in [defaults.yml](https://github.com/idelchi/godyl/blob/main/defaults.yml) which is also embedded in the binary.

See [Default Configuration File]({{ site.baseurl }}/configuration/defaults) for more details.
