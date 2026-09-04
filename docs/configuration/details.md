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
  # Require .zip extension on Windows, but allow deviations for other platforms
  - pattern: .zip
    match: '{{ if eq .OS "windows" }}required{{ else }}weighted{{ end }}'
    type: endswith
```

See [tools.yml](https://github.com/idelchi/godyl/blob/main/tools.yml) for various examples of hints.

## AI Fallback and Hints

Hints remain the primary, repeatable way to teach `godyl` how a project's
release assets are named. When the root `ai` setting is enabled, AI is used only
when the normal platform rules and configured hints produce equally ranked,
qualified candidates. Successful and unmatched results are never sent to a
model.

The selected model receives the target platform and only the tied best asset
names. The response must equal one of those names exactly; it cannot invent a
download URL or bypass a requirement, exclusion, invalid hint, or source error.
See [AI fallback]({{ site.baseurl }}/commands/index#ai-fallback) for
configuration.

Use `godyl install --suggest` when matching does not produce a usable result and
you want help improving the configuration rather than selecting an asset for
only the current run. Suggested replacement hints are parsed and rerun through
the deterministic matcher before being displayed as verified.

## Conditional Logic

You can use conditional logic to customize behavior based on the platform:

```yaml
# Skip Windows installation
skip:
  reason: "Tool is not available on Windows"
  condition: '{{ eq .OS "windows" }}'
```

## Custom Commands

You can run custom commands after installation (or without installation).

```yaml
commands:
  - "mkdir -p {{ .Output }}"
  - "chmod +x {{ .Output }}/{{ .Exe }}"
  - "{{ .Output }}/{{ .Exe }} --configure"
```

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
