#!/bin/sh

set -eu

TEST_DIR=$(mktemp -d)
trap 'rm -rf "${TEST_DIR}"' EXIT

if [ "${#}" -eq 0 ]; then
	go run . dump tools --embedded >"${TEST_DIR}/tools.yml"
else
	TAGS=$(printf ',%s' "${@}")
	TAGS=${TAGS#,}
	go run . dump tools --embedded --tags="${TAGS}" >"${TEST_DIR}/tools.yml"
fi

while read -r OS ARCH; do
	go run . \
		--env-file=tokens.env \
		--no-cache \
		--tmp="${TEST_DIR}" \
		install "${TEST_DIR}/tools.yml" \
		--output="${TEST_DIR}/${OS}-${ARCH}" \
		--strategy=force \
		--os="${OS}" \
		--arch="${ARCH}"
done <<'EOF'
darwin arm64
windows amd64
linux arm64
linux amd64
EOF
