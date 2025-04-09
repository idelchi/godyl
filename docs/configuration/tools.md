---
layout: default
title: Tools Format
---

# Tools Format

Tools can be defined in a YAML file (typically `tools.yml`).

See [tools.yml](https://github.com/idelchi/godyl/blob/main/tools.yml) for getting started.

You can use a simple form or a full form for tool definitions.

### Simple Form

```yaml
- idelchi/godyl
```

This is the simplest form to download the latest release of `godyl` from the the default source type (`github`).

If the path is a URL, it will be considered as a `source.url` type. Otherwise, it will be assumed to be a `source.github` type in the form `owner/repo`.

### Full Form

For more complex configurations, you can use the extended form:

{% raw  %}

```yaml
name: godyl
description: Asset downloader
version:
  version: v0.1.0
path: https://github.com/idelchi/godyl/releases/download/v0.1.0/godyl_linux_amd64.tar.gz
output: ~/.local/bin
exe:
  name: godyl
  patterns:
    - "{{ .Exe }}{{ .EXTENSION }}$"
platform:
  os: linux
  architecture:
    type: amd64
aliases:
  - gd
source:
  type: github
tags:
  - cli
  - downloader
strategy: upgrade
```

A complete reference for all fields is available below.

### Full form

Below are all configuration options along with examples

```yaml
name: idelchi/godyl
description: Asset downloader for GitHub releases, URLs, and Go projects
version:
  version: v0.1.0 # For `github` and `gitlab` sources, leave empty to get the latest release
  commands:
    - --version
    - version
  patterns:
    - '.*?(\d+\.\d+\.\d+).*'
path: "https://github.com/idelchi/godyl/releases/download/v0.0.1/godyl_{{ .OS }}_{{ .ARCH }}.tar.gz" # For `github` and `gitlab` sources, leave empty to get the latest release
output: ~/.local/bin
exe:
  name: godyl
  patterns: "^{{ .OS }}-{{ .Exe }}-{{ .ARCH }}{{ .EXTENSION }}$"
platform:
  os: windows
  architecture:
    type: arm
    version: 7
  library: gnu
  extension: exe
  distribution: windows
aliases:
  - gd
values:
  customField: customValue
fallbacks:
  - go
hints:
  - pattern: amd64
    weight: |-
      {{- if eq .ARCH "amd64" -}}
      1
      {{- else -}}
      0
      {{- end -}}
  - pattern: "{{ .Exe }}"
    must: true
source:
  type: github|gitlab|url|go|none
  github:
    repo: godyl
    owner: idelchi
    token: secret
  gitlab:
    project: godyl
    namespace: idelchi/go-projects
    token: secret
    server: https://gitlab.self-hosted.com
  url:
    token:
      token: secret
      header: Authorization
      scheme: Bearer
      headers:
        Content-Type:
          - application/json
          - application/x-www-form-urlencoded
        Accept:
          - application/json
          - application/x-www-form-urlencoded
  go:
    command: cmd/godyl
commands:
  pre:
    commands:
      - "mkdir -p {{ .Output }}"
    allow_failure: true
    exit_on_error: false
  post:
    commands:
      - "chmod +x {{ .Output }}/{{ .Exe }}"
    allow_failure: false
    exit_on_error: true
tags:
  - downloader
strategy: none|upgrade|force
extensions:
  - .exe
skip:
  - reason: "godyl is not available for Darwin"
    condition: '{{ eq .OS "darwin" }}'
mode: find|extract
env:
  GH_TOKEN: $GODYL_GITHUB_TOKEN
no_verify_ssl: false
```

Most of the fields also support simplified forms which will be described below.

## Templating

Many fields in the configuration support templating with variables like:

- `{{ .Name }}` - The name of the tool
- `{{ .Output }}` - The output path
- `{{ .Exe }}` - The executable name
- `{{ .OS }}` - The operating system
- `{{ .ARCH }}` - The architecture
- `{{ .EXTENSION }}` - The file extension for the platform

For example, to set a path that adapts to the current platform:

```yaml
path: https://example.com/download/{{ .Name }}_{{ .OS }}_{{ .ARCH }}.tar.gz
```

## Available Fields

Below is a comprehensive list of fields that can be used to configure each tool:

### Name

**Required**: Yes

The name of the tool to download. Used as display name and for inferring other fields.

```yaml
name: idelchi/godyl
```

### Description

**Optional**: Yes

A description of the tool, for documentation purposes.

```yaml
description: Asset downloader for GitHub releases, URLs, and Go projects
```

### Version

**Optional**: Yes (will be inferred if not provided)

The version of the tool to download.

Simple form:

```yaml
version: v0.1.0
```

Full form:

```yaml
version:
  version: v0.1.0
  commands:
    - "--version"
  patterns:
    - "v\\d+\\.\\d+\\.\\d+"
```

- `version.version`: The version to download
- `version.commands`: Commands to run to get the installed version (for upgrades)
- `version.patterns`: Regex patterns to extract the version from command output (for upgrades)

### Path

**Optional**: Yes (will be inferred if not provided)

The path to the tool to download. Must be a URL to a file.

```yaml
path: https://github.com/idelchi/godyl/releases/download/v0.1.0/godyl_linux_amd64.tar.gz
```

The most common use-case is to have it inferred from the `source` field configuration.

### Output

**Required** (can be set from defaults or flags)

The directory where the tool will be installed.

```yaml
output: ~/.local/bin
```

### Exe

**Optional**: Yes (will be inferred if not provided)

Information about the executable.

Simple form:

```yaml
exe: godyl
```

Full form:

```yaml
exe:
  name: godyl
  patterns:
    - "^{{ .Exe }}{{ .EXTENSION }}$"
    - ".*/{{ .Exe }}{{ .EXTENSION }}$"
```

- `exe.name`: The name of the executable
- `exe.patterns`: Regex patterns to find the executable in the downloaded archive

### Platform

**Optional**: Yes (will be inferred from the system)

Platform and architecture information.

```yaml
platform:
  os: linux
  architecture:
    type: amd64
    version: ""
  library: gnu
  extension: ""
  distribution: debian
```

### Aliases

**Optional**: Yes

Aliases for the tool. Will create symlinks (or copies on Windows).

```yaml
aliases:
  - gd
  - godl
```

### Values

**Optional**: Yes

Arbitrary values that can be used in templates.

```yaml
values:
  customField: customValue
```

### Fallbacks

**Optional**: Yes

Fallback strategies if no matches were made in releases.

```yaml
fallbacks:
  - go
```

### Hints

**Optional**: Yes

Hints to help Godyl find the correct tool.

```yaml
hints:
  - pattern: "{{ .Exe }}"
    weight: 1
  - pattern: "{{ .OS }}"
    must: true
```

### Source

**Required** (can be set from defaults)

Information about the source of the tool.

GitHub source:

```yaml
source:
  type: github
  github:
    repo: godyl
    owner: idelchi
    token:
```

URL source:

```yaml
source:
  type: url
  url:
    token:
      token:
      header: Authorization
      scheme: Bearer
```

Go source:

```yaml
source:
  type: go
  go:
    command: cmd/godyl
```

### Commands

**Optional**: Yes

Commands to run before and after installation.

```yaml
commands:
  pre:
    - "mkdir -p {{ .Output }}"
  post:
    - "chmod +x {{ .Output }}/{{ .Exe }}"
```

### Tags

**Optional**: Yes

Tags to filter tools.

```yaml
tags:
  - cli
  - downloader
```

### Strategy

**Required** (can be set from defaults)

Strategy for updating the tool.

```yaml
strategy: upgrade
```

Valid values:

- `none`: Skip if the tool already exists
- `upgrade`: Upgrade if a newer version is available
- `force`: Always download and install

### Extensions

**Optional**: Yes

Extensions to consider when matching tools.

```yaml
extensions:
  - .zip
  - .tar.gz
```

### Skip

**Optional**: Yes

Conditions under which to skip the tool.

```yaml
skip:
  - condition: '{{ eq .OS "windows" }}'
    reason: "Tool is not available on Windows"
```

### Mode

**Required** (can be set from defaults)

Mode for downloading and installing.

```yaml
mode: find
```

Valid values:

- `find`: Download, extract, and find the executable
- `extract`: Download and extract directly to the output directory

{% endraw  %}
