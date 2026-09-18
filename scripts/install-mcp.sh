#!/usr/bin/env bash
# Install the checksum-verified native release downloaded in chapter 00.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
. "$ROOT/scripts/common.sh"
. "$ROOT/scripts/versions.env"
[ $# -eq 0 ] || { [ $# -eq 1 ] && [ "$1" = --from-local ]; } || die 'Usage: scripts/install-mcp.sh [--from-local]'
asset="$(beans_mcp_asset)"
source_dir="$ROOT/dist/beans-mcp/$BEANS_MCP_VERSION"
if [ "${1:-}" != --from-local ]; then
  "$ROOT/scripts/download-mcp.sh"
fi
[ -f "$source_dir/$asset" ] && [ -f "$source_dir/SHA256SUMS" ] || die 'Run scripts/get-materials.sh first.'
expected="$(awk -v n="$asset" '$2 == n {print $1}' "$source_dir/SHA256SUMS")"
actual="$(sha256_file "$source_dir/$asset")"
[ -n "$expected" ] && [ "$expected" = "$actual" ] || die 'Adapter checksum mismatch; refusing installation.'
mkdir -p "$BIN_DIR"
target="$BIN_DIR/beans-mcp$(executable_suffix)"
install -m 0755 "$source_dir/$asset" "$target"
"$target" --version
