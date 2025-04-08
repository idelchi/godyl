---
layout: default
title: Advanced Features
---

# Advanced Features

This page covers advanced features of Godyl, including templating, platform inference, and special use cases.

## Template Variables

Godyl supports templating in many configuration fields. The following variables are available:

| Variable              | Description                                               |
| --------------------- | --------------------------------------------------------- |
| `{{ .Name }}`         | The name of the tool or project                           |
| `{{ .Output }}`       | The output path template for built artifacts              |
| `{{ .Exe }}`          | The name of the executable                                |
| `{{ .Env.<> }}`       | Any environment variable                                  |
| `{{ .Values.<> }}`    | Custom values for templating                              |
| `{{ .Version }}`      | The version of the tool or project                        |
| `{{ .OS }}`           | The operating system (e.g., `linux`, `darwin`, `windows`) |
| `{{ .ARCH }}`         | The architecture type (e.g., `amd64`, `arm64`)            |
| `{{ .ARCH_VERSION }}` | The version of the architecture, if applicable            |
| `{{ .LIBRARY }}`      | The system library (e.g., `gnu`, `musl`)                  |
| `{{ .EXTENSION }}`    | The file extension specific to the platform               |
| `{{ .DISTRIBUTION }}` | The distribution name (e.g., `debian`, `alpine`)          |

## Templating Functions

All functions from the [slim-sprig](https://github.com/go-task/slim-sprig) library are available for use in templates. Some common functions include:

- `eq`: Equal comparison
- `ne`: Not equal comparison
- `lt`, `le`, `gt`, `ge`: Less than, less than or equal, greater than, greater than or equal
- `and`, `or`, `not`: Logical operations
- `contains`, `hasPrefix`, `hasSuffix`: String operations
- `sub`, `add`, `mul`, `div`: Arithmetic operations
- `list`: Create a list
- `has`: Check if a list contains a value

## Platform Inference

Godyl can infer platform details from asset names. Here's how different platforms are recognized:

### Operating Systems

| OS      | Inferred from           |
| ------- | ----------------------- |
| Linux   | linux                   |
| Darwin  | darwin, macos, mac, osx |
| Windows | windows, win            |
| FreeBSD | freebsd                 |
| Android | android                 |
| NetBSD  | netbsd                  |
| OpenBSD | openbsd                 |

### Architectures

| Architecture    | Inferred from                           |
| --------------- | --------------------------------------- |
| AMD64           | amd64, x86_64, x64, win64               |
| ARM64           | arm64, aarch64                          |
| AMD32           | amd32, x86, i386, i686, win32, 386, 686 |
| ARM32 (v7)      | armv7, armv7l, armhf                    |
| ARM32 (v6)      | armv6, armv6l                           |
| ARM32 (v5)      | armv5, armel                            |
| ARM32 (unknown) | arm                                     |

### Libraries

| Library    | Inferred from |
| ---------- | ------------- |
| GNU        | gnu           |
| Musl       | musl          |
| MSVC       | msvc          |
| LibAndroid | android       |

## Using Hints

Hints help Godyl choose the right asset to download, especially when multiple similar assets are available. For example:

```yaml
hints:
  # Prefer the exact executable name
  - pattern: "{{ .Exe }}"
    weight: 1

  # Prefer .zip files on Windows
  - pattern: zip
    weight: '{{ if eq .OS "windows" }}1{{ else }}0{{ end }}'

  # Prefer the correct ARM version
  - pattern: "armv{{.ARCH_VERSION}}"
    weight: 3
  - pattern: "armv{{sub .ARCH_VERSION 1}}"
    weight: 2
  - pattern: "armv{{sub .ARCH_VERSION 2}}"
    weight: 1
```

Setting `must: true` requires that the pattern matches, otherwise the asset is excluded:

```yaml
hints:
  - pattern: "{{ .OS }}"
    must: true
```

## Conditional Logic

You can use conditional logic to customize behavior based on the platform:

```yaml
# Skip Windows installation
skip:
  reason: "Tool is not available on Windows"
  condition: '{{ eq .OS "windows" }}'

# Use different version based on OS
version: |-
  {{- if eq .OS "windows" -}}
    v0.1.0
  {{- else -}}
    v0.2.0
  {{- end -}}

# Choose source type based on OS
source:
  type: |-
    {{- if eq .OS "windows" -}}
      url
    {{- else -}}
      github
    {{- end -}}
```

## Extension Handling

You can specify which file extensions to prefer:

```yaml
extensions:
  - '{{ if eq .OS "windows" }}.exe{{ else }}{{ end }}'
  - '{{ if has .OS (list "windows" "darwin") }}.zip{{ else }}{{ end }}'
  - .tar.gz
```

## Custom Commands

You can run custom commands before or after installation:

```yaml
commands:
  pre:
    - "mkdir -p {{ .Output }}"
  post:
    - "chmod +x {{ .Output }}/{{ .Exe }}"
    - "{{ .Output }}/{{ .Exe }} --configure"
```

## Alternative Installation Methods

For tools that aren't available as binaries, you can use the `commands` source type:

```yaml
source:
  type: commands
  commands:
    - "pip install {{ .Name }}=={{ .Version }}"
```

## YQ Filtering Example

Extract a subset of tools from the embedded tools.yml using yq:

```sh
godyl dump tools | yq --yaml-output '[.[] | try (select(.tags != null and (.tags[] == "docker")))]' > my-tools.yml
```

This creates a new tools.yml file containing only tools tagged with "docker".
