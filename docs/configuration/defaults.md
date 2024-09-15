---
layout: default
title: Default Configuration
parent: Configuration
nav_order: 1
---

# Default Configuration

`godyl` uses a default configuration to provide sensible defaults for all tools. This can be overridden by providing your own `defaults.yml` file.

Furthermore, many values can be overridden as described in the [Configuration]({{ site.baseurl }}/configuration/configuration) section.

## Default Configuration File

The default configuration is embedded in the `godyl` binary and looks like this:

{% raw %}

```yaml
default:
  output: .bin-{{ .OS }}-{{ .ARCH_LONG }}
  exe:
    patterns:
      - "^{{ .Exe }}{{ .EXTENSION }}$"
      - ".*/{{ .Exe }}{{ .EXTENSION }}$"
  hints:
    - pattern: "{{ .Exe }}"
    - pattern: "armv{{ .ARCH_VERSION }}"
      weight: 4
    - pattern: "armv{{sub .ARCH_VERSION 1}}"
      weight: 3
    - pattern: "armv{{sub .ARCH_VERSION 2}}"
      weight: 2
    - pattern: "arm[^v].*"
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
        header: Authorization
        scheme: Bearer
  strategy: sync
  version:
    commands:
      - --version
      - -v
      - -version
      - version
    patterns:
      - '.*?(v?\d+\.\d+\.\d+).*'
      - '.*?(v?\d+\.\d+).*'
```

See the [Tools]({{ site.baseurl }}/configuration/tools) for full configuration options.

You can compose new configurations using the `inherit` keyword.

As an example, the configuration below allows you to inherit from `default` and customize the output directory:

```yaml
default: <as above>

my-custom-output-directory:
  inherit:
    - default
  output: /usr/local/bin
```

`inherit´ are applied in the order they are defined, with any set field overriding the parent.

{% endraw %}
