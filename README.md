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

> [!WARNING]
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

| Flag               | Environment Variable  | Default        | Description                                    |
| ------------------ | --------------------- | -------------- | ---------------------------------------------- |
| `--help`, `-h`     | `GODYL_HELP`          | `false`        | Show help message and exit                     |
| `--version`        | `GODYL_VERSION`       | `false`        | Show version information and exit              |
| `--dot-env`        | `GODYL_DOT_ENV`       | `.env`         | Path to .env file                              |
| `--defaults`, `-d` | `GODYL_DEFAULTS`      | `defaults.yml` | Path to defaults file                          |
| `--show-config`    | `GODYL_SHOW_CONFIG`   | `false`        | Show the parsed configuration and exit         |
| `--show-defaults`  | `GODYL_SHOW_DEFAULTS` | `false`        | Show the parsed default configuration and exit |
| `--show-env`       | `GODYL_SHOW_ENV`      | `false`        | Show the parsed environment variables and exit |
| `--show-platform`  | `GODYL_SHOW_PLATFORM` | `false`        | Detect the platform and exit                   |
| `--update`         | `GODYL_UPDATE`        | `false`        | Update the tools                               |
| `--dry`            | `GODYL_DRY`           | `false`        | Run without making any changes (dry run)       |
| `--log`            | `GODYL_LOG`           | `info`         | Log level (debug, info, warn, error)           |
| `--parallel`, `-j` | `GODYL_PARALLEL`      | `0`            | Number of parallel downloads (0 is unlimited)  |
| `--output`         | `GODYL_OUTPUT`        | `""`           | Output path for the downloaded tools           |
| `--tags`, `-t`     | `GODYL_TAGS`          | `["!native"]`  | Tags to filter tools by                        |
| `--source`         | `GODYL_SOURCE`        | `github`       | Source from which to install the tools         |
| `--strategy`       | `GODYL_STRATEGY`      | `none`         | Strategy to use for updating tools             |
| `--github-token`   | `GODYL_GITHUB_TOKEN`  | `""`           | GitHub token for authentication                |

The path to the tools file is provided as a positional argument, defaulting to `tools.yml`.

An example [tools.yml](./tools.yml) is provided.

## Tools

A YAML file controls the tools to download and install. Alternative, if the positional argument to the tool is not a YAML file, it will be treated as a single tool name.

Examples are provided in [tools.yml](./tools.yml) and

```yaml
- ajeetdsouza/zoxide
```

Above is the `simple` form to attempt to download the latest release of `zoxide` from `ajeetdsouza/zoxide`.

The full form is

```yaml
name: string
description: string
version: string
path: string
checksum: string
output: string
exe:
  name: string
  patterns:
    - string
platform:
  os: string
  architecture:
    type: string
    version: string
  library: string
  extension: string
  distribution: string
aliases: []
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
    repo: string
    owner: string
    token: string
    data:
      exe: string
      version: string
  url:
    url: string
    token: string
  go:
    command: string
  commands: []
tags:
  - string
strategy: string
extensions:
  - string
skip:
  - condition: string
    reason: string
test: []
allowFailure: boolean
post: []
mode: string
settings: {}
env:
  key: string
```

### Badges

![Inferred](https://img.shields.io/badge/Inferred-blue)
![Required](https://img.shields.io/badge/Required-red)
![Optional](https://img.shields.io/badge/Optional-green)
![Not Implemented](https://img.shields.io/badge/Not%20Implemented-gray)

### Name

![Required](https://img.shields.io/badge/Required-Yes-green)

| Template      | Templated | As Template |
| ------------- | --------- | ----------- |
| `{{ .Name }}` | No        | No          |

`name` is the name of the tool to download.

#### Usage

- Used to infer `exe.name` if not given
- Used to infer `sources.github.repo` and `sources.github.owner` if not given, by splitting on `/`

### Description

![Optional](https://img.shields.io/badge/Optional-Yes-blue)

| Template | Templated | As Template |
| -------- | --------- | ----------- |
| `N/A`    | No        | No          |

`description` is an optional description of the tool, for documentation purposes.

### Version

![Inferred](https://img.shields.io/badge/Inferred-blue)
![Optional](https://img.shields.io/badge/Optional-green)

| Template         | Templated | As Template |
| ---------------- | --------- | ----------- |
| `{{ .Version }}` | No        | Yes         |

`version` is the version of the tool to download.

#### Usage

- Will be inferred and populated by the `source` method if not given

### Path

![Inferred](https://img.shields.io/badge/Inferred-blue)
![Optional](https://img.shields.io/badge/Optional-green)

| Template | Templated | As Template |
| -------- | --------- | ----------- |
| `N/A`    | Yes       | No          |

`path` is the path to the tool to download. Currently only supports URLs.

#### Usage

- Will be inferred and populated by the `source` method if not given

### Checksum

![Not Implemented](https://img.shields.io/badge/Not%20Implemented-gray)

### Output

![Required](https://img.shields.io/badge/Required-red)

| Template | Templated | As Template |
| -------- | --------- | ----------- |
| `N/A`    | No        | No          |

`output` is the path to the directory where the tool will be installed.

#### Usage

- Set according to [flags and environment variables](#configuration) or [defaults](#defaults) if not given

### Exe

![Inferred](https://img.shields.io/badge/Inferred-blue)
![Optional](https://img.shields.io/badge/Optional-green)

| Template              | Templated | As Template |
| --------------------- | --------- | ----------- |
| `{{ .Exe.Name }}`     | Yes       | Yes         |
| `{{ .Exe.Patterns }}` | Yes       | No          |

`exe` is a dictionary containing the name of the executable and patterns to use for finding the executable.

#### Usage

- `exe.name` is the name of the executable, inferred from `name` if not given
- `exe.patterns` is a list of patterns to use for finding the executable, inferred from `name` if not given
- Set according to [defaults](#defaults) if not given

### Platform

![Inferred](https://img.shields.io/badge/Inferred-blue)
![Optional](https://img.shields.io/badge/Optional-green)

| Template             | Templated | As Template |
| -------------------- | --------- | ----------- |
| `{{ .Platform.<> }}` | No        | Yes         |

`platform` is a dictionary containing the platform and architecture information.

#### Usage

- Any field not given will be inferred from the system
- Set according to [defaults](#defaults) if not given

### Aliases

![Optional](https://img.shields.io/badge/Optional-green)

| Template | Templated | As Template |
| -------- | --------- | ----------- |
| `N/A`    | No        | No          |

`aliases` is a list of aliases for the tool. Will be used to create symlinks (or copies on `Windows`) for the tool.

### Values

![Optional](https://img.shields.io/badge/Optional-green)

| Template           | Templated | As Template |
| ------------------ | --------- | ----------- |
| `{{ .Values.<> }}` | No        | Yes         |

`values` is an arbitrary values map, which can be used for templating in other fields.

### Fallbacks

![Optional](https://img.shields.io/badge/Optional-green)

| Template | Templated | As Template |
| -------- | --------- | ----------- |
| `N/A`    | No        | No          |

`fallbacks` is a list of fallback strategies to use if the tool cannot be found.

They will be tried in order until the tool is found or all have been tried.

### Hints

![Optional](https://img.shields.io/badge/Optional-green)

| Template               | Templated | As Template |
| ---------------------- | --------- | ----------- |
| `{{ .Hints.Weight }}`  | Yes       | No          |
| `{{ .Hints.Pattern }}` | Yes       | No          |
| `{{ .Hints.Regex }}`   | No        | No          |
| `{{ .Hints.Must }}`    | No        | No          |

`hints` is a list of hints for matching, which can be used to help `godyl` find the correct tool.

#### Usage

- Set according to [defaults](#defaults) if not given

### Source

#### Type

![Required](https://img.shields.io/badge/Required-red)

| Template | Templated | As Template |
| -------- | --------- | ----------- |
| `N/A`    | No        | No          |

`source.type` is the source type. Accepted values are:

- `github`
- `url`
- `go`
- `commands`

#### GitHub

![Inferred](https://img.shields.io/badge/Inferred-blue)
![Optional](https://img.shields.io/badge/Optional-green)

| Template | Templated | As Template |
| -------- | --------- | ----------- |
| `N/A`    | No        | No          |

`source.github` is a dictionary containing the owner, repository and token of the tool.

#### Usage

- `repo` and `owner` will be inferred from `name` if not given
- `token` will be set according to [flags and environment variables](#configuration) if not given

#### URL

![Optional](https://img.shields.io/badge/Optional-green)

| Template | Templated | As Template |
| -------- | --------- | ----------- |
| `N/A`    | No        | No          |

`source.url` is a dictionary containing the URL and token of the tool.

#### Usage

- `token` will be set according to [flags and environment variables](#configuration) if not given

#### Go

![Inferred](https://img.shields.io/badge/Inferred-blue)
![Optional](https://img.shields.io/badge/Optional-green)

| Template | Templated | As Template |
| -------- | --------- | ----------- |
| `N/A`    | No        | No          |

`source.go` is a dictionary containing the command to run to install the tool.

## Defaults

A default configuration may be used to specify default settings for all tools. These will override (or extend in some case) the settings for each tool.

The following is embedded and used by default if no configuration is provided:

[config.yml](./cmd/godyl/defaults.yml)

The example above defines:

- The default output directory for all tools (`~/.local/bin`)
- Patterns to use for when searching for the executable (`"{{ .Exe.Name }}{{ .Platform.Extension }}$"`)
- Hints to:

  - use the executable name as a pattern (useful for repositories with multiple binaries, such as `ahmetb/kubectx`)

    ```yaml
    pattern: "{{ .Exe.Name }}"
    weight: 1
    ```

  - prefer `.zip` files for Windows

    ```yaml
    pattern: zip
    weight: '{{ if eq .Platform.OS "windows" }}1{{ else }}0{{ end }}'
    ```

- `find` mode for downloading, extracting and finding the executable
- The default source type as `github`
- `none` strategy to skip tools which already exist
- Settings the environment variable `GH_TOKEN` to the value of `GODYL_GITHUB_TOKEN`

The full set of default options are:

```yaml
exe:
  name: string
  patterns:
    - string
output: string
platform:
  os: string
  architecture:
    type: string
    version: string
  library: string
  extension: string
  distribution: string
values:
  key: any
fallbacks:
  - string
hints:
  - pattern: string
    weight: int
    regex: boolean
    must: boolean
source:
  type: string
  github:
    repo: string
    owner: string
    token: string
  url:
    url: string
    token: string
  go:
    command: string
  commands:
    - string
tags:
  - string
strategy: string
extensions:
  - string
env:
  key: string
mode: string
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
