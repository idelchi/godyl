---
layout: default
title: paths
parent: Commands
nav_order: 10
---

# Paths Command

`paths` prints out all active filesystem paths used by godyl.

These include:

- `config path`
- `cache path`
- `go path`
- `temp download path`

The config entry is the same active file printed by `godyl config path`. See
[Configuration File Location]({{ site.baseurl }}/configuration/index#configuration-file-location)
for its lookup order and platform-specific global location.

## Syntax

```sh
godyl [flags] paths
```
