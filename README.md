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

- flags
- environment variables
- `.env` file

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

The path to the file containing the tool installation instructions is provided as a positional argument, defaulting to `tools.yml`.

An example [tools.yml](./tools.yml) is provided.

In general, settings can be set in the following ways (order of priority):

- as a field in the [tools.yml](./tools.yml) definition

  ```yaml
  output: ~/.local/bin
  ```

- as a flag to the tool

  ```sh
  godyl --output ~/.local/bin
  ```

- as an environment variable

  ```sh
  GODYL_OUTPUT=~/.local/bin godyl
  ```

- in an `.env` file

  ```
  GODYL_OUTPUT=~/.local/bin
  ```

- by setting the value in a `defaults.yml` file (see [defaults](#defaults))

  ```yaml
  output: ~/.local/bin
  ```

If none of the above are fulfilled, the default configuration embedded from [defaults.yml](./cmd/godyl/defaults.yml) will be used.

## Tools

A YAML file controls the tools to download and install. Alternatively, if the positional argument to the tool is not a YAML file, it will be treated as a single tool name or URL.

Examples are provided in [tools.yml](./tools.yml).

### Simple form

```yaml
- ajeetdsouza/zoxide
```

Above is the `simple` form to attempt to download the latest release of `zoxide` from `ajeetdsouza/zoxide`.

If it is a simply two-part string, it will be considered as a `source.github` type.
If it is a URL, it will be considered as a `source.url` type.

### Full form

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

Any field that accepts a list, can also be provided as a string.

For example:

```yaml
aliases: z
fallbacks: go
exe:
  patterns: "{{ .Exe.Name }}{{ .Platform.Extension }}$"
```

is equivalent to:

```yaml
aliases:
  - z
fallbacks:
  - go
exe:
  patterns:
    - "{{ .Exe.Name }}{{ .Platform.Extension }}$"
```

### Name

![Required](https://img.shields.io/badge/Required-Yes-green)

| Template      | Templated | As Template |
| ------------- | --------- | ----------- |
| `{{ .Name }}` | ![no]     | ![yes]      |

`name` is the name of the tool to download.

#### Usage

- Used to infer `exe.name` if not given
- Used to infer `sources.github.repo` and `sources.github.owner` if not given, by splitting on `/`

### Description

![Optional](https://img.shields.io/badge/Optional-Yes-blue)

| Template | Templated | As Template |
| -------- | --------- | ----------- |
| ![na]    | ![no]     | ![no]       |

`description` is an optional description of the tool, for documentation purposes.

### Version

![Inferred](https://img.shields.io/badge/Inferred-blue)
![Optional](https://img.shields.io/badge/Optional-green)

| Template         | Templated | As Template |
| ---------------- | --------- | ----------- |
| `{{ .Version }}` | ![no]     | ![yes]      |

[version]: https://img.shields.io/badge/Version-{{ .Version }}-blue
`version` is the version of the tool to download.

#### Usage

- Will be inferred and populated by the `source` method if not given

### Path

![Inferred](https://img.shields.io/badge/Inferred-blue)
![Optional](https://img.shields.io/badge/Optional-green)

| Template | Templated | As Template |
| -------- | --------- | ----------- |
| ![na]    | ![yes]    | ![no]       |

`path` is the path to the tool to download. Currently only supports URLs.

#### Usage

- Will be inferred and populated by the `source` method if not given

### Checksum

![Not Implemented](https://img.shields.io/badge/Not%20Implemented-gray)

### Output

![Required](https://img.shields.io/badge/Required-red)

| Template | Templated | As Template |
| -------- | --------- | ----------- |
| ![na]    | ![no]     | ![no]       |

`output` is the path to the directory where the tool will be installed.

#### Usage

- Set according to [flags and environment variables](#configuration) or [defaults](#defaults) if not given

### Exe

![Inferred](https://img.shields.io/badge/Inferred-blue)
![Optional](https://img.shields.io/badge/Optional-green)

| Template              | Templated | As Template |
| --------------------- | --------- | ----------- |
| `{{ .Exe.Name }}`     | ![yes]    | ![yes]      |
| `{{ .Exe.Patterns }}` | ![yes]    | ![no]       |

`exe` is a dictionary containing the name of the executable and patterns to use for finding the executable.

#### Usage

- `exe.name` is the name of the executable, inferred from `name` if not given
- `exe.patterns` is a list of patterns to use for finding the executable, inferred from `name` if not given
- Set according to [defaults](#defaults) if not given

#### Alternative form

```yaml
exe: zoxide
```

is equivalent to:

```yaml
exe:
  name: zoxide
```

### Platform

![Inferred](https://img.shields.io/badge/Inferred-blue)
![Optional](https://img.shields.io/badge/Optional-green)

| Template             | Templated | As Template |
| -------------------- | --------- | ----------- |
| `{{ .Platform.<> }}` | ![no]     | ![yes]      |

`platform` is a dictionary containing the platform and architecture information.

#### Usage

- Any field not given will be inferred from the system
- Set according to [defaults](#defaults) if not given

### Aliases

![Optional](https://img.shields.io/badge/Optional-green)

| Template | Templated | As Template |
| -------- | --------- | ----------- |
| ![na]    | ![no]     | ![no]       |

`aliases` is a list of aliases for the tool. Will be used to create symlinks (or copies on `Windows`) for the tool.

### Values

![Optional](https://img.shields.io/badge/Optional-green)

| Template           | Templated | As Template |
| ------------------ | --------- | ----------- |
| `{{ .Values.<> }}` | ![no]     | ![yes]      |

`values` is an arbitrary values map, which can be used for templating in other fields.

### Fallbacks

![Optional](https://img.shields.io/badge/Optional-green)

| Template | Templated | As Template |
| -------- | --------- | ----------- |
| ![na]    | ![no]     | ![no]       |

`fallbacks` is a list of fallback strategies to use if the tool cannot be found.

They will be tried in order until the tool is found or all have been tried.

### Hints

![Optional](https://img.shields.io/badge/Optional-green)

| Template               | Templated | As Template |
| ---------------------- | --------- | ----------- |
| `{{ .Hints.Weight }}`  | ![yes]    | ![no]       |
| `{{ .Hints.Pattern }}` | ![yes]    | ![no]       |
| `{{ .Hints.Regex }}`   | ![no]     | ![no]       |
| `{{ .Hints.Must }}`    | ![no]     | ![no]       |

`hints` is a list of hints for matching, which can be used to help `godyl` find the correct tool.

#### Usage

- Set according to [defaults](#defaults) if not given

### Source

#### Type

![Required](https://img.shields.io/badge/Required-red)

| Template | Templated | As Template |
| -------- | --------- | ----------- |
| ![na]    | ![no]     | ![no]       |

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
| ![na]    | ![no]     | ![no]       |

`source.github` is a dictionary containing the owner, repository and token of the tool.

#### Usage

- `repo` and `owner` will be inferred from `name` if not given, or set according to [defaults](#defaults) (not recommended)
- `token` will be set according to [flags and environment variables](#configuration) or [defaults](#defaults) if not given

#### URL

![Optional](https://img.shields.io/badge/Optional-green)

| Template | Templated | As Template |
| -------- | --------- | ----------- |
| ![na]    | ![no]     | ![no]       |

`source.url` is a dictionary containing the URL and token of the tool.

#### Usage

- `token` will be set according to [flags and environment variables](#configuration) if not given

#### Go

![Inferred](https://img.shields.io/badge/Inferred-blue)
![Optional](https://img.shields.io/badge/Optional-green)

| Template | Templated | As Template |
| -------- | --------- | ----------- |
| ![na]    | ![no]     | ![no]       |

`source.go` is a dictionary containing the command to run to install the tool.

#### Commands

![Optional](https://img.shields.io/badge/Optional-green)

| Template | Templated | As Template |
| -------- | --------- | ----------- |
| ![na]    | ![yes]    | ![no]       |

`commands` is a list of commands to run to install the tool.

### Tags

![Inferred](https://img.shields.io/badge/Inferred-blue)
![Optional](https://img.shields.io/badge/Optional-green)

| Template | Templated | As Template |
| -------- | --------- | ----------- |
| ![na]    | ![no]     | ![no]       |

`tags` is a list of tags to filter tools by.

#### Usage

- the `name` of the tool will always be added as a tag

### Strategy

![Required](https://img.shields.io/badge/Required-red)

| Template | Templated | As Template |
| -------- | --------- | ----------- |
| ![na]    | ![no]     | ![no]       |

`strategy` is the strategy to use for updating the tool.

Accepted values are:

- `none`
- `upgrade`
- `force`

#### Usage

- Set according to [flags and environment variables](#configuration) or [defaults](#defaults) if not given
- `none` will skip the tool if it already exists
- `upgrade` will attempt to parse the version of an existing tool and upgrade if necessary
- `force` will always download and install the tool

### Extensions

![Optional](https://img.shields.io/badge/Optional-green)

| Template | Templated | As Template |
| -------- | --------- | ----------- |
| ![na]    | ![no]     | ![no]       |

`extensions` is a list of extensions to add to the tool name when matching the tools (used only for `github` source type).

#### Usage

- Set according to [defaults](#defaults) if not given
- Can be used to for example prefer `.zip` files for Windows

### Skip

![Optional](https://img.shields.io/badge/Optional-green)

| Template | Templated | As Template |
| -------- | --------- | ----------- |
| ![na]    | ![yes]    | ![no]       |

`skip` is a list of conditions under which the tool should be skipped.

#### Usage

- `reason` is a description of why the tool was skipped
- `condition` is a condition to check, and can use templating

#### Alternative form

```yaml
skip: <condition>
```

is equivalent to:

```yaml
skip:
  - condition: <condition>
```

### Test

![Not Implemented](https://img.shields.io/badge/Not%20Implemented-gray)

### Allow Failure

![Not Implemented](https://img.shields.io/badge/Not%20Implemented-gray)

### Post

![Optional](https://img.shields.io/badge/Optional-green)

| Template | Templated | As Template |
| -------- | --------- | ----------- |
| ![na]    | ![yes]    | ![no]       |

`post` is a list of commands to run after the tool has been installed.

### Mode

![Required](https://img.shields.io/badge/Required-red)

| Template | Templated | As Template |
| -------- | --------- | ----------- |
| ![na]    | ![no]     | ![no]       |

`mode` is the mode to use for the tool.

Accepted values are:

- `find`
- `download`

#### Usage

- `find` will download, extract and find the executable
- `download` will download the tool and extract it directly to the output directory
- Set according to [flags and environment variables](#configuration) or [defaults](#defaults) if not given

## Defaults

A default configuration may be used to specify default settings for all tools. These will override (or extend in some case) the settings for each tool.

The following is embedded and used by default if no default configuration is provided:

[defaults.yml](./cmd/godyl/defaults.yml)

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
exe: {}
output: str
platform: {}
values: {}
fallbacks: []
hints: {}
source: {}
tags: []
strategy: str
extensions: []
env: {}
mode: str
```

<!-- Badges -->

[yes]: https://img.shields.io/badge/Yes-green
[no]: https://img.shields.io/badge/No-red
[inferred]: https://img.shields.io/badge/Inferred-blue
[required]: https://img.shields.io/badge/Required-red
[optional]: https://img.shields.io/badge/Optional-green
[not-implemented]: https://img.shields.io/badge/Not%20Implemented-gray
[na]: https://img.shields.io/badge/N%2FA-lightgrey

<!-- Badges -->
