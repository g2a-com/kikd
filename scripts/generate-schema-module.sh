#!/usr/bin/env bash

set -e

REPO_DIR="$(git rev-parse --show-toplevel)"
MODULE_DIR="${REPO_DIR}/internal/schema"

# Install modules
if [ ! -d "${REPO_DIR}/node_modules" ]; then
	npm ci
fi

# Generate schemas.go file
(
	cd "${REPO_DIR}"
	echo '//go:generate ../../scripts/generate-schema-module.sh'
	echo 'package schema'
	echo
	echo '// THIS FILE WAS GENERATED USING SCRIPT, DO NOT CHANGE IT BY HAND!'
	echo '// RUN "go generate ./internal/schema" TO UPDATE.'
	echo
	echo 'var SCHEMAS = map[string][]byte{'
	for FILEPATH in assets/schemas/public/*/*/*.yaml; do
		API_VERSION="$(echo "${FILEPATH}" | cut -d / -f 4-5)"
		BASENAME=$(basename -s ".yaml" "${FILEPATH}")
		KIND=$(echo "${BASENAME:0:1}" | tr '[:lower:]' '[:upper:]')${BASENAME:1}

		if [ "${KIND}" = Index ]; then
			continue
		fi

		printf "\t\"%s/%s\": []byte(\`\n" "${API_VERSION}" "${KIND}"
		"${REPO_DIR}/scripts/normalize-schema.mjs" "${FILEPATH}" | sed -e 's/^/\t\t/;'
		printf "\t\`),\n"
	done
	echo '}'
) >"${MODULE_DIR}/schemas.go"
