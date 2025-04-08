---
layout: default
title: Default Configuration
---

# Default Configuration

Godyl uses a default configuration to provide sensible defaults for all tools. This can be overridden by providing your own `defaults.yml` file.

Furthermore, many values can be overridden as described in the [Configuration](configuration/configuration#content-start) section.

## Default Configuration File

The default configuration is embedded in the Godyl binary and looks like this:

```yaml
output: .bin-{{ .OS }}-{{ .ARCH_LONG }}
exe:
  patterns:
    - "^{{ .Exe }}{{ .EXTENSION }}$"
    - ".*/{{ .Exe }}{{ .EXTENSION }}$"
hints:
  - pattern: "{{ .Exe }}"
    weight: 1
  - pattern: "armv{{ .ARCH_VERSION }}"
    weight: 4
  - pattern: "armv{{sub .ARCH_VERSION 1}}"
    weight: 3
  - pattern: "armv{{sub .ARCH_VERSION 2}}"
    weight: 2
  - pattern: "arm[^v].*"
    weight: 1
  - pattern: "musleabihf"
    weight: |-
      {{- if and (eq .DISTRIBUTION "alpine") (eq .ARCH "arm") -}}
      1
      {{- else -}}
      0
      {{- end -}}
extensions:
  - '{{ if eq .OS "windows" }}.exe{{ else }}{{ end }}'
  - '{{ if eq .OS "windows" }}.zip{{ else }}{{ end }}'
  - '{{ if eq .OS "darwin" }}.zip{{ else }}{{ end }}'
  - .tar.gz
mode: find
source:
  type: github
  url:
    token:
      scheme: Basic
strategy: none
env:
  GH_TOKEN: $GODYL_GITHUB_TOKEN
version:
  commands:
    - --version
    - -v
    - -version
    - version
  patterns:
    - '.*?(\d+\.\d+\.\d+).*'
    - '.*?(\d+\.\d+).*'
```
