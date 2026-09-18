#!/usr/bin/env bash
# Maintainer build. Attendees download binaries; they do not need Go.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
. "$ROOT/scripts/common.sh"
. "$ROOT/scripts/versions.env"
[ $# -eq 0 ] || { [ $# -eq 1 ] && [ "$1" = --release ]; } || die 'Usage: scripts/build-mcp.sh [--release]'
commit="$(git -C "$ROOT" rev-parse HEAD)"
built_at="$(git -C "$ROOT" show -s --format=%cI HEAD)"
flags="-s -w -X main.version=$BEANS_MCP_VERSION -X main.commit=$commit -X main.buildAt=$built_at"
cd "$ROOT/mcp/beans"
if [ "${1:-}" != --release ]; then
  mkdir -p "$BIN_DIR"
  CGO_ENABLED=0 go build -trimpath -ldflags "$flags" -o "$BIN_DIR/beans-mcp$(executable_suffix)" ./cmd/beans-mcp
  exit 0
fi
dest="$ROOT/dist/beans-mcp/$BEANS_MCP_VERSION"
mkdir -p "$dest"
# Keep these aligned with scripts/common.sh and the native CI matrix.
for target in darwin/arm64 darwin/amd64 linux/amd64 linux/arm64 windows/amd64; do
  os="${target%/*}"; arch="${target#*/}"; suffix=''
  [ "$os" != windows ] || suffix=.exe
  info "Building $target"
  CGO_ENABLED=0 GOOS="$os" GOARCH="$arch" go build -trimpath -ldflags "$flags" \
    -o "$dest/beans-mcp_${BEANS_MCP_VERSION}_${os}_${arch}${suffix}" ./cmd/beans-mcp
done
copy_file LICENSE "$dest/LICENSE"
copy_file THIRD-PARTY-NOTICES.md "$dest/THIRD-PARTY-NOTICES.md"
jq -n --arg version "$BEANS_MCP_VERSION" --arg commit "$commit" --arg built_at "$built_at" \
  --arg go "$(go env GOVERSION)" \
  '{name:"beans-mcp",version:$version,commit:$commit,built_at:$built_at,go:$go,targets:["darwin/arm64","darwin/amd64","linux/amd64","linux/arm64","windows/amd64"]}' > "$dest/release-manifest.json"
cd "$dest"
for file in beans-mcp_* LICENSE THIRD-PARTY-NOTICES.md release-manifest.json; do
  printf '%s  %s\n' "$(sha256_file "$file")" "$file"
done > SHA256SUMS
ok "Release assets: $dest"
