---
layout: default
title: Configuration
---

# Configuration

Godyl is configured in several ways

Flags are set (in order of priority) by:

1. Command-line flags
2. Environment variables
3. `.env` file(s)
4. `.yml` file

All of above will be merged for run-time settings and to form defaults for each tool being processed.

## Commandline flags

See [Command Reference]({{ site.baseurl }}/commands/commands#content-start) and sub-commands for details on available flags and their default values.

## Environment variables

Environment variables are available for all flag arguments and are prefixed with `GODYL_` and further with `_<subcommand>` for each subcommand.

For the `install` and `download` subcommands, the subcommand prefix is `TOOL_`.

The `.env` files follow the same format. It further supports a `YAML`-like syntax, where you can use `:` to separate keys and values.

For both cases, all environment variables are also loaded into the run-time environment.

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

See [Default Configuration File]({{ site.baseurl }}/configuration/defaults#content-start) for more details.

## Related Topics

- [Tool Configuration Format]({{ site.baseurl }}/configuration/tools#content-start)
