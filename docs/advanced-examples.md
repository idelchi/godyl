---
layout: default
title: Advanced Examples
---

# Advanced Examples

This page provides advanced examples and patterns for using Godyl in complex scenarios.

## Creating a Development Environment

This example shows how to use Godyl to set up a complete development environment with various tools.

Create a `dev-tools.yml` file:

```yaml
# Go development tools
- name: golangci-lint
  source:
    type: github
    github:
      owner: golangci
      repo: golangci-lint
  tags:
    - go
    - linter

- name: gofumpt
  source:
    type: github
    github:
      owner: mvdan
      repo: gofumpt
  tags:
    - go
    - formatter

- name: delve
  source:
    type: github
    github:
      owner: go-delve
      repo: delve
  hints:
    - pattern: dlv
      weight: 2
  tags:
    - go
    - debugger

# Kubernetes tools
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

# Infrastructure tools
- name: terraform
  source:
    type: github
    github:
      owner: hashicorp
      repo: terraform
  tags:
    - infra
    - cli

- name: awscli
  skip:
    - condition: '{{ eq .OS "windows" }}'
      reason: "AWS CLI requires a different installation process on Windows"
  source:
    type: commands
    commands:
      - pip install awscli
  tags:
    - infra
    - cli
    - aws

# Shell tools
- name: jq
  source:
    type: github
    github:
      owner: stedolan
      repo: jq
  tags:
    - cli
    - json

- name: yq
  source:
    type: github
    github:
      owner: mikefarah
      repo: yq
  tags:
    - cli
    - yaml
```

Install all tools:

```sh
godyl install dev-tools.yml --output ~/.local/bin
```

Install only Go development tools:

```sh
godyl install dev-tools.yml --tags go --output ~/.local/bin
```

Install Kubernetes and infrastructure tools:

```sh
godyl install dev-tools.yml --tags 'k8s,infra' --output ~/.local/bin
```

## Platform-specific Installations

This example shows how to handle platform-specific configurations and alternatives.

```yaml
# Linux-specific tools
- name: htop
  skip:
    - condition: '{{ ne .OS "linux" }}'
      reason: "htop is only fully functional on Linux"
  source:
    type: github
    github:
      owner: htop-dev
      repo: htop
  tags:
    - linux
    - system

# Platform alternatives
- name: ps
  skip:
    - condition: '{{ eq .OS "linux" }}'
      reason: "Using htop instead on Linux"
  source:
    type: github
    github:
      owner: PowerShell
      repo: PowerShell
  hints:
    - pattern: pwsh
      weight: 1
  tags:
    - system

# Architecture-specific optimizations
- name: ripgrep
  source:
    type: github
    github:
      owner: BurntSushi
      repo: ripgrep
  hints:
    - pattern: '{{ if eq .ARCH "amd64" }}x86_64{{ else }}{{ .ARCH }}{{ end }}'
      weight: 3
    - pattern: "musl"
      weight: '{{ if eq .LIBRARY "musl" }}2{{ else }}0{{ end }}'
  tags:
    - cli
    - search
```

## Advanced Templating

This example shows complex templating techniques.

```yaml
# Dynamic version selection
- name: node
  version: |-
    {{- if eq .OS "windows" -}}
      18.16.1
    {{- else if eq .OS "darwin" -}}
      20.3.0
    {{- else -}}
      {{- if eq .ARCH "arm64" -}}
        20.3.0
      {{- else -}}
        18.16.1
      {{- end -}}
    {{- end -}}
  source:
    type: github
    github:
      owner: nodejs
      repo: node
  hints:
    - pattern: node-v{{ .Version }}-{{ .OS }}-{{ .ARCH }}
      weight: 10
  tags
```
