#!/bin/sh
set -e

# Usage function
usage() {
  cat <<EOF
Usage: ${0} [OPTIONS]
Oneoff extracts a list of tools using godyl.

All arguments are passed to 'godyl download' command, and as such, you're advised to check the documentation at https://github.com/idelchi/godyl.

Environment variables:
GODYL_GITHUB_TOKEN: GitHub token to use for downloading assets from GitHub.
DISABLE_SSL: Disable SSL verification when downloading assets.

Example:

    curl -sSL https://raw.githubusercontent.com/idelchi/godyl/refs/heads/dev/scripts/oneoff.sh | sh -s -- -o ./bin idelchi/gogen

EOF
  exit 1
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
  tmp=$(mktemp -d)
  trap 'rm -rf "${tmp}"' EXIT

  curl ${DISABLE_SSL:+-k} -sSL \
    "https://raw.githubusercontent.com/idelchi/scripts/refs/heads/dev/install.sh" |
    INSTALLER_TOOL=godyl GODYL_GITHUB_TOKEN=${GODYL_GITHUB_TOKEN} \
      sh -s -- \
      -d "${tmp}" \
      ${DISABLE_SSL:+-k}

  [ -n "$@" ] && printf "Calling godyl download with arguments: '${command}'\n"

  # Install tools using godyl
  "${tmp}/godyl" --help
  exit 0
  GODYL_GITHUB_TOKEN=${GODYL_GITHUB_TOKEN} "${tmp}/godyl" download "${command}" ${DISABLE_SSL:+-k}

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
  # Check for required commands
  need_cmd curl

  # Install tools
  command="$@"

  install_tools "${command}"
}


main "$@"
