---
layout: default
title: Tools Format
---

# Tools YAML Format

The `tools.yml` file is used to define the tools that Godyl should download and install. This page describes the format of this file in detail.

## Basic Structure

A `tools.yml` file contains a list of tools to download:

```yaml
- name: tool1
  # tool1 configuration...
- name: tool2
  # tool2 configuration...
```

## Available Fields

Below is a comprehensive list of fields that can be used to configure each tool:

### Name

**Required**: Yes

The name of the tool to download. Used as display name and for inferring other fields.

```yaml
name: godyl
```

### Description

**Optional**: Yes

A description of the tool, for documentation purposes.

```yaml
description: Asset downloader for GitHub releases, URLs, and Go projects
```

### Version

**Optional**: Yes (will be inferred if not provided)

The version of the tool to download. Can be a simple string or a complex object with commands for version detection.

Simple form:

```yaml
version: v0.1.0
```

Full form:

```yaml
version:
  version: v0.1.0
  commands:
    - "{{ .Exe }} --version"
  patterns:
    - "v\\d+\\.\\d+\\.\\d+"
```

- `version.version`: The version to download
- `version.commands`: Commands to run to get the installed version (for upgrades)
- `version.patterns`: Regex patterns to extract the version from command output

### Path

**Optional**: Yes (will be inferred if not provided)

The path to the tool to download. Currently only supports URLs.

```yaml
path: https://github.com/idelchi/godyl/releases/download/v0.1.0/godyl_linux_amd64.tar.gz
```

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

Fallback strategies if the tool cannot be found.

```yaml
fallbacks:
  - go
  - python
```

### Hints

**Optional**: Yes

Hints to help Godyl find the correct tool.

```yaml
hints:
  - pattern: "{{ .Exe }}"
    weight: 1
  - pattern: "{{ .OS }}"
    weight: 2
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
    token: ${GITHUB_TOKEN}
```

URL source:

```yaml
source:
  type: url
  url:
    token:
      token: ${URL_TOKEN}
      header: Authorization
```

Go source:

```yaml
source:
  type: go
  go:
    command: ./cmd/godyl
```

Commands source:

```yaml
source:
  type: commands
  commands:
    - "pip install godyl=={{ .Version }}"
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
