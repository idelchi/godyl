---
layout: default
title: Configuration
---

# Configuration

Godyl can be configured in several ways (in order of priority):

1. Command-line flags
2. Environment variables
3. `.env` file(s)
4. `.yml` file
5. `defaults.yml` file (if not present, the embedded defaults will be used)

All of above will be merged together to form defaults for each tool being processed.

## Commandline flags

See [Command Reference](/godyl/commands/index#content-start) and sub-commands for details on available flags and their default values.

## Environment variables

Environment variables are available for all flag arguments and are prefixed with `GODYL_` and further with `_<subcommand>` for each subcommand.

For the `install` and `download` subcommands, the subcommand prefix is `TOOL_`.

The `.env` files follow the same format. It further supports a `YAML`-like syntax, where you can use `:` to separate keys and values.

## YAML Configuration

A `yaml` file can be used to set values for the flags, following the same subcommand convention as environment variables, but without the `GODYL_` prefix.

For example, to set the `output` flag for the `install` subcommand, and the `full` flag for the `dump tools` subcommand, you would use the following format in your YAML file:

```yaml
tool:
  output: ~/.local/bin

dump:
  tools:
    full: true
```

## Defaults Configuration

The `defaults.yml` file is used to set default values for all tools. It supports the same fields as the `tools.yml` file.

An example of sane default values are provided in [defaults.yml](https://github.com/idelchi/godyl/blob/main/defaults.yml) which is
also embedded in the binary.

See [Default Configuration File](/godyl/commands/defaults#content-start) for more details.

## Related Topics

- [Tool Configuration Format](/godyl/commands/tools#content-start)
