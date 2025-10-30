---
layout: default
title: Tools Format
parent: Configuration
nav_order: 2
---

{% raw %}

# Tools Format

Tools can be defined in a YAML file. `godyl` will look for a file named `tools.yml` in the current directory, if not specified.

See [tools.yml](https://github.com/idelchi/godyl/blob/main/tools.yml) for a sample configuration.

You can use a full or simple form for the tool definitions. Most fields follow the same format, with a full and simple form.

## Configuration

```yaml
- idelchi/envprof
```

This is the simplest form to download the latest release of `envprof` from the the default source type (for example, `github` or `gitlab`).

For more complex configurations, you can use the extended form:

```yaml
- name: idelchi/envprof
  description: Environment profile manager
  exe:
    name: envprof
    patterns:
      - "**/{{ .Exe }}{{ .EXTENSION }}"
  source:
    type: github
  tags:
    - cli
    - env
  strategy: sync
```

A complete reference for all fields is available below.

### Full form

```yaml
# Name of the tool to download. Used as display name and for inferring other fields.
- name: idelchi/envprof
  # Optional description of the tool, for documentation purposes.
  description: Profile-based environment variable manager
  # Version tracking of the tool. Specifies the target version, as well as how to parse the current version.
  version:
    # For `github` and `gitlab` sources, leave empty to fetch the latest release from the API.
    # The version is always available as {{ .Version }}, except when not set.
    # It is then only available after the version has been determined.
    # Allows for using wildcards like `v1.*` or `1.2.*` to fetch the latest matching version.
    version: v0.1.0
    # Commands to run to get the current installed version (for syncs),
    # whenever not available in the cache.
    commands:
      - --version
      - version
    # Regex patterns to extract the version from command output (for syncs),
    # whenever not available in the cache.
    patterns:
      - '.*?(\d+\.\d+\.\d+).*'
  # The download url. For `github` and `gitlab` sources,
  # leave empty to populate from the API.
  url: "https://github.com/idelchi/envprof/releases/download/v0.0.1/envprof_{{ .OS }}_{{ .ARCH }}.tar.gz"
  # Checksum information to verify the download.
  checksum:
    # The type of checksum. Supported types are `sha256`, `sha512`, `sha1`, `md5`, `file`, and `none`.
    type: sha256|sha512|sha1|md5|file|none
    # Value can be one of the following:
    # For `Type=[sha256 sha512 sha1 md5]` it can be:
    #  - A value (the checksum)
    #  - when prefixed by `url:` or `path:`, a URL or file path containing a single checksum value (or a BSD or GNU style checksum file, see `entry`).
    # For `Type=file` it can be:
    #  - A URL or file path containing BSD or GNU style checksums
    #  - Empty to determine it from the source [gitlab, github].
    #      Either a checksum asset will be used, or the digest field from the asset (if available).
    value: "[abc123...|url:https://example.com/checksum.txt|path:./checksum.txt]|https://example.com/checksums.txt"
    # For `Type=file`, pattern to match to select the correct checksum file from the assets.
    pattern: "checksum*.txt"
    # Entry may be used when value contains `url:` or `path:` which points to a file with multiple entries.
    # Mainly used as workaround when the `go-getter` library cannot select the correct asset.
    # Try `type: file` first, and if it doesn't work, use `type: sha256` with `value: url:...` and `entry`.
    entry: "{{ .File }}"
  # The output directory where the tool will be placed.
  output: ~/.local/bin # [`--output`]
  # The executable name. Specifies the desired output name of the executable,
  # as well as the patterns to find it in the downloads.
  exe:
    # The name to use for the executable.
    # Will be inferred from `name` using source-specific rules if not provided.
    # If no suffix is provided, the platform-specific suffix will be added.
    name: envprof
    # Glob patterns to find the executable in the downloads.
    # Uses globstar, so you can use `**` to match any number of directories.
    patterns:
      - "**/{{ .OS }}-{{ .Exe }}*{{ .EXTENSION }}"
  # A list of fallback strategies to try if the main source strategy fails.
  # Will be used in the order they are defined.
  fallbacks:
    - go
  # Hints to find the correct asset to download.
  hints:
    # Pattern to match the asset name. See `type` for allowed syntax.
    - pattern: "*amd64*"
      # Weight of the hint. Defaults to 1 if not provided.
      weight: |-
        {{- if eq .ARCH "amd64" -}}
        1
        {{- else -}}
        0
        {{- end -}}
      # Method to use for matching the pattern.
      # default: glob
      type: glob|regex|globstar|startswith|endswith|contains
      # Whether to use the weight for matches, require the match, or exclude the match.
      # default: weighted
      match: weighted|required|excluded
  source:
    type: github|gitlab|url|go|none # [`--source`]
    github:
      # Inferred from first part of `name` if not provided
      owner: idelchi
      # Inferred from last part of `name` if not provided
      repo: envprof
      token: secret # [`--github-token`]
    gitlab:
      # Inferred from first part of `name` if not provided
      namespace: idelchi/go-projects
      # Inferred from last part of `name` if not provided
      project: envprof
      token: secret # [`--gitlab-token`]
      server: https://gitlab.self-hosted.com
      no-token: false # Suppress usage of token
    url:
      token: secret # [`--url-token`]
      headers:
        Authorization:
          - "Bearer {{ .Tokens.URL }}"
        Content-Type:
          - application/json
          - application/x-www-form-urlencoded
    go:
      # Specifies the base url of go projects.
      base: github.com
      # Specifies the path for the `go install` command.
      # Useful when the installable is not in the default path (e.g `cmd/<tool>` or `.`).
      command: cmd/envprof
      # Whether to download and install Go if not available locally.
      download_if_missing: true
  # Run custom commands after the installation (or only commands if `source.type` is `none`).
  commands:
    # The list of commands to run.
    commands:
      - "mkdir -p {{ .Output }}"
    # Whether to suppress failures in the commands.
    allow_failure: true
    # Whether to exit immediately on error.
    exit_on_error: false
  # List of tags to filter tools.
  tags:
    - env
  # Strategy for updating existing tools.
  strategy: none|sync|existing|force
  # Skip the tool if the condition is met.
  skip:
    - reason: "envprof is not available for Darwin"
      condition: '{{ eq .OS "darwin" }}'
  # The mode for downloading and installing the tool.
  # `find` will download, extract, and find the executable.
  # `extract` will download and extract directly to the output directory.
  mode: find|extract
  # A collection of arbitrary values.
  # Will be available as `{{ .Values.<name> }}` anywhere templating is supported.
  values:
    customScalar: scalarValue
    customMap:
      key1: value1
      key2: value2
    customList:
      - item1
      - item2
  # A collection of environment variables.
  # Will be accessible as `{{ .Env.<ENV_VAR> }}` anywhere templating is supported.
  env:
    GH_TOKEN: $GODYL_GITHUB_TOKEN
  # Disable SSL verification.
  no_verify_ssl: true
  # Disable cache usage
  no_cache: true
  # Disable checksum verification
  no_verify_checksum: true
  # A list of defaults to inherit from.
  inherit:
    - default
```

Most of the fields also support simplified forms which is described below.

## Templating

Many fields in the configuration support templating with variables like:

- `{{ .Name }}` - The name of the tool (`name`)
- `{{ .Env }}` - The environment variables (`env`), accessed with `{{ .Env.<ENV_VAR> }}`
- `{{ .Values }}` - The custom values (`values`), accessed with `{{ .Values.<name> }}`
- `{{ .Exe }}` - The executable name (`exe.name`)
- `{{ .Source }}` - The source type (`source.type`)
- `{{ .Output }}` - The output path (`output`)
- `{{ .Tokens }}` - The tokens for the source type (`source.tokens`), accessed with `{{ .Tokens.<source> }}`
- `{{ .Version }}` - The version to fetch (`version.version`).
- `{{ .URL }}` - The download URL (`url`), only available after the URL has been rendered.

> **Note:** If the `version.version` field is unset, the template variable will only be available after the API call has been made.

Special cases, only available in `checksum.value` and `checksum.entry`:

- `{{ .File }}` is the base name of the file in the `url`
- `{{ .Base }}` is the full url without `{{ .File }}`
- `{{ .URLWithoutExtensions }}` is the full url without any file extensions (multiple extensions removed)

Platform-specific variables are upper-cased and available as:

- `{{ .OS }}` - The operating system
- `{{ .ARCH }}` - The architecture
- `{{ .ARCH_ALIASES }}` - The architecture aliases as an array
- `{{ .ARCH_VERSION }}` - The architecture version
- `{{ .ARCH_LONG }}` - The architecture with version
- `{{ .IS_ARM }}` - Whether the architecture is ARM
- `{{ .IS_X86 }}` - Whether the architecture is x86
- `{{ .LIBRARY }}` - The library used for the platform
- `{{ .DISTRIBUTION }}` - The distribution used for the platform
- `{{ .EXTENSION }}` - The file extension for the platform

Examples:

```yaml
url: https://example.com/download/{{ .Name }}_{{ .OS }}_{{ .ARCH }}.tar.gz
url: https://releases.hashicorp.com/terraform/{{ .Version | trimPrefix "v" }}/terraform_{{ .Version | trimPrefix "v" }}_{{ .OS }}_{{ .ARCH }}.zip
```

## Available Fields

Below is a comprehensive list of fields that can be used to configure each tool.

For each tool, you can see whether it is required, supports templating, and whether it is exported as a template variable.

### `name`

🔴 Required • 🧩 Templated • 📤 Exports as: `{{ .Name }}`

The name of the tool to download. Used as display name and for inferring other fields.

```yaml
name: idelchi/envprof
```

Used for inference in [`exe`](#exe) and [`source`](#source)

### `description`

A description of the tool, for documentation purposes.

```yaml
description: Asset downloader for GitHub releases, URLs, and Go projects
```

### `version`

📤 Exports as: `{{ .Version }}`

The version of the tool to download. Will be inferred by the source type if not provided.

Simple form:

```yaml
version: v0.1.0
```

Full form:

```yaml
version:
  version: v0.1.0
  commands:
    - --version
  patterns:
    # Match "anything-v0.1.0" or "anything-0.1.0"
    - '.*?(v?\d+\.\d+).*'
```

### `url`

🧩 Templated • 📤 Exports as: `{{ .URL }}`

The url of the tool to download. Must be a URL to a file. Will be inferred by the source type if not provided.

```yaml
url: https://github.com/idelchi/envprof/releases/download/v0.1.0/envprof_linux_amd64.tar.gz
```

The most common use-case is to have it inferred from the `source` field configuration for the `github` and `gitlab` sources.

### `output`

🧩 Templated • 📤 Exports as: `{{ .Output }}`

The directory where the tool will be installed.

```yaml
output: ./bin/{{ .OS }}
```

### `exe`

📤 Exports as: `{{ .Exe }}`

Information about the executable. `exe.name` will be inferred from `name` and the source type if not explicitly provided.

Simple form:

```yaml
exe: envprof
```

Full form:

```yaml
exe:
  name: envprof
  patterns:
    - "**/{{ .Exe }}{{ .EXTENSION }}"
```

### `values`

📤 Exports as: `{{ .Values }}`

Arbitrary values that can be used in templates.

```yaml
values:
  protocol: https
```

Use as:

```yaml
url: {{ .Values.protocol }}://example.com/download/{{ .Name }}_{{ .OS }}_{{ .ARCH }}.tar.gz
```

### `fallbacks`

Fallback strategies if no matches were made in releases.

```yaml
fallbacks:
  - go
```

### `hints`

🧩 Templated

Hints to help `godyl` find the correct tool.

```yaml
hints:
  - pattern: "*{{ .Exe }}*"
    weight: 1
  - pattern: "^{{ .OS }}"
    type: regex
    match: required
```

If weight is not provided, it will be set to 1.

### `source`

🧩 Templated • 📤 Exports as: `{{ .Source }}`

Information about the source of the tool.

GitHub source:

```yaml
source:
  type: github
  github:
    repo: envprof
    owner: idelchi
    token:
```

URL source:

```yaml
source:
  type: url
  url:
    token:
    headers:
```

Go source:

```yaml
source:
  type: go
  go:
    base: github.com
    command: cmd/envprof
    download_if_missing: true
```

> **Note**: Choosing `go` as a source or fallback, without having a local installation and the `go` command available, will result in
> the download of the latest version of `go`, if `download_if_missing` is set to `true`.

#### Templating

Only the `token` fields support templating.

Furthermore, the `tokens` themselves are available as:

```yaml
{{ .Tokens.GitHub }}
{{ .Tokens.GitLab }}
{{ .Tokens.URL }}
```

### `commands`

🧩 Templated

Commands to run after the main source has been executed.

```yaml
commands:
  - |
    if ! command -v wget; then
      echo "wget is not installed. Please install wget to proceed."
      exit 1
    fi
  - mkdir -p {{ .Output }}
  - wget -qO- https://github.com/idelchi/envprof/releases/download/{{ .Version }}/envprof_{{ .OS }}_{{ .ARCH }}.tar.gz | tar -xz -C {{ .Output }}
```

### `tags`

Tags to filter tools.

```yaml
tags:
  - cli
  - downloader
```

The name of the current tool will always be added to the list of tags.

### `strategy`

Strategy for updating the tool.

```yaml
strategy: sync
```

Valid values:

- `none`: Skip if the tool already exists
- `sync`: Sync the tool to the desired version
- `existing`: Only sync if the tool already exists
- `force`: Always download and install

### `skip`

🧩 Templated

Conditions under which to skip the tool.

```yaml
skip:
  - condition: '{{ eq .OS "windows" }}'
    reason: "Tool is not available on Windows"
```

#### Templating

Only the `condition` field supports templating.

### `mode`

Mode for downloading and installing.

```yaml
mode: find
```

Valid values:

- `find`: Download, extract, and find the executable
- `extract`: Download and extract directly to the output directory

### `checksum`

🧩 Templated (only the `value`)

Checksum information to verify the download.

```yaml
checksum:
  type: sha256
  value: "abc123..."
  pattern: "checksum*.txt"
  entry: "{{ .File }}"
```

The combination `type: file` and empty `value` will fetch the checksum file from the source (only `github` & `gitlab` supported).

{% endraw %}
