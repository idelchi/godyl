#!/bin/sh
set -e

# Tool specific variables
TOOL="godyl"
DISABLE_SSL="${GODYL_DISABLE_SSL:-false}"

need_cmd() {
  if ! command -v "${1}" >/dev/null 2>&1; then
    printf "Required command '${1}' not found"
    exit 1
  fi
}

main() {
  # Check for required commands
  need_cmd curl

  # Call the installation script with the provided arguments
  curl ${DISABLE_SSL:+-k} -sSL https://raw.githubusercontent.com/idelchi/scripts/refs/heads/main/install.sh | INSTALLER_TOOL=${TOOL} sh -s -- "$@"
}

main "$@"
