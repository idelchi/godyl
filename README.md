docker pull --platform linux/arm64 ubuntu:20.04
docker pull --platform linux/arm/v7 ubuntu:20.04

docker run --rm -it --platform linux/arm64 ubuntu:latest
docker run --rm -it --platform linux/arm/v7 ubuntu:latest

# godyl

## Getting Started

### Prerequisites

- Go 1.23 or higher

### Installation

    go install github.com/idelchi/

### Usage

```
- name:
  version: {{ inferred from GitHub if source }}
  path: {{ templated last to allow population of all fields first }}
  exe: {{ inferred from name }}
  platform:
    os: {{ inferred if not given }}
    architecture: {{ inferred if not given }}
    library: {{ inferred if not given }}
  aliases: []
  values: {}
  fallbacks: []
  hints:
    - pattern:
      weight:
      regex:
      must:
  source: {}

- owner/repository
- path
```
