#!/usr/bin/env bash
# Install the checksum-verified native release downloaded in chapter 00.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
. "$ROOT/scripts/common.sh"
. "$ROOT/scripts/versions.env"
[ $# -eq 0 ] || { [ $# -eq 1 ] && [ "$1" = --from-local ]; } || die 'Usage: scripts/install-mcp.sh [--from-local]'
case "$(uname -s)/$(uname -m)" in
  Darwin/arm64) platform=darwin_arm64 ;;
  Darwin/x86_64) platform=darwin_amd64 ;;
  Linux/aarch64|Linux/arm64) platform=linux_arm64 ;;
  Linux/x86_64) platform=linux_amd64 ;;
  *) die 'Supported hosts: macOS and Linux. Rehearsed host: macOS arm64.' ;;
esac
asset="beans-mcp_${BEANS_MCP_VERSION}_${platform}"
source_dir="$ROOT/dist/beans-mcp/$BEANS_MCP_VERSION"
[ -f "$source_dir/$asset" ] && [ -f "$source_dir/SHA256SUMS" ] || die 'Run scripts/get-materials.sh first.'
expected="$(awk -v n="$asset" '$2 == n {print $1}' "$source_dir/SHA256SUMS")"
actual="$(shasum -a 256 "$source_dir/$asset" | awk '{print $1}')"
[ -n "$expected" ] && [ "$expected" = "$actual" ] || die 'Adapter checksum mismatch; refusing installation.'
mkdir -p "$BIN_DIR"
install -m 0755 "$source_dir/$asset" "$BIN_DIR/beans-mcp"
"$BIN_DIR/beans-mcp" --version
