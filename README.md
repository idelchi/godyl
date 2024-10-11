# godyl

Tool is work in progress and needs both cleaning up and documenting.

`godyl` helps with batch-downloading and installing statically compiled binaries from:

- GitHub releases
- Go projects
- URLs

As an alternative to above, custom commands can be used as well.

`godyl` will infer the platform and architecture from the system it is running on, and will attempt to download the appropriate binary.

This uses simple heuristics to infer the correct binary to download, and will not work for all projects.

Most properties can be overridden and `hints` can be used to help `godyl` make the correct decision.

## Installation

### From source

```sh
go install github.com/idelchi/godyl/cmd/godyl@latest
```

## From [installation script](https://raw.githubusercontent.com/idelchi/gocry/refs/heads/dev/scripts/install.sh)

```sh
curl -sSL https://raw.githubusercontent.com/idelchi/godyl/refs/heads/dev/scripts/install.sh | sh -s -- -v v0.0 -o ~/.local/bin
```

## Configuration

A configuration may be used to specify default settings for all tools. These will override (or extend in some case) the settings for each tool.

The following is embedded and used by default if no configuration is provided:

[config.yml](./cmd/godyl/config.yml)

```yaml
defaults:
  output: ~/.local/bin
  exe:
    patterns:
      - "{{ .Exe.Name }}.*"
  source:
    type: github
  hints:
    - pattern: "{{ .Exe.Name }}"
      weight: 1
```

The example above defines:

- The default output directory for all tools
- A pattern to use for when searching for the executable
- The default source to use if not specified
- A hint to use the executable name as a pattern (useful for repositories with multiple binaries, such as `ahmetb/kubectx`)

The full set of configuration options are:

```yaml
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
```

## Tools

A YAML file controls the tools to download and install. Alternative, if the positional argument to the tool is not a YAML file, it will be treated as a single tool name.

Examples are provided in [tools.yml](./tools.yml) and

```yaml
- ajeetdsouza/zoxide
```

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

| Field           | Type     | Template | Default              |
| --------------- | -------- | -------- | -------------------- |
| output          | string   | yes      | ~/.local/bin         |
| exe.patterns    | string[] | yes      | "{{ .Exe.Name }}.\*" |
| hints[].pattern | string   | yes      | "{{ .Exe.Name }}"    |
| hints[].weight  | number   | yes      | 1                    |
| source.type     | string   | no       | GitHub               |

| Field        | Type     | Template | Alt-form                                | Inferrence                                      | Implemented |
| ------------ | -------- | -------- | --------------------------------------- | ----------------------------------------------- | ----------- |
| name         | string   | yes      | -                                       | no                                              | yes         |
| description  | string   | yes      | no                                      | no                                              | yes         |
| version      | string   | yes      | no                                      | if left blank, from chosen source type          | yes         |
| path         | string   | yes      | no                                      | if left blank, from chosen source type          | yes         |
| checksum     | string   | no       | no                                      | no                                              | **no**      |
| output       | string   | yes      | no                                      | no                                              | yes         |
| exe          | dict     | no       | exe: string -> exe.name: string         | yes                                             | yes         |
| exe.name     | string   | yes      | no                                      | if left blank, from ´name`if on form`<>/<name>` | yes         |
| exe.patterns | string[] | yes      | patterns: string -> patterns[0]: string | no                                              | yes         |
| platform     | dict     | yes      | no                                      | yes                                             | yes         |
| platform     | dict     | yes      | no                                      | yes                                             | yes         |
| platform     | dict     | yes      | no                                      | yes                                             | yes         |
| platform     | dict     | yes      | no                                      | yes                                             | yes         |

| aliases | string[] | no | no | no | yes |
| values | dict | yes | no | no | yes |
| fallbacks | string[] | no | no | no | yes |
| hints.pattern | string | yes | no | no | yes |
| hints.weight | number | yes | no | no | yes |
| hints.regex | boolean | no | no | no | yes |
| hints.must | boolean | no | no | no | yes |
| source.type | string | no | yes | no | yes |
| source.github.owner | string | no | no | no | yes |
| source.github.repo | string | no | no | no | yes |
| source.github.token | string | no | no | no | yes |
| source.url.url | string | no | no | no | yes |
| source.url.token | string | no | no | no | yes |
| source.go.command | string | no | no | no | yes |
| tags | string[] | yes | no | no | yes |
| strategy | string | no | no | no | yes |
| extensions | string[] | no | no | no | yes |
| skip | boolean | no | no | no | yes |
| test | string[] | yes | no | no | yes |
