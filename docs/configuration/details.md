---
layout: default
title: Details
parent: Configuration
nav_order: 3
---

{% raw %}

## Using Hints

Hints help `godyl` choose the right asset to download, especially when multiple similar assets are available. For example:

```yaml
hints:
  # Prefer the exact executable name
  - pattern: "*{{ .Exe }}*"
    weight: 1

  # Prefer .zip files on Windows
  - pattern: "*.zip"
    weight: '{{ if eq .OS "windows" }}1{{ else }}0{{ end }}'
```

Setting `must: true` requires that the pattern matches, otherwise the asset is excluded:

```yaml
hints:
  - pattern: "*{{ .OS }}*"
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
```

## Custom Commands

You can run custom commands after installation (or without installation).

```yaml
commands:
  - "mkdir -p {{ .Output }}"
  - "chmod +x {{ .Output }}/{{ .Exe }}"
  - "{{ .Output }}/{{ .Exe }} --configure"
```

## Alternative Installation Methods

For tools that aren't available as binaries, you can use commands:

```yaml
source:
  type: |
    {{- if eq .OS "windows" -}}
    url
    {{- else -}}
    github
    {{- end -}}
```

## Create subset tools.yml

Extract a subset of tools from the embedded tools.yml:

```sh
godyl dump tools --tags docker > docker-tools.yml
```

This creates a file containing only tools tagged with "docker".

## Platform Inference

`godyl` tries to infer platform details from asset names. Here's how different platforms are recognized:

### Operating Systems

| OS      | Inferred from           |
| :------ | :---------------------- |
| Linux   | linux                   |
| Darwin  | darwin, macos, mac, osx |
| Windows | windows, win            |
| FreeBSD | freebsd                 |
| Android | android                 |
| NetBSD  | netbsd                  |
| OpenBSD | openbsd                 |

### Architectures

| Architecture    | Inferred from                           |
| :-------------- | :-------------------------------------- |
| AMD64           | amd64, x86_64, x64, win64               |
| ARM64           | arm64, aarch64                          |
| AMD32           | amd32, x86, i386, i686, win32, 386, 686 |
| ARM32 (v7)      | armv7, armv7l, armhf                    |
| ARM32 (v6)      | armv6, armv6l                           |
| ARM32 (v5)      | armv5, armel                            |
| ARM32 (unknown) | arm                                     |

### Libraries

| Library    | Inferred from |
| :--------- | :------------ |
| GNU        | gnu           |
| Musl       | musl          |
| MSVC       | msvc          |
| LibAndroid | android       |

{% endraw %}
