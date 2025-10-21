---
layout: default
title: Default Configuration
parent: Configuration
nav_order: 1
---

# Default Configuration

`godyl` uses a default configuration to provide sensible defaults for all tools. This can be overridden by providing your own `defaults.yml` file.

Furthermore, many values can be overridden as described in the [Configuration]({{ site.baseurl }}/configuration/index) section.

## Default Configuration File

The default configuration is embedded in the `godyl` binary and looks like this:

{% raw %}

```yaml
default:
  # Output as for example `.bin-linux-amd64/<tool>`
  output: .bin-{{ .OS }}-{{ .ARCH_LONG }}
  checksum:
    type: file
  exe:
    patterns:
      # Search for the executable in all subdirectories
      - "**/{{ .Exe }}*{{ .EXTENSION }}"
  hints:
    # General rules
    # Name of the executable (without extension)
    - pattern: "{{ .Exe }}"
      type: contains

    # .tar.gz format is commonly used
    - pattern: .tar.gz
      type: endswith

    # .gz format is commonly used
    - pattern: .gz
      type: endswith
      weight: 1

    # .exe extensions for Windows
    - pattern: .exe
      weight: '{{ if eq .OS "windows" }}1{{ else }}0{{ end }}'
      type: endswith

      # .zip extensions for Windows
    - pattern: .zip
      weight: '{{ if eq .OS "windows" }}1{{ else }}0{{ end }}'
      type: endswith

      # .zip extensions for macOS
    - pattern: .zip
      weight: '{{ if eq .OS "darwin" }}1{{ else }}0{{ end }}'
      type: endswith

    # ARM32 specifics
    # Give priority to the current version
    - pattern: armv{{ .ARCH_VERSION }}
      weight: 4
      type: contains

    # Less priority to the previous version
    - pattern: armv{{sub .ARCH_VERSION 1}}
      weight: 3
      type: contains

    # Even less priority to the previous version
    - pattern: armv{{sub .ARCH_VERSION 2}}
      weight: 2
      type: contains

    # Matches strings starting with 'arm' followed by any single character except 'v'
    - pattern: arm[^v]
      type: regex

    # Alpine Linux specifics
    # Match musl-based ARM binaries
    - pattern: musleabihf
      weight: |-
        {{- if and (eq .DISTRIBUTION "alpine") (eq .ARCH "arm") -}}
        1
        {{- else -}}
        0
        {{- end -}}
      type: contains

    # Avoid files ending with .txt (commonly the checksum.txt file)
    - pattern: .txt
      type: endswith
      match: excluded

    # Avoid files ending with .sha256 (commonly the checksum file)
    - pattern: .sha256
      type: endswith
      match: excluded
  source:
    type: github
    go:
      base: github.com
      download_if_missing: false
  mode: find
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

linux:
  inherit: default
  output: ~/.local/bin

windows:
  inherit: default
  output: $USERPROFILE\AppData\Local\Programs

darwin:
  inherit: default
  output: ~/bin

gitlab:
  inherit: default
  source:
    type: gitlab
    gitlab:
      no-token: true
```

{% endraw %}

See the [Tools]({{ site.baseurl }}/configuration/tools) for full configuration options.

You can compose new configurations using the `inherit` keyword.

As an example, the configuration below allows you to inherit from `default` and customize the output directory:

{% raw %}

```yaml
default: <as above>

my-custom-output-directory:
  inherit:
    - default
  output: /usr/local/bin
```

`inherit` entries are applied in the order they are defined, with any set field overriding the parent.

The example below shows how to inherit from `default` and override the `gitlab` source type to use a self-hosted GitLab instance,
supporting both the `gitlab` source type and optionally the `url` source type, passing the `PRIVATE-TOKEN` header.

```yaml
gitlab-selfhosted:
  inherit: default
  source:
    type: gitlab
    gitlab:
      server: https://gitlab.self-hosted.com
    url:
      headers:
        PRIVATE-TOKEN:
          - "{{ .Tokens.GitLab }}"
```

{% endraw %}
