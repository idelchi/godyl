---
layout: default
title: Template Reference
---

# Template Reference

Godyl uses Go templates for various configuration fields. This page provides a reference for the templating system in Godyl.

## Template Variables

The following variables are available in templates:

| Variable              | Description                         | Example Value                                |
| --------------------- | ----------------------------------- | -------------------------------------------- |
| `{{ .Name }}`         | The name of the tool                | `godyl`                                      |
| `{{ .Output }}`       | The output path                     | `~/.local/bin`                               |
| `{{ .Exe }}`          | The name of the executable          | `godyl`                                      |
| `{{ .Env.<> }}`       | Any environment variable            | `{{ .Env.HOME }}` -> `/home/user`            |
| `{{ .Values.<> }}`    | Custom values for templating        | `{{ .Values.customField }}` -> `customValue` |
| `{{ .Version }}`      | The version of the tool             | `v0.1.0`                                     |
| `{{ .OS }}`           | The operating system                | `linux`, `darwin`, `windows`                 |
| `{{ .ARCH }}`         | The architecture type               | `amd64`, `arm64`                             |
| `{{ .ARCH_VERSION }}` | The version of the architecture     | `7` (for ARM v7)                             |
| `{{ .LIBRARY }}`      | The system library                  | `gnu`, `musl`                                |
| `{{ .EXTENSION }}`    | The file extension for the platform | `.exe` on Windows, empty on Linux/MacOS      |
| `{{ .DISTRIBUTION }}` | The distribution name               | `debian`, `alpine`                           |

## Template Functions

Godyl includes all functions from the [slim-sprig](https://github.com/go-task/slim-sprig) library. Here are some commonly used functions:

### String Functions

| Function     | Description                        | Example                                              |
| ------------ | ---------------------------------- | ---------------------------------------------------- |
| `lower`      | Convert to lowercase               | `{{ lower "HELLO" }}` -> `hello`                     |
| `upper`      | Convert to uppercase               | `{{ upper "hello" }}` -> `HELLO`                     |
| `trim`       | Trim whitespace                    | `{{ trim " hello " }}` -> `hello`                    |
| `trimPrefix` | Trim prefix                        | `{{ trimPrefix "v" "v1.0.0" }}` -> `1.0.0`           |
| `trimSuffix` | Trim suffix                        | `{{ trimSuffix ".tar.gz" "file.tar.gz" }}` -> `file` |
| `contains`   | Check if string contains substring | `{{ contains "hello" "he" }}` -> `true`              |
| `hasPrefix`  | Check if string has prefix         | `{{ hasPrefix "hello" "he" }}` -> `true`             |
| `hasSuffix`  | Check if string has suffix         | `{{ hasSuffix "hello" "lo" }}` -> `true`             |
| `replace`    | Replace substring                  | `{{ replace "hello" "l" "x" -1 }}` -> `hexxo`        |

### Comparison Functions

| Function | Description           | Example                                                 |
| -------- | --------------------- | ------------------------------------------------------- |
| `eq`     | Equal                 | `{{ eq .OS "linux" }}` -> `true` if OS is linux         |
| `ne`     | Not equal             | `{{ ne .OS "windows" }}` -> `true` if OS is not windows |
| `lt`     | Less than             | `{{ lt 1 2 }}` -> `true`                                |
| `le`     | Less than or equal    | `{{ le 1 1 }}` -> `true`                                |
| `gt`     | Greater than          | `{{ gt 2 1 }}` -> `true`                                |
| `ge`     | Greater than or equal | `{{ ge 1 1 }}` -> `true`                                |

### Logical Operators

| Function | Description | Example                           |
| -------- | ----------- | --------------------------------- |
| `and`    | Logical AND | `{{ and true false }}` -> `false` |
| `or`     | Logical OR  | `{{ or true false }}` -> `true`   |
| `not`    | Logical NOT | `{{ not true }}` -> `false`       |

### Mathematical Functions

| Function | Description    | Example                |
| -------- | -------------- | ---------------------- |
| `add`    | Addition       | `{{ add 1 2 }}` -> `3` |
| `sub`    | Subtraction    | `{{ sub 3 2 }}` -> `1` |
| `mul`    | Multiplication | `{{ mul 2 3 }}` -> `6` |
| `div`    | Division       | `{{ div 6 3 }}` -> `2` |
| `mod`    | Modulo         | `{{ mod 5 2 }}` -> `1` |

### List Functions

| Function | Description                  | Example                                                                    |
| -------- | ---------------------------- | -------------------------------------------------------------------------- |
| `list`   | Create a list                | `{{ list "a" "b" "c" }}` -> `[a b c]`                                      |
| `first`  | Get first item               | `{{ first (list "a" "b") }}` -> `a`                                        |
| `rest`   | Get all but first item       | `{{ rest (list "a" "b" "c") }}` -> `[b c]`                                 |
| `last`   | Get last item                | `{{ last (list "a" "b") }}` -> `b`                                         |
| `has`    | Check if list contains value | `{{ has (list "linux" "darwin") .OS }}` -> `true` if OS is linux or darwin |

## Templated Fields

The following fields in the tool configuration support templating:

| Field             | Example                                                                |
| ----------------- | ---------------------------------------------------------------------- |
| `output`          | `output: bin/{{ .OS }}-{{ .ARCH }}`                                    |
| `skip.condition`  | `condition: '{{ eq .OS "windows" }}'`                                  |
| `version`         | `version: v{{ if eq .OS "windows" }}0.1.0{{ else }}0.2.0{{ end }}`     |
| `source.type`     | `type: {{ if eq .OS "windows" }}url{{ else }}github{{ end }}`          |
| `exe.patterns`    | `patterns: ["^{{ .OS }}-{{ .Exe}}"]`                                   |
| `extensions`      | `extensions: ['{{ if eq .OS "windows" }}.exe{{ else }}{{ end }}']`     |
| `commands`        | `commands: ["pip install {{ .Exe }}=={{ .Version }}"]`                 |
| `hints[].pattern` | `pattern: "{{ .OS }}"`                                                 |
| `hints[].weight`  | `weight: '{{ if eq .ARCH "arm" }}1{{ else }}0{{ end }}'`               |
| `path`            | `path: "https://example.com/{{ .Name }}_{{ .OS }}_{{ .ARCH }}.tar.gz"` |

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

The `-` in the delimiters (`{{-` and `-}}`) removes whitespace before or after the template action.

## Loop Logic

You can use loops in templates:

```yaml
commands:
  post:
    - |-
      {{- range $key, $value := .Values.config -}}
      echo "{{ $key }}={{ $value }}" >> config.ini
      {{- end -}}
```

## Examples

### Platform-specific Paths

```yaml
path: https://example.com/downloads/{{ .Name }}-{{ .Version }}-{{ .OS }}-{{ .ARCH }}.tar.gz
```

### Skip Windows Installation

```yaml
skip:
  condition: '{{ eq .OS "windows" }}'
  reason: "Tool is not available on Windows"
```

### Architecture-specific Hints

```yaml
hints:
  - pattern: armv7
    weight: '{{ if eq .ARCH "arm" }}{{ if eq .ARCH_VERSION "7" }}10{{ else }}1{{ end }}{{ else }}0{{ end }}'
```

### OS-dependent Extensions

```yaml
extensions:
  - '{{ if eq .OS "windows" }}.exe{{ else }}{{ end }}'
  - '{{ if has .OS (list "windows" "darwin") }}.zip{{ else }}{{ end }}'
  - .tar.gz
```

### Value-based Templating

```yaml
values:
  repo: myrepo
  owner: myowner

source:
  type: github
  github:
    owner: "{{ .Values.owner }}"
    repo: "{{ .Values.repo }}"
```
