# godyl

[![Go Reference](https://pkg.go.dev/badge/github.com/idelchi/godyl.svg)](https://pkg.go.dev/github.com/idelchi/godyl)
[![Go Report Card](https://goreportcard.com/badge/github.com/idelchi/godyl)](https://goreportcard.com/report/github.com/idelchi/godyl)
[![Build Status](https://github.com/idelchi/godyl/actions/workflows/go.yml/badge.svg)](https://github.com/idelchi/godyl/actions/workflows/go.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

godyl helps with batch-downloading and installing statically compiled binaries from:

- GitHub releases
- URLs
- Go projects

As an alternative to above, custom commands can be used as well.

`godyl` will infer the platform and architecture from the system it is running on, and will attempt to download the appropriate binary.

This uses simple heuristics to infer the correct binary to download, and will not work for all projects.

Most properties can be overridden, with `hints` and `skip` used to help `godyl` make the correct decision.

> [!NOTE]
> This repo is a work in progress!
> Needing both cleaning up and documenting.

> **Warning**
> This repo is a work in progress!
> Needing both cleaning up and documenting.

---

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Customization](#customization)
- [Contributing](#contributing)
- [License](#license)

---

## Installation

### From source

```sh
go install github.com/idelchi/godyl/cmd/godyl@latest
```

## From installation script

```sh
curl -sSL https://raw.githubusercontent.com/idelchi/godyl/refs/heads/dev/scripts/install.sh | sh -s -- -v v0.0 -o ~/.local/bin
```

## Configuration

The tools can be configured by (in order of priority)

- flags to the tool
- environment variables
- a `.env` file

The following flags and their corresponding environment variables are available:

| Flag              | Environment Variable  | Description                                    |
| ----------------- | --------------------- | ---------------------------------------------- |
| `--output`        | `GODYL_OUTPUT`        | Output path for the downloaded tools           |
| `--tags`          | `GODYL_TAGS`          | Tags to filter tools by                        |
| `--defaults`      | `GODYL_DEFAULTS`      | Path to defaults file                          |
| `--update`        | `GODYL_UPDATE`        | Update the tools                               |
| `--strategy`      | `GODYL_STRATEGY`      | Strategy to use for updating tools             |
| `--dry`           | `GODYL_DRY`           | Run without making any changes (dry run)       |
| `--log`           | `GODYL_LOG`           | Log level (DEBUG, INFO, WARN, ERROR)           |
| `--github-token`  | `GODYL_GITHUB_TOKEN`  | GitHub token for authentication                |
| `--source`        | `GODYL_SOURCE`        | Source from which to install the tools         |
| `--dot-env`       | `GODYL_DOT_ENV`       | Path to .env file                              |
| `--help`          | `GODYL_HELP`          | Show help message and exit                     |
| `--show-config`   | `GODYL_SHOW_CONFIG`   | Show the parsed configuration and exit         |
| `--show-defaults` | `GODYL_SHOW_DEFAULTS` | Show the parsed default configuration and exit |
| `--show-env`      | `GODYL_SHOW_ENV`      | Show the parsed environment variables and exit |
| `--show-platform` | `GODYL_SHOW_PLATFORM` | Show the detected platform and exit            |
| `--version`       | `GODYL_VERSION`       | Show version information and exit              |
| `--parallel`      | `GODYL_PARALLEL`      | Number of parallel downloads                   |

The path to the tools file is provided as a positional argument, defaulting to `tools.yml`.

An example [tools.yml](./tools.yml) is provided.

## Defaults

A default configuration may be used to specify default settings for all tools. These will override (or extend in some case) the settings for each tool.

The following is embedded and used by default if no configuration is provided:

[config.yml](./cmd/godyl/defaults.yml)

The example above defines:

- The default output directory for all tools
- Patterns to use for when searching for the executable
- Hints to:
  - use the executable name as a pattern (useful for repositories with multiple binaries, such as `ahmetb/kubectx`)
  - prefer `.zip` files for Windows
- `find` mode for downloading, extracting and finding the executable
- The default source type as GitHub
- `none` strategy to skip tools which already exist
- Settings the environment variable `GH_TOKEN` to the value of `GODYL_GITHUB_TOKEN`

The full set of configuration options are:

type Defaults struct {
Exe Exe
Output string
Platform detect.Platform
Values map[string]any
Fallbacks []string
Hints match.Hints
Source sources.Source
Tags Tags
Strategy Strategy
Extensions []string
Env env.Env
Mode Mode
}

```yaml
exe:
  name: string
  patterns: []

output: string
platform:
  os: string
  architecture:
    type: string
    version: string
  distribution: string
  library: string
  extension: string

values: {}
fallbacks: []
hints:
  - pattern: string
    weight: int
    regex: boolean
    must: boolean
source:
  type: string
  github:
    owner: string
    repo: string
    token: string
  url:
    url: string
    token: string
  go:
    command: string

tags: []
strategy: string
extensions: []
env: {}
mode: string
```

# path to tools file

tools: string

# list of tags to filter tools by

tags: []

# whether to update `godyl` itself

update: boolean

# update strategy for `godyl`

update-strategy: string

# dry run to output chosen tools

dry: boolean

# log level, one of debug, info, warn, error

log: string

# help message

help: boolean

# show configuration

show: boolean

# show version

version: boolean

# number of parallel downloads

parallel: int

````

## Tools

A YAML file controls the tools to download and install. Alternative, if the positional argument to the tool is not a YAML file, it will be treated as a single tool name.

Examples are provided in [tools.yml](./tools.yml) and

```yaml
- ajeetdsouza/zoxide
````

Above is the `simple` form to attempt to download the latest release of `zoxide` from `ajeetdsouza/zoxide`.

The full form is

```yaml
# Name of the tool, can use Go templates
name: ajeetdsouza/zoxide
# Description of the tool
description: A smart autojump tool
# Version of the tool, can use Go templates
version: v{{ .Values.Version }}
# Path to fetch the tool, can use Go templates. Will be inferred if not given
path: ""
# Checksum for the downloaded file (NOT IMPLEMENTED)
checksum: ""
# Output path for the tool
output: "{{ .Output }}"
exe:
  # Name of the executable itself, inferred from name if not given, can use Go templates
  name:
  # Patterns to use for finding the executable, can use Go templates
  patterns:
    - "{{ .Exe.Name }}.*"
platform: "{{ .Platform }}" # Platform detection. Any field not given will be detected from the system.
aliases: # Aliases for the tool
  - z
values: # Arbitrary values map, can be used for templating in other fields
  version: v0.9.6
fallbacks: # List of fallback strategies
  - go
hints: # Hints for matching, can use Go templates in pattern and weight fields
  - pattern: ""
    weight: 1
    regex: false
    must: false
source:
  type: # Source type, can be github, go, or url
  github:
    owner:
    repo:
    token:
  url:
    url:
    token:
  go:
    command:
tags: # Tags for categorizing tools, can use Go templates
  - terminal
strategy: none # Strategy for installation, can be none, upgrade or force
extensions:
  - .gz
skip: false # Whether to skip installation (evaluated as boolean)
test: # Test commands, can use Go templates
  - zoxide --version
```

## Settings

In general, settings can be set in the following ways:

- as a field in the tool definition

  ```yaml
  output: ~/.local/bin
  ```

- as a flag to the tool

  ```sh
  godyl --defaults.output ~/.local/bin
  ```

- as an environment variable

  ```sh
  GODYL_DEFAULTS_OUTPUT=~/.local/bin godyl
  ```

- in the configuration file

  ```yaml
  defaults:
    output: ~/.local/bin
  ```

  pflag.Bool("version", false, "Show the version information and exit")
  pflag.BoolP("help", "h", false, "Show the help information and exit")
  pflag.BoolP("show", "s", false, "Show the configuration and exit")
  pflag.StringP("config", "c", config.Get(), "Path to configuration file")
  pflag.IntP("parallel", "j", 0, "Number of parallel downloads")

  // Selected custom flags
  pflag.String("defaults.source.github.token", "", "GitHub token for API requests")
  pflag.String("defaults.strategy", "none", "")
  pflag.String("defaults.output", "~/.local/bin", "")

  pflag.String("log", string(logger.INFO), "")

  pflag.Bool("dry", false, "")

  pflag.Bool("update", false, "")
  pflag.String("update-strategy", string(tools.Upgrade), "")

  pflag.StringSliceP("tags", "t", []string{"!native"}, "Tags to filter tools by")

The following `config` parameters are available:

| Field               | Type          | `config.yml`          | Flag                           | Environment Variable          | Default                                 |
| ------------------- | ------------- | --------------------- | ------------------------------ | ----------------------------- | --------------------------------------- |
| output              | string        | defaults.output       | `--defaults.output`            | `GODYL_DEFAULTS_OUTPUT`       | `~/.local/bin`                          |
| exe                 | list of dicts | defaults.exe          |                                |                               | [config.yml](./cmd/godyl/config.yml#L3) |
| hints               | list of dicts | defaults.hints        |                                |                               | [config.yml](./cmd/godyl/config.yml#L5) |
| source.type         | string        | defaults.exe.patterns | `defaults.source.type`         | `GODYL_DEFAULTS_SOURCE_TYPE`  | `github`                                |
| source.github.token | string        | defaults.github.token | `defaults.source.github.token` | `GODYL_DEFAULTS_GITHUB_TOKEN` |                                         |

| Field | Type          | `config.yml` | Flag | Environment Variable | Default                                                                                   |
| ----- | ------------- | ------------ | ---- | -------------------- | ----------------------------------------------------------------------------------------- |
| exe   | list of dicts | defaults.exe | -    | -                    | <pre>exe:<br>&nbsp;&nbsp;patterns: "{{ .Exe.Name }}.\*"<br></pre> [refer](#configuration) |

### output

`output` is the path to the directory where the tool will be installed.

This can be set with (in order of priority):

- as a field in the tool definition

  ```yaml
  output: ~/.local/bin
  ```

- as a flag to the tool

  ```sh
  godyl --defaults.output ~/.local/bin
  ```

- as an environment variable

  ```sh
  GODYL_DEFAULTS_OUTPUT=~/.local/bin godyl
  ```

- in the configuration file

  ```yaml
  defaults:
    output: ~/.local/bin
  ```

| Field           | Type     | Template | Default            |
| --------------- | -------- | -------- | ------------------ |
| output          | string   | yes      | ~/.local/bin       |
| exe.patterns    | string[] | yes      | {{ .Exe.Name }}.\* |
| hints[].pattern | string   | yes      | {{ .Exe.Name }}    |
| hints[].weight  | number   | yes      | 1                  |
| source.type     | string   | no       | GitHub             |

| Field       | Type     | Template | Alt-form                                | Inferrence                             | Implemented |
| ----------- | -------- | -------- | --------------------------------------- | -------------------------------------- | ----------- |
| name        | string   | yes      | no                                      | no                                     | yes         |
| description | string   | yes      | no                                      | no                                     | yes         |
| version     | string   | yes      | no                                      | if left blank, from chosen source type | yes         |
| path        | string   | yes      | no                                      | if left blank, from chosen source type | yes         |
| checksum    | string   | no       | no                                      | no                                     | **no**      |
| output      | string   | yes      | no                                      | no                                     | yes         |
| aliases     | string[] | no       | `aliases: string -> aliases[0]: string` | no                                     | yes         |

| Field        | Type     | Template | Alt-form                                  | Inferrence                                        | Implemented |
| ------------ | -------- | -------- | ----------------------------------------- | ------------------------------------------------- | ----------- |
| exe          | dict     | no       | `exe: string -> exe.name: string`         | yes                                               | yes         |
| exe.name     | string   | yes      | no                                        | if left blank, from `name` if on form `<>/<name>` | yes         |
| exe.patterns | string[] | yes      | `patterns: string -> patterns[0]: string` | no                                                | yes         |

| Field                         | Type   | Template | Alt-form | Inferrence                                                              | Implemented |
| ----------------------------- | ------ | -------- | -------- | ----------------------------------------------------------------------- | ----------- |
| platform                      | dict   | yes      | no       | any attribute left blank under `platform.` will have its value inferred | yes         |
| platform.os                   | string | no       | no       | yes                                                                     | yes         |
| platform.architecture         | dict   | no       | no       | yes                                                                     | yes         |
| platform.architecture.type    | dict   | no       | no       | yes                                                                     | yes         |
| platform.architecture.version | dict   | no       | no       | yes                                                                     | yes         |
| platform.distribution         | string | no       | no       | yes                                                                     | yes         |
| platform.library              | string | no       | no       | yes                                                                     | yes         |
| platform.extension            | string | no       | no       | yes                                                                     | yes         |

| Field      | Type     | Template | Alt-form | Inferrence | Implemented |
| ---------- | -------- | -------- | -------- | ---------- | ----------- |
| values     | dict     | yes      | no       | no         | yes         |
| fallbacks  | string[] | no       | no       | no         | yes         |
| tags       | string[] | yes      | no       | no         | yes         |
| strategy   | string   | no       | no       | no         | yes         |
| extensions | string[] | no       | no       | no         | yes         |
| skip       | boolean  | no       | no       | no         | yes         |
| test       | string[] | yes      | no       | no         | yes         |

| Field         | Type    | Template | Alt-form | Inferrence | Implemented |
| ------------- | ------- | -------- | -------- | ---------- | ----------- |
| hints.pattern | string  | yes      | yes      | no         | yes         |
| hints.weight  | string  | yes      | yes      | no         | yes         |
| hints.regex   | boolean | no       | no       | no         | yes         |
| hints.must    | boolean | no       | no       | no         | yes         |

| Field       | Type   | Template | Alt-form | Inferrence | Implemented |
| ----------- | ------ | -------- | -------- | ---------- | ----------- |
| source.type | string | no       | no       | no         | yes         |

| Field               | Type   | Template | Alt-form | Inferrence | Implemented |
| ------------------- | ------ | -------- | -------- | ---------- | ----------- |
| source.github       | dict   | no       | no       | no         | yes         |
| source.github.owner | string | no       | no       | no         | yes         |
| source.github.repo  | string | no       | no       | no         | yes         |
| source.github.token | string | no       | no       | no         | yes         |

| Field            | Type   | Template | Alt-form | Inferrence | Implemented |
| ---------------- | ------ | -------- | -------- | ---------- | ----------- |
| source.url       | dict   | no       | no       | no         | yes         |
| source.url.url   | string | no       | no       | no         | yes         |
| source.url.token | string | no       | no       | no         | yes         |

| Field             | Type   | Template | Alt-form | Inferrence | Implemented |
| ----------------- | ------ | -------- | -------- | ---------- | ----------- |
| source.go         | dict   | no       | no       | no         | yes         |
| source.go.command | string | no       | no       | no         | yes         |
