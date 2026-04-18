#!/usr/bin/env bash
# Detect and apply a new Anytype API spec version.
#
# Reads the latest published version from
#   https://raw.githubusercontent.com/anyproto/anytype-api/refs/heads/main/docs/reference/versions.json
# compares it against the version currently pinned in the repo (the
# `APIVersion` constant in internal/client/client.go), and — if a newer
# version is available — rewrites the pinned version in:
#
#   - internal/client/client.go  (APIVersion constant)
#   - GNUmakefile                (OPENAPI_URL default)
#   - codegen/openapi.yaml       (re-downloaded from the new URL)
#
# and regenerates the IR + Go code via `make generate` when the generator
# binaries are available.
#
# Output: the resolved latest / current versions are printed to STDOUT, and
# machine-readable values are written to GITHUB_OUTPUT when the variable
# is set (i.e. when running inside GitHub Actions).
set -euo pipefail

VERSIONS_JSON_URL="${VERSIONS_JSON_URL:-https://raw.githubusercontent.com/anyproto/anytype-api/refs/heads/main/docs/reference/versions.json}"

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
CLIENT_FILE="$REPO_ROOT/internal/client/client.go"
MAKEFILE="$REPO_ROOT/GNUmakefile"
OPENAPI_FILE="$REPO_ROOT/codegen/openapi.yaml"

require() {
	command -v "$1" >/dev/null 2>&1 || {
		echo "error: required command '$1' is not on PATH" >&2
		exit 2
	}
}

require curl
require jq

emit_output() {
	local key="$1" value="$2"
	echo "$key=$value"
	if [[ -n "${GITHUB_OUTPUT:-}" ]]; then
		echo "$key=$value" >>"$GITHUB_OUTPUT"
	fi
}

current_version=$(grep -E '^\s*APIVersion\s*=' "$CLIENT_FILE" | head -1 | sed -E 's/.*"([^"]+)".*/\1/')
if [[ -z "$current_version" ]]; then
	echo "error: could not extract APIVersion from $CLIENT_FILE" >&2
	exit 1
fi

versions_json=$(curl -fsSL "$VERSIONS_JSON_URL")
latest_version=$(echo "$versions_json" | jq -r '.[0].version')
latest_download=$(echo "$versions_json" | jq -r '.[0].downloadUrl')

if [[ -z "$latest_version" || "$latest_version" == "null" ]]; then
	echo "error: could not parse latest version from $VERSIONS_JSON_URL" >&2
	exit 1
fi

emit_output current_version "$current_version"
emit_output latest_version "$latest_version"
emit_output latest_download_url "$latest_download"

if [[ "$current_version" == "$latest_version" ]]; then
	emit_output changed "false"
	echo "Already on the latest Anytype API version ($current_version)."
	exit 0
fi

emit_output changed "true"
echo "New Anytype API version available: $current_version -> $latest_version"

# 1. Update the pinned version string.
sed -i.bak -E "s|(APIVersion\s*=\s*\")[^\"]+(\")|\1$latest_version\2|" "$CLIENT_FILE"
rm -f "$CLIENT_FILE.bak"

# 2. Update the makefile default URL so `make fetch-spec` targets the new file.
sed -i.bak -E "s|(OPENAPI_URL\s*\\?=\\s*).*|\1$latest_download|" "$MAKEFILE"
rm -f "$MAKEFILE.bak"

# 3. Re-download the OpenAPI spec.
mkdir -p "$(dirname "$OPENAPI_FILE")"
curl -fsSL "$latest_download" -o "$OPENAPI_FILE"

# 4. Regenerate the IR + Framework code, if the generators are installed.
regenerate_if_possible() {
	if ! command -v tfplugingen-openapi >/dev/null 2>&1; then
		echo "tfplugingen-openapi not installed; skipping regeneration." >&2
		return
	fi
	if ! command -v tfplugingen-framework >/dev/null 2>&1; then
		echo "tfplugingen-framework not installed; skipping regeneration." >&2
		return
	fi
	(cd "$REPO_ROOT" && make generate-spec generate-code fmt)
}
regenerate_if_possible

echo "Anytype API bumped to $latest_version."
