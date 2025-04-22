---
layout: default
title: Template Reference
---

# Template Reference

`godyl` uses Go templates for various configuration fields. This page provides a reference for the templating system in `godyl`.

## Template Variables

The following variables are available in templates:

{% raw  %}

| Variable              | Description                                | Example Value                                |
| --------------------- | ------------------------------------------ | -------------------------------------------- |
| `{{ .Name }}`         | The name of the tool                       | `idelchi/godyl`                              |
| `{{ .Output }}`       | The output path                            | `~/.local/bin`                               |
| `{{ .Exe }}`          | The name of the executable                 | `godyl`                                      |
| `{{ .Env.<> }}`       | Any environment variable                   | `{{ .Env.HOME }}` -> `/home/user`            |
| `{{ .Values.<> }}`    | Custom values for templating               | `{{ .Values.customField }}` -> `customValue` |
| `{{ .Version }}`      | The version of the tool                    | `v0.1.0`                                     |
| `{{ .OS }}`           | The operating system                       | `linux`, `darwin`, `windows`                 |
| `{{ .ARCH }}`         | The architecture type                      | `amd64`, `arm64`                             |
| `{{ .ARCH_VERSION }}` | The version of the architecture            | `7` (for ARM v7)                             |
| `{{ .ARCH_LONG }}`    | The long version of the architecture       | `armv7l`                                     |
| `{{ .IS_ARM }}`       | Whether the architecture is ARM (32 or 64) | `true` or `false`                            |
| `{{ .IS_X86 }}`       | Whether the architecture is x86 (32 or 64) | `true` or `false`                            |
| `{{ .LIBRARY }}`      | The system library                         | `gnu`, `musl`                                |
| `{{ .EXTENSION }}`    | The file extension for the platform        | `.exe` on Windows, empty on Linux/MacOS      |
| `{{ .DISTRIBUTION }}` | The distribution name                      | `debian`, `alpine`                           |

## Template Functions

`godyl` includes all functions from the [slim-sprig](https://github.com/go-task/slim-sprig) library.

## Templated Fields

The following fields in the tool configuration support templating:

| Field               | Example                                                                |
| ------------------- | ---------------------------------------------------------------------- |
| `output`            | `output: bin/{{ .OS }}-{{ .ARCH }}`                                    |
| `skip.condition`    | `condition: '{{ eq .OS "windows" }}'`                                  |
| `version`           | `version: v{{ if eq .OS "windows" }}0.1.0{{ else }}0.2.0{{ end }}`     |
| `source.type`       | `type: {{ if eq .OS "windows" }}url{{ else }}github{{ end }}`          |
| `exe.patterns`      | `patterns: ["^{{ .OS }}-{{ .Exe}}"]`                                   |
| `extensions`        | `extensions: ['{{ if eq .OS "windows" }}.exe{{ else }}{{ end }}']`     |
| `commands.pre/post` | `["pip install {{ .Exe }}=={{ .Version }}"]`                           |
| `hints[].pattern`   | `pattern: "{{ .OS }}"`                                                 |
| `hints[].weight`    | `weight: '{{ if eq .ARCH "arm" }}1{{ else }}0{{ end }}'`               |
| `path`              | `path: "https://example.com/{{ .Name }}_{{ .OS }}_{{ .ARCH }}.tar.gz"` |

## Conditional Logic

You can use conditional logic in templates:

```yaml
version: |-
  {{- if eq .OS "windows" -}}
    v0.1.0
  {{- else if eq .OS "darwin" -}}
    v0.2.0
  {{- else -}}
    v0.3.0
  {{- end -}}
```

## Examples

### Platform-specific Paths

```yaml
path: https://releases.hashicorp.com/terraform/{{ .Version | trimPrefix "v" }}/terraform_{{ .Version | trimPrefix "v" }}_{{ .OS }}_{{ .ARCH }}.zip
```

### Skip Windows Installation

```yaml
skip:
  condition: '{{ eq .OS "windows" }}'
  reason: "Tool is not available on Windows"
```

{% endraw  %}
