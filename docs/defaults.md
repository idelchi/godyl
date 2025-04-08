---
layout: default
title: Default Configuration
---

# Default Configuration

Godyl uses a default configuration to provide sensible defaults for all tools. This can be overridden by providing your own `defaults.yml` file.

## Default Configuration File

The default configuration is embedded in the Godyl binary and looks like this:

```yaml
output: ~/.local/bin

exe:
  patterns:
    - "^{{ .Exe }}{{ .EXTENSION }}$"
    - ".*/{{ .Exe }}{{ .EXTENSION }}$"

hints:
  - pattern: "{{ .Exe }}"
    weight: 1
  - pattern: zip
    weight: '{{ if eq .OS "windows" }}1{{ else }}0{{ end }}'
  - pattern: "armv{{.ARCH_VERSION}}"
    weight: 3
  - pattern: "armv{{sub .ARCH_VERSION 1}}"
    weight: 2
  - pattern: "armv{{sub .ARCH_VERSION 2}}"
    weight: 1

extensions:
  - '{{ if eq .OS "windows" }}.exe{{ else }}{{ end }}'
  - '{{ if has .OS (list "windows" "darwin") }}.zip{{ else }}{{ end }}'
  - .tar.gz

mode: find
source:
  type: github
strategy: none
env:
  GH_TOKEN: '{{ .Env.GODYL_GITHUB_TOKEN }}'
```

## Overriding Defaults

You can override these defaults by:

1. Creating your own `defaults.yml` file
2. Passing it to Godyl with the `--defaults` flag:

```sh
godyl --defaults my-defaults.yml install tools.yml
```

3. Or setting the `GODYL_DEFAULTS` environment variable:

```sh
GODYL_DEFAULTS=my-defaults.yml godyl install tools.yml
```

## Available Default Fields

The following fields can be set in the `defaults.yml` file:

```yaml
exe:       # Default executable settings
output:    # Default output directory
platform:  # Default platform settings
values:    # Default values for templating
fallbacks: # Default fallback strategies
hints:     # Default hints for matching
source:    # Default source settings
tags:      # Default tags
strategy:  # Default update strategy
env:       # Default environment variables
mode:      # Default mode (find or extract)
```

## Examples

### Custom Output Directory

```yaml
output: /usr/local/bin
```

### Alternative GitHub Token Environment Variable

```yaml
env:
  GH_TOKEN: '{{ .Env.GITHUB_TOKEN }}'
```

### Different Default Source

```yaml
source:
  type: gitlab
```

### Automatic Upgrades

```yaml
strategy: upgrade
```

### Custom Executable Patterns

```yaml
exe:
  patterns:
    - "^bin/{{ .OS }}/{{ .Exe }}{{ .EXTENSION }}$"
    - "^{{ .Exe }}_{{ .OS }}_{{ .ARCH }}{{ .EXTENSION }}$"
```

## Viewing Default Configuration

You can view the current default configuration with:

```sh
godyl dump defaults
```

## Precedence

The precedence for configuration settings, from highest to lowest, is:

1. Command-line flags
2. Environment variables
3. `.env` file(s)
4. Tool-specific configuration in `tools.yml`
5. Custom `defaults.yml` file
6. Embedded default configuration