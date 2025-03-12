#!/bin/sh
set -e

GITHUB_TOKEN=${GODYL_GITHUB_TOKEN:-${GITHUB_TOKEN}}

# Usage function
usage() {
  cat <<EOF
Usage: ${0} [OPTIONS]
Installs Kubernetes-related tools using godyl.

This script will install:
- helm
- kubectl (with alias 'kc')
- k9s
- kubectx
- kubens
- task

All arguments are passed to 'godyl install' command, and as such, you're advised to check the documentation at https://github.com/idelchi/godyl.

Environment variables:

  GODYL_GITHUB_TOKEN/GITHUB_TOKEN       GitHub token to use for downloading assets from GitHub.
  DISABLE_SSL                           Disable SSL verification when downloading assets.

Example:

    curl -sSL https://raw.githubusercontent.com/idelchi/godyl/refs/heads/dev/scripts/k8s.sh | sh -s -- -o ./bin
EOF
  exit 1
}

# Parse arguments
parse_args() {
  # Handle known options with getopts
  while getopts ":h" opt; do
    case "${opt}" in
      h) usage ;;
    esac
    shift $((OPTIND - 1))
    OPTIND=1
  done
}

# Create and handle temporary directory
setup_temp_dir() {
  if [ -z "${TEMP_DIR}" ]; then
    TEMP_DIR=$(mktemp -d)
    debug "Created temporary directory: ${TEMP_DIR}"
  else
    mkdir -p "${TEMP_DIR}"
    debug "Using specified temporary directory: ${TEMP_DIR}"
  fi

  # Set trap to clean up temporary directory
  trap 'rm -rf "${TEMP_DIR}"' EXIT
}

# Install godyl and tools
install_tools() {
  local args="${1}"

  tmp=$(mktemp -d)
  trap 'rm -rf "${tmp}"' EXIT

  curl ${DISABLE_SSL:+-k} -sSL \
    "https://raw.githubusercontent.com/idelchi/scripts/refs/heads/dev/install.sh" |
    INSTALLER_TOOL=godyl GODYL_GITHUB_TOKEN=${GITHUB_TOKEN} \
      sh -s -- \
      -d "${tmp}" \
      ${DISABLE_SSL:+-k}

  [ -n "${args}" ] && printf "Calling 'godyl install' with arguments: '${args}'\n"

  # Install tools using godyl
  GODYL_GITHUB_TOKEN=${GITHUB_TOKEN} "${tmp}/godyl" install ${args} ${DISABLE_SSL:+-k} - <<YAML
- name: helm/helm
  path: https://get.helm.sh/helm-{{ .Version }}-{{ .OS }}-{{ .ARCH }}.tar.gz
- name: kubernetes/kubernetes
  exe: kubectl
  path: https://dl.k8s.io/{{ .Version }}/bin/{{ .OS }}/{{ .ARCH }}/kubectl{{ .EXTENSION }}
  aliases: kc
- derailed/k9s
- name: ahmetb/kubectx
- name: ahmetb/kubectx
  exe: kubens
- name: go-task/task
YAML

  rm -rf "${tmp}"
  printf "All tools installed successfully\n"
}

need_cmd() {
  if ! command -v "${1}" >/dev/null 2>&1; then
    printf "Required command '${1}' not found"
    exit 1
  fi
}

main() {
  parse_args "$@"

  # Check for required commands
  need_cmd curl

  # Install tools
  args="$@"
  install_tools "${args}"
}

main "$@"
