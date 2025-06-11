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

## Commandline flags

See [Command Reference]({{ site.baseurl }}/commands) and sub-commands for details on available flags and their default values.

## Environment variables

{% raw %}

Environment variables are available for all flag arguments and are prefixed with `GODYL_` and further with `_<SUBCOMMAND>` for each subcommand.

The `.env` files follow the same format.

Examples:

```sh
# Set the output directory for the `install` command
GODYL_INSTALL_OUTPUT=~/.local/bin

# Set the full flag for the `dump tools` command
GODYL_DUMP_TOOLS_FULL=true
```

All environment variables are also loaded into the run-time environment, regardless of whether they came directly from the environment or from a `.env` file.

As such, they can all be accessed with `{{ .Env.<ENV_VAR> }}`.

{% endraw %}

## YAML Configuration

A `yaml` file can be used to set values for the flags, following the same subcommand convention as environment variables.

For example, to set:

- the `env-file` flag on the root command
- `output` flag for the `install` subcommand
- `format` flag for the `dump` subcommand
- the `full` flag for the `dump tools` subcommand,

you would use the following format in your YAML file:

```yaml
# Root command
env-file:
  - .env

# `install` subcommand
tool:
  output: ~/.local/bin

# `dump` subcommand
dump:
  format: json
  # `dump tools` subcommand
  tools:
    full: true
```

## Defaults Configuration

The `defaults.yml` file is used to set default values for all tools. It supports the same fields as the `tools.yml` file.

An example of sane default values are provided in [defaults.yml](https://github.com/idelchi/godyl/blob/main/defaults.yml) which is also embedded in the binary.

See [Default Configuration File]({{ site.baseurl }}/configuration/defaults) for more details.
