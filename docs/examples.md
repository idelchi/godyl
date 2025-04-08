---
layout: default
title: Examples
---

# Godyl Examples

This page provides practical examples of using Godyl for different scenarios.

## Basic Examples

### Download a Single Tool

Download the latest version of a tool from GitHub:

```sh
godyl download idelchi/godyl --output ~/.local/bin
```

### Install Multiple Tools from a YAML File

Create a `tools.yml` file:

```yaml
- name: godyl
  source:
    type: github
    github:
      owner: idelchi
      repo: godyl

- name: kubectl
  source:
    type: url
    path: https://dl.k8s.io/release/v1.27.3/bin/{{ .OS }}/{{ .ARCH }}/kubectl{{ .EXTENSION }}
  mode: extract
```

Then install the tools:

```sh
godyl install tools.yml --output ~/.local/bin
```

### Download a Specific Version

```sh
godyl download idelchi/godyl --version v0.1.0 --output ~/.local/bin
```

### Download for a Different Platform

```sh
godyl download idelchi/godyl --os linux --arch arm64 --output ~/.local/bin
```

## Advanced Examples

### Create a Kubernetes Tools Collection

Create a `k8s-tools.yml` file:

```yaml
- name: kubectl
  source:
    type: url
    path: https://dl.k8s.io/release/v1.27.3/bin/{{ .OS }}/{{ .ARCH }}/kubectl{{ .EXTENSION }}
  mode: extract
  tags:
    - k8s
    - cli

- name: helm
  source:
    type: github
    github:
      owner: helm
      repo: helm
  tags:
    - k8s
    - cli

- name: k9s
  source:
    type: github
    github:
      owner: derailed
      repo: k9s
  tags:
    - k8s
    - cli
    - tui

- name: kubectx
  source:
    type: github
    github:
      owner: ahmetb
      repo: kubectx
  tags:
    - k8s
    - cli
  hints:
    - pattern: kubectx
      weight: 2
```

Install all Kubernetes tools:

```sh
godyl install k8s-tools.yml --output ~/.local/bin
```

### Filter Tools by Tags

Install only CLI tools:

```sh
godyl install tools.yml --tags cli --output ~/.local/bin
```

Exclude TUI tools:

```sh
godyl install tools.yml --tags '!tui' --output ~/.local/bin
```

### Using Custom Commands

Create a `pip-tools.yml` file for Python packages:

```yaml
- name: black
  source:
    type: commands
    commands:
      - pip install black=={{ .Version }}
  version: "23.3.0"
  tags:
    - python
    - formatter

- name: mypy
  source:
    type: commands
    commands:
      - pip install mypy=={{ .Version }}
  version: "1.3.0"
  tags:
    - python
    - type-checker
```

Install Python tools:

```sh
godyl install pip-tools.yml
```

### Using a Different Strategy

Always reinstall tools:

```sh
godyl install tools.yml --strategy force --output ~/.local/bin
```

Upgrade tools if a newer version is available:

```sh
godyl install tools.yml --strategy upgrade --output ~/.local/bin
```

### Create Aliases for Tools

```yaml
- name: kubectl
  source:
    type: url
    path: https://dl.k8s.io/release/v1.27.3/bin/{{ .OS }}/{{ .ARCH }}/kubectl{{ .EXTENSION }}
  mode: extract
  aliases:
    - k
    - kctl
```

This will create symlinks (or copies on Windows) named `k` and `kctl` pointing to `kubectl`.

### Skip Tools Based on Platform

```yaml
- name: linux-only-tool
  source:
    type: github
    github:
      owner: example
      repo: linux-only-tool
  skip:
    - condition: '{{ ne .OS "linux" }}'
      reason: "This tool is only available on Linux"
```

### Using Custom Hints

```yaml
- name: terraform
  source:
    type: github
    github:
      owner: hashicorp
      repo: terraform
  hints:
    - pattern: "{{ .OS }}_{{ .ARCH }}"
      must: true
    - pattern: zip
      weight: 2
```

### Dumping and Reusing Embedded Tools

```sh
# Dump the embedded tools configuration
godyl dump tools > my-tools.yml

# Modify the file if needed

# Install the tools
godyl install my-tools.yml --output ~/.local/bin
```
